package clients

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	ghClient *github.Client
	log      zerolog.Logger
}

// NewGithubClient creates a new instance of the GithubClient.
func NewGithubClient(token string, logger zerolog.Logger) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GithubClient{
		ghClient: client,
		log:      logger,
	}
}

// CreateBranch creates a new branch from the specified base branch.
func (gc *GithubClient) CreateBranch(ctx context.Context, owner, repo, newBranch, baseBranch string) error {
	baseRef, _, err := gc.ghClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return fmt.Errorf("failed to get base branch reference: %w", err)
	}

	newRef := &github.Reference{
		Ref:    github.String("refs/heads/" + newBranch),
		Object: &github.GitObject{SHA: baseRef.Object.SHA},
	}

	_, _, err = gc.ghClient.Git.CreateRef(ctx, owner, repo, newRef)
	if err != nil {
		return fmt.Errorf("failed to create new branch reference: %w", err)
	}

	return nil
}

// GetLatestCommitFromBranch gets the latest commit SHA from a branch.
func (gc *GithubClient) GetLatestCommitFromBranch(ctx context.Context, owner, repo, branch string) (string, error) {
	ref, _, err := gc.ghClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return "", fmt.Errorf("failed to get branch reference: %w", err)
	}

	commitSHA := ref.GetObject().GetSHA()
	return commitSHA, nil
}

// AddCommitToBranch adds a commit to the specified branch.
func (gc *GithubClient) AddCommitToBranch(ctx context.Context, owner, repo, branch, filePath, commitMessage string, fileContent []byte) error {
	ref, _, err := gc.ghClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return err
	}

	blob := &github.Blob{
		Content:  github.String(string(fileContent)),
		Encoding: github.String("utf-8"),
	}
	blobRes, _, err := gc.ghClient.Git.CreateBlob(ctx, owner, repo, blob)
	if err != nil {
		return err
	}

	treeEntry := &github.TreeEntry{
		Path: github.String(filePath),
		Mode: github.String("100644"),
		Type: github.String("blob"),
		SHA:  blobRes.SHA,
	}
	baseTree, _, err := gc.ghClient.Git.GetTree(ctx, owner, repo, *ref.Object.SHA, false)
	if err != nil {
		return err
	}
	tree, _, err := gc.ghClient.Git.CreateTree(ctx, owner, repo, *baseTree.SHA, []*github.TreeEntry{treeEntry})
	if err != nil {
		return err
	}

	parentCommit, _, err := gc.ghClient.Git.GetCommit(ctx, owner, repo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	commit := &github.Commit{
		Message: github.String(commitMessage),
		Tree:    tree,
		Parents: []*github.Commit{parentCommit},
	}
	newCommit, _, err := gc.ghClient.Git.CreateCommit(ctx, owner, repo, commit)
	if err != nil {
		return err
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err = gc.ghClient.Git.UpdateRef(ctx, owner, repo, ref, false)
	if err != nil {
		return err
	}

	return nil
}

// CodeFilter is used to filter files based on a regex pattern.
type CodeFilter struct {
	FileRegexPattern string
}

// GetCommitCode gets the code from a commit.
func (gc *GithubClient) GetCommitCode(context context.Context, owner, repo, commitSHA string, codeFilter CodeFilter) (map[string]string, error) {
	commit, _, err := gc.ghClient.Git.GetCommit(context, owner, repo, commitSHA)
	if err != nil {
		return nil, err
	}

	treeSHA := commit.Tree.GetSHA()
	tree, _, err := gc.ghClient.Git.GetTree(context, owner, repo, treeSHA, true)
	if err != nil {
		return nil, fmt.Errorf("error getting git tree: %w", err)
	}

	files := make(map[string]string)
	// temporary so we are not processing same files
	Shuffle(tree.Entries)
	for _, entry := range tree.Entries {
		if len(files) >= 5 {
			break // Stop processing once we have 5 entries. Temporary fix for context length limit reached
		}
		if entry.GetType() != "blob" {
			continue
		}

		re, err := regexp.Compile(codeFilter.FileRegexPattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex: %w", err)
		}
		if re.MatchString(entry.GetPath()) {
			blob, resp, err := gc.ghClient.Git.GetBlob(context, owner, repo, entry.GetSHA())
			if err != nil {
				return nil, fmt.Errorf("error getting blob: %w", err)
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return nil, fmt.Errorf("unexpected status code getting blob: %d", resp.StatusCode)
			}
			content, err := base64.StdEncoding.DecodeString(blob.GetContent())
			if err != nil {
				return nil, fmt.Errorf("error base64 decoding blob string: %w", err)
			}
			// if first line has substring DO NOT EDIT then skip
			if IsDoNotEditFile(content) {
				continue
			}

			files[entry.GetPath()] = string(content)
		}
	}

	return files, nil
}

// AddCommentToPullRequestReview adds a comment to a pull request review.
//
// Parameters:
//   - ctx: The context for the API request.
//   - owner: The owner of the repository.
//   - repo: The repository name.
//   - prNumber: The pull request number.
//   - commentBody: The content of the comment to be added.
//   - commitID: The SHA of the commit to comment on.
//   - path: The file path in the repository where the comment should be added.
//   - position: The position in the diff where the comment should be added.
//
// Returns:
//   - The created PullRequestComment object.
//   - An error if any occurred during the API request.
func (gc *GithubClient) AddCommentToPullRequestReview(ctx context.Context, owner, repo string, prNumber int, commentBody, commitID, path string, startLine, line int) (*github.PullRequestComment, error) {
	comment := &github.PullRequestComment{
		// Text content of the comment
		Body: github.String(commentBody),
		// SHA of the commit to comment on
		CommitID: github.String(commitID),
		// Filepath which the comment applies
		Path: github.String(path),
		// Position in the diff where the comment should be applied
		// Position: github.Int(position),
		// First line of range you want to comment on
		StartLine: github.Int(startLine),
		// Last line of range you want to comment on
		Line: github.Int(line),
	}

	prComment, _, err := gc.ghClient.PullRequests.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	return prComment, nil
}

// CreatePullRequest creates a new pull request.
func (gc *GithubClient) CreatePullRequest(ctx context.Context, owner, repo, title, head, base, body string) (*github.PullRequest, error) {
	pr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body),
	}

	pullRequest, _, err := gc.ghClient.PullRequests.Create(ctx, owner, repo, pr)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request: %w", err)
	}

	return pullRequest, nil
}

