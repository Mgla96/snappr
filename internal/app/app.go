package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Mgla96/snappr/internal/errors"

	"github.com/Mgla96/snappr/internal/adapters/clients"
	"github.com/Mgla96/snappr/internal/config"

	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
)

type role string

const (
	system                  role = "system"
	user                    role = "user"
	ErrUnmarshalLLMResponse      = errors.New("error unmarshalling LLM response")
	ErrMissingGHToken            = errors.New("GH_TOKEN environment variable is required. Please set it before running this command")
	ErrMissingLLMToken           = errors.New("LLM_TOKEN environment variable is required. Please set it before running this command")
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// githubClient is an interface for interacting with github
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . githubClient
type githubClient interface {
	CreateBranch(ctx context.Context, owner, repo, newBranch, baseBranch string) error
	GetLatestCommitFromBranch(ctx context.Context, owner, repo, branch string) (string, error)
	AddCommitToBranch(ctx context.Context, owner, repo, branch, filePath, commitMessage string, fileContent []byte) error
	GetCommitCode(context context.Context, owner, repo, commitSHA string, codeFilter clients.CodeFilter) (map[string]string, error)
	AddCommentToPullRequestReview(ctx context.Context, owner, repo string, prNumber int, commentBody, commitID, path string, startLine, line int) (*github.PullRequestComment, error)
	CreatePullRequest(ctx context.Context, owner, repo, title, head, base, body string) (*github.PullRequest, error)
	MergePullRequest(ctx context.Context, owner, repo string, prNumber int, commitMessage string) (*github.PullRequestMergeResult, error)
	ListPullRequests(ctx context.Context, owner, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, error)
	GetPRCode(ctx context.Context, owner, repo string, prNumber int, opts *github.ListOptions, codeFilter clients.CodeFilter) (map[string]string, error)
	GetPRPatch(ctx context.Context, owner, repo string, prNumber int) (string, error)
}

// llmClient is an interface for interacting with a llm model
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . llmClient
type llmClient interface {
	GenerateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, model clients.ModelType) (string, error)
}

type PRCreation struct {
	// Title of the pull request
	Title string `json:"title"`
	// Body of the pull request
	Body string `json:"body"`
	// UpdatedFiles is a list of files that have been updated in the pull request
	UpdatedFiles []PRCreationFile `json:"updated_files"`
}

type PRCreationFile struct {
	// Path is the file path of the file that has been updated
	Path string `json:"path"`
	// FullContent is the full content of the file that has been updated
	FullContent   string `json:"full_content"`
	CommitMessage string `json:"commit_message"`
}

type PRChanges struct {
	Files []FileChange `json:"files"`
}

type FileChange struct {
	Path        string `json:"path"`
	FullContent string `json:"full_content"`
	Patch       string `json:"patch"`
}

type PRReviewMap map[string][]PRCommentInfo
type PRCommentInfo struct {
	CommentBody string
	// The start_line is the first line in the pull request diff that your multi-line comment applies to.
	StartLine int
	// The line of the blob in the pull request diff that the comment applies to. For a multi-line comment, the last line of the range that your comment applies to.
	Line int
	// The start_side is the starting side of the diff that the comment applies to. Can be LEFT or RIGHT.
	StartSide string
}

type App struct {
	cfg          *config.Config
	githubClient githubClient
	llmClient    llmClient
	log          zerolog.Logger
}

// New creates a new instance of the App.
func New(
	cfg *config.Config, githubClient githubClient, llmClient llmClient, logger zerolog.Logger) *App {
	app := &App{
		cfg:          cfg,
		githubClient: githubClient,
		llmClient:    llmClient,
		log:          logger,
	}

	return app
}

func (a *App) chatCompletionWithRetry(ctx context.Context, messages []openai.ChatCompletionMessage, model clients.ModelType, retries int) (string, error) {
	var response string
	var err error
	for i := 0; i < retries; i++ {
		response, err = a.llmClient.GenerateChatCompletion(ctx, messages, model)
		if err == nil {
			break
		}
	}
	return response, err
}

