package app

import (
	"fmt"
	"os"

	"github.com/Mgla96/snappr/internal/config"

	"gopkg.in/yaml.v3"
)

const (
	promptBaseSystemPrompt = `You will only receive and respond with JSON and will act as an API.`
	promptCreatePR         = `As an expert Go software engineer, refactor Go code snippets from JSON (map[string]string, where key=file path, value=code snippet) to enhance performance, readability, and best practices adherence.
- **Return Format:** Provide code review feedback in a JSON that can unmarshal to the PRCreation spec below.
type PRCreation struct {
	// Title of the pull request
	Title string ` + "`json:\"title\"`" + `
	// Body of the pull request
	Body string ` + "`json:\"body\"`" + `
	// UpdatedFiles is a list of files that have been updated in the pull request
	UpdatedFiles []PRCreationFile ` + "`json:\"updated_files\"`" + `
}
type PRCreationFile struct {
	// Path is the file path of the file that has been updated
	Path string ` + "`json:\"path\"`" + `
	// FullContent is the full content of the file that has been updated.
	FullContent string ` + "`json:\"full_content\"`" + `
	CommitMessage string ` + "`json:\"commit_message\"`" + `
}
**Objective:** Only apply changes that clearly improve performance, readability, or best practices. Understand code objectives and provide new code to get closer to that objective. If no changes are needed for a file, do not add the file to the UpdatedFiles list. If you can solve a TODO comment in the code, please do so. If a change requires changes to other files that you can't update then just add a TODO comment.
`
	promptCodeReview = `As an expert Go software engineer tasked with code reviewing. Receive a JSON object of PRChanges struct below:
type PRChanges struct {
	Files []FileChange
}
type FileChange struct {
	Path        string
	FullContent string
	Patch       string
}
**Requirements:**
**Return Format:** Provide code review feedback in a JSON structured as:
// PRReviewMap is a map of file paths to a list of PRCommentInfo
type PRReviewMap map[string][]PRCommentInfo
type PRCommentInfo struct {
	CommentBody string
	// The start_line is the first line in the pull request diff that your multi-line comment applies to.
	StartLine int
	// The line of the blob in the pull request diff that the comment applies to. For a multi-line comment, the last line of the range that your comment applies to.
	Line int
	// Side of the diff for the last line. Can be LEFT or RIGHT.
	Side clients.Side
	// The start_side is the starting side of the diff that the comment applies to. Can be LEFT or RIGHT.
	StartSide clients.Side
}
**Objective:** Deliver actionable, line-specific feedback on only the code that was changed as part of the git diff. The git diff provides the exact lines you need to look at.
StartLine and Line must both fall within the same @@ hunk range.
`
)

type SnapprUserConfig struct {
	PromptWorkflows  []config.PromptWorkflow  `yaml:"promptWorkflows"`
	KnowledgeSources []config.KnowledgeSource `yaml:"knowledgeSources"`
}

// NewDefaultPromptAndKnowledgeConfig creates a config.yaml file with the default prompt and knowledge config.
func NewDefaultPromptAndKnowledgeConfig(configPath string) error {
	initUserConfig := SnapprUserConfig{
		PromptWorkflows: []config.PromptWorkflow{
			{
				Name: "createPR",
				Steps: []config.PromptWorkflowStep{
					{
						Prompt:      promptBaseSystemPrompt,
						InputSource: "text",
					},
					{
						Prompt:      promptCreatePR,
						InputSource: "text",
					},
				},
			},
			{
				Name: "codeReview",
				Steps: []config.PromptWorkflowStep{
					{
						Prompt:      promptBaseSystemPrompt,
						InputSource: "text",
					},
					{
						Prompt:      promptCodeReview,
						InputSource: "text",
					},
				},
			},
		},
		KnowledgeSources: []config.KnowledgeSource{
			{
				Name:  "exampleFileSource",
				Type:  config.KnowledgeSourceTypeFile,
				Value: "exampleFile.txt",
			},
			{
				Name:  "exampleTextSource",
				Type:  config.KnowledgeSourceTypeText,
				Value: "example text",
			},
			{
				Name:  "effectiveGo",
				Type:  config.KnowledgeSourceTypeURL,
				Value: "https://go.dev/doc/effective_go",
			},
		},
	}

	yamlData, err := yaml.Marshal(&initUserConfig)
	if err != nil {
		return fmt.Errorf("error while marshaling to YAML: %w", err)
	}

	err = os.WriteFile(configPath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error while writing to file: %w", err)
	}

	return nil
}