// MergePullRequest merges a pull request.
func (gc *GithubClient) MergePullRequest(ctx context.Context, owner, repo string, prNumber int, commitMessage string) (*github.PullRequestMergeResult, error) {
	mergeResult, _, err := gc.ghClient.PullRequests.Merge(ctx, owner, repo, prNumber, commitMessage, &github.PullRequestOptions{})
	if err != nil {
		return nil, fmt.Errorf("error merging pull request: %w", err)
	}

	return mergeResult, nil
}

// ListPullRequests lists all pull requests in a repository.
func (gc *GithubClient) ListPullRequests(ctx context.Context, owner, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	pullRequests, _, err := gc.ghClient.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing pull requests: %w", err)
	}

	return pullRequests, nil
}

// GetPRCode gets the code from a pull request.
func (gc *GithubClient) GetPRCode(ctx context.Context, owner, repo string, prNumber int, opts *github.ListOptions) (map[string]string, error) {
	commitFiles, _, err := gc.ghClient.PullRequests.ListFiles(ctx, owner, repo, prNumber, opts)
	if err != nil {
		return nil, fmt.Errorf("error listing PR files: %w", err)
	}

	files := make(map[string]string)
	for _, commitFile := range commitFiles {
		if strings.HasSuffix(commitFile.GetFilename(), ".go") {
			blob, resp, err := gc.ghClient.Git.GetBlob(ctx, owner, repo, commitFile.GetSHA())
			if err != nil {
				return nil, fmt.Errorf("error getting blob: %w", err)
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return nil, fmt.Errorf("unexpected status code getting blob: %d", resp.StatusCode)
			}
			content, err := base64.StdEncoding.DecodeString(blob.GetContent())
			if err != nil {
				return nil, fmt.Errorf("error base64 decoding blob string: %w", err)
			}

			files[commitFile.GetFilename()] = string(content)
		}
	}

	return files, nil
}

func (gc *GithubClient) GetPRDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	diffOpts := &github.RawOptions{Type: github.Diff}
	diff, _, err := gc.ghClient.PullRequests.GetRaw(ctx, owner, repo, prNumber, *diffOpts)
	if err != nil {
		gc.log.Error().Err(err).Msg("Failed to retrieve the pull request diff")
		return "", fmt.Errorf("error retrieving PR diff: %w", err)
	}

	return diff, nil
}

func (gc *GithubClient) GetPRPatch(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	diffOpts := &github.RawOptions{Type: github.Patch}
	diff, _, err := gc.ghClient.PullRequests.GetRaw(ctx, owner, repo, prNumber, *diffOpts)
	if err != nil {
		gc.log.Error().Err(err).Msg("Failed to retrieve the pull request diff")
		return "", fmt.Errorf("error retrieving PR diff: %w", err)
	}

	return diff, nil
}

// Shuffle shuffles a slice of any type
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		panic("Shuffle: not a slice")
	}

	rand.Seed(time.Now().UnixNano())
	n := rv.Len()
	swap := reflect.Swapper(slice)

	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}

// IsDoNotEditFile checks if the file has a DO NOT EDIT comment.
func IsDoNotEditFile(data []byte) bool {
	lines := bytes.SplitN(data, []byte("\n"), 2)
	if len(lines) > 0 && bytes.Contains(lines[0], []byte("DO NOT EDIT")) {
		return true
	}
	return false
}