func extractKnowledgeSourceData(knowledgeSources string, cfgKnowledgeSources []config.KnowledgeSource) ([]string, error) {
	knowledgeSourceData := []string{}
	if knowledgeSources != "" {
		for _, source := range strings.Split(knowledgeSources, ",") {
			if source == "" {
				continue
			}
			data, err := RetrieveKnowledge(source, cfgKnowledgeSources)
			if err != nil {
				return nil, fmt.Errorf("error retrieving knowledge source: %w", err)
			}
			if data == NotImplementedMessage {
				// TODO(mgottlieb): Log a warning here
				continue
			}
			knowledgeSourceData = append(knowledgeSourceData, data)
		}
	}
	return knowledgeSourceData, nil
}

// ExecuteCreatePR executes the create PR workflow.
func (a *App) ExecuteCreatePR(ctx context.Context, commitSHA, branch, workflowName, knowledgeSources, fileRegexPattern string, printOnly bool) error {
	// Get code on GitHub from commit
	code, err := a.githubClient.GetCommitCode(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, commitSHA, clients.CodeFilter{
		FileRegexPattern: fileRegexPattern,
	})
	if err != nil {
		return fmt.Errorf("error getting commit code: %w", err)
	}

	if len(code) == 0 {
		return fmt.Errorf("no code found in commit")
	}

	// convert map[string]string where key is the file path and value is the file content to string
	// to a format the LLM can understand for context to help generate updated code
	codeJson, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	a.log.Debug().Msgf("Code to be reviewed: %s", codeJson)

	// Parse code and feed to LLM with prompt
	promptWorkflow := GetWorkflowByName(workflowName, a.cfg.Input.PromptWorkflows)
	if promptWorkflow == nil {
		return fmt.Errorf("workflow not found: %s", workflowName)
	}

	knowledgeSourceData, err := extractKnowledgeSourceData(knowledgeSources, a.cfg.Input.KnowledgeSources)
	if err != nil {
		return fmt.Errorf("error extracting knowledge source data: %w", err)
	}

	var messages []openai.ChatCompletionMessage

	for _, step := range promptWorkflow.Steps {
		messages = append(messages, openai.ChatCompletionMessage{Role: string(system), Content: step.Prompt})
	}
	for i, knowledgeSourceData := range knowledgeSourceData {
		messages = append(messages, openai.ChatCompletionMessage{Role: string(system), Content: "start knowledge context " + fmt.Sprintf("%d", i) + " {" + knowledgeSourceData + "} end knowledge context.\n"})
	}
	messages = append(messages, openai.ChatCompletionMessage{Role: string(user), Content: string(codeJson)})

	response, err := a.chatCompletionWithRetry(ctx, messages, a.cfg.LLM.DefaultModel, a.cfg.LLM.Retries)
	if err != nil {
		return fmt.Errorf("error generating chat completion: %w", err)
	}

	if response == "" {
		return fmt.Errorf("received empty response from LLM")
	}

	jsonFromResp := extractJSON(response)
	if printOnly {
		a.log.Info().Msg(jsonFromResp)
		return nil
	}

	// Get updated code from LLM response
	updatedCode, err := unmarshalTo[PRCreation]([]byte(jsonFromResp))
	if err != nil {
		return fmt.Errorf("error unmarshalling updated code JSON: %w: %v", ErrUnmarshalLLMResponse, err)
	}

	// Commit and push code to GitHub
	err = a.githubClient.CreateBranch(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, branch, "main")
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	for _, file := range updatedCode.UpdatedFiles {
		err = a.githubClient.AddCommitToBranch(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, branch, file.Path, file.CommitMessage, []byte(file.FullContent))
		if err != nil {
			return fmt.Errorf("failed to add commit to branch for file: %s, %w", file.Path, err)
		}
	}

	// Create pull request
	pr, err := a.githubClient.CreatePullRequest(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, updatedCode.Title, branch, "main", updatedCode.Body)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}
	a.log.Info().Msgf("Created pull request: %s", pr.GetHTMLURL())

	return nil
}

// ExecutePRReview executes the PR review workflow.
func (a *App) ExecutePRReview(ctx context.Context, commitSHA string, prNumber int, workflowName, knowledgeSources, fileRegexPattern string, printOnly bool) error {
	// Get code on GitHub from commit
	code, err := a.githubClient.GetPRCode(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, prNumber, nil, clients.CodeFilter{
		FileRegexPattern: fileRegexPattern,
	})
	if err != nil {
		return fmt.Errorf("error getting commit code: %w", err)
	}

	if len(code) == 0 {
		return fmt.Errorf("no code found in commit")
	}

	// convert map[string]string where key is the file path and value is the file content to string
	// to a format the LLM can understand for context to help generate updated code
	codeJson, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	a.log.Debug().Msgf("Code to be reviewed: %s", codeJson)

	rawDiff, err := a.githubClient.GetPRPatch(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, prNumber)
	if err != nil {
		return fmt.Errorf("error getting raw diff: %w", err)
	}

	parsedRawDiff := parseRawDiff(rawDiff)
	fileChanges := make([]FileChange, len(code))
	ctr := 0
	for filepath := range code {
		fileChanges[ctr] = FileChange{
			Path:        filepath,
			FullContent: code[filepath],
			Patch:       parsedRawDiff[filepath],
		}
		ctr += 1
	}

	prChanges := PRChanges{
		Files: fileChanges,
	}

	res, err := json.Marshal(prChanges)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	promptWorkflow := GetWorkflowByName(workflowName, a.cfg.Input.PromptWorkflows)
	if promptWorkflow == nil {
		return fmt.Errorf("workflow not found: %s", workflowName)
	}

	knowledgeSourceData, err := extractKnowledgeSourceData(knowledgeSources, a.cfg.Input.KnowledgeSources)
	if err != nil {
		return fmt.Errorf("error extracting knowledge source data: %w", err)
	}

	var messages []openai.ChatCompletionMessage

	for _, step := range promptWorkflow.Steps {
		messages = append(messages, openai.ChatCompletionMessage{Role: string(system), Content: step.Prompt})
	}

	for i, knowledgeSourceData := range knowledgeSourceData {
		messages = append(messages, openai.ChatCompletionMessage{Role: string(system), Content: "start knowledge context " + fmt.Sprintf("%d", i) + " {" + knowledgeSourceData + "} end knowledge context.\n"})
	}

	messages = append(messages, openai.ChatCompletionMessage{Role: string(user), Content: string(res)})

	response, err := a.chatCompletionWithRetry(ctx, messages, a.cfg.LLM.DefaultModel, a.cfg.LLM.Retries)
	if err != nil {
		return fmt.Errorf("error generating chat completion: %w", err)
	}

	if response == "" {
		return fmt.Errorf("received empty response from LLM")
	}
	jsonFromResp := extractJSON(response)
	if printOnly {
		fmt.Print(jsonFromResp)
		return nil
	}
	// Get updated code from LLM response
	prReviewMap, err := unmarshalTo[PRReviewMap]([]byte(jsonFromResp))
	if err != nil {
		return fmt.Errorf("error unmarshalling updated code JSON: %w: %v", ErrUnmarshalLLMResponse, err)
	}

	for filePath, prCommentInfoList := range prReviewMap {
		for _, prCommentInfo := range prCommentInfoList {
			prReviewComment, err := a.githubClient.AddCommentToPullRequestReview(ctx, a.cfg.Github.Owner, a.cfg.Github.Repo, prNumber, prCommentInfo.CommentBody, commitSHA, filePath, prCommentInfo.StartLine, prCommentInfo.Line)
			if err != nil {
				a.log.Error().Err(err).Msgf("Failed to add comment for filepath: %s", filePath)
			} else {
				a.log.Info().Msgf("added comment: %s", *prReviewComment.Body)
			}
		}
	}

	return nil
}

func parseRawDiff(diff string) map[string]string {
	diffMap := make(map[string]string)
	diffSections := strings.Split(diff, "diff --git")
	for _, section := range diffSections {
		if strings.TrimSpace(section) == "" {
			continue
		}
		lines := strings.SplitN(section, "\n", 2)
		header := lines[0]
		patch := lines[1]
		// Get the file path from the header
		headerParts := strings.Fields(header)
		filePath := strings.TrimPrefix(headerParts[len(headerParts)-1], "b/")
		diffMap[filePath] = patch
	}
	return diffMap
}

// CheckGHToken checks if the GH_TOKEN environment variable is set.
func CheckGHToken() error {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		return ErrMissingGHToken
	}
	return nil
}

// CheckLLMToken checks if the LLM_TOKEN environment variable is set.
func CheckLLMToken() error {
	token := os.Getenv("LLM_TOKEN")
	if token == "" {
		return ErrMissingLLMToken
	}
	return nil
}

func CheckTokens() error {
	err := CheckGHToken()
	if err != nil {
		return fmt.Errorf("error no github token: %w", err)
	}
	err = CheckLLMToken()
	if err != nil {
		return fmt.Errorf("error no llm token: %w", err)
	}
	return nil
}
