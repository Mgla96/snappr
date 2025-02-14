<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# app

```go
import "github.com/Mgla96/snappr/internal/app"
```

Package app provides the main business logic for the snappr application.

## Index

- [Constants](<#constants>)
- [func CheckGHToken\(\) error](<#CheckGHToken>)
- [func CheckLLMToken\(\) error](<#CheckLLMToken>)
- [func CheckTokens\(\) error](<#CheckTokens>)
- [func GetWorkflowByName\(name string, workflowList \[\]config.PromptWorkflow\) \*config.PromptWorkflow](<#GetWorkflowByName>)
- [func NewDefaultPromptAndKnowledgeConfig\(configPath string\) error](<#NewDefaultPromptAndKnowledgeConfig>)
- [func RetrieveKnowledge\(sourceName string, knowledgeSources \[\]KnowledgeSource\) \(string, error\)](<#RetrieveKnowledge>)
- [type App](<#App>)
  - [func New\(cfg \*config.Config, githubClient githubClient, llmClient llmClient, logger zerolog.Logger\) \*App](<#New>)
  - [func Setup\(\) \(\*App, error\)](<#Setup>)
  - [func SetupNoEnv\(cfg \*config.Config\) \*App](<#SetupNoEnv>)
  - [func \(a \*App\) ExecuteCreatePR\(ctx context.Context, commitSHA, branch, workflowName, fileRegexPattern string, printOnly bool\) error](<#App.ExecuteCreatePR>)
  - [func \(a \*App\) ExecutePRReview\(ctx context.Context, commitSHA string, prNumber int, workflowName, fileRegexPattern string, printOnly bool\) error](<#App.ExecutePRReview>)
- [type FileChange](<#FileChange>)
- [type KnowledgeSource](<#KnowledgeSource>)
- [type KnowledgeSourceType](<#KnowledgeSourceType>)
- [type PRChanges](<#PRChanges>)
- [type PRCommentInfo](<#PRCommentInfo>)
- [type PRCreation](<#PRCreation>)
- [type PRCreationFile](<#PRCreationFile>)
- [type PRReviewMap](<#PRReviewMap>)
- [type SnapprUserConfig](<#SnapprUserConfig>)


## Constants

<a name="ErrUnmarshalLLMResponse"></a>

```go
const (
    ErrUnmarshalLLMResponse = errors.New("error unmarshalling LLM response")
    ErrMissingGHToken       = errors.New("GH_TOKEN environment variable is required. Please set it before running this command")
    ErrMissingLLMToken      = errors.New("LLM_TOKEN environment variable is required. Please set it before running this command")
)
```

<a name="KnowledgeSourceTypeFile"></a>

```go
const (
    KnowledgeSourceTypeFile KnowledgeSourceType = "file"
    KnowledgeSourceTypeURL  KnowledgeSourceType = "url"
    KnowledgeSourceTypeAPI  KnowledgeSourceType = "api"
    KnowledgeSourceTypeText KnowledgeSourceType = "text"
    NotImplementedMessage   string              = "Not implemented"
)
```

<a name="CheckGHToken"></a>
## func CheckGHToken

```go
func CheckGHToken() error
```

CheckGHToken checks if the GH\_TOKEN environment variable is set.

<a name="CheckLLMToken"></a>
## func CheckLLMToken

```go
func CheckLLMToken() error
```

CheckLLMToken checks if the LLM\_TOKEN environment variable is set.

<a name="CheckTokens"></a>
## func CheckTokens

```go
func CheckTokens() error
```



<a name="GetWorkflowByName"></a>
## func GetWorkflowByName

```go
func GetWorkflowByName(name string, workflowList []config.PromptWorkflow) *config.PromptWorkflow
```

GetWorkflowByName returns workflow information by name from a list of workflows.

<a name="NewDefaultPromptAndKnowledgeConfig"></a>
## func NewDefaultPromptAndKnowledgeConfig

```go
func NewDefaultPromptAndKnowledgeConfig(configPath string) error
```

NewDefaultPromptAndKnowledgeConfig creates a config.yaml file with the default prompt and knowledge config.

<a name="RetrieveKnowledge"></a>
## func RetrieveKnowledge

```go
func RetrieveKnowledge(sourceName string, knowledgeSources []KnowledgeSource) (string, error)
```



<a name="App"></a>
## type App



```go
type App struct {
    // contains filtered or unexported fields
}
```

<a name="New"></a>
### func New

```go
func New(cfg *config.Config, githubClient githubClient, llmClient llmClient, logger zerolog.Logger) *App
```

New creates a new instance of the App.

<a name="Setup"></a>
### func Setup

```go
func Setup() (*App, error)
```

Setup sets up the application utilizing environment variables.

<a name="SetupNoEnv"></a>
### func SetupNoEnv

```go
func SetupNoEnv(cfg *config.Config) *App
```

SetupNoEnv sets up the application from a config struct instead of utilizing environment variables.

<a name="App.ExecuteCreatePR"></a>
### func \(\*App\) ExecuteCreatePR

```go
func (a *App) ExecuteCreatePR(ctx context.Context, commitSHA, branch, workflowName, fileRegexPattern string, printOnly bool) error
```

ExecuteCreatePR executes the create PR workflow.

<a name="App.ExecutePRReview"></a>
### func \(\*App\) ExecutePRReview

```go
func (a *App) ExecutePRReview(ctx context.Context, commitSHA string, prNumber int, workflowName, fileRegexPattern string, printOnly bool) error
```

ExecutePRReview executes the PR review workflow.

<a name="FileChange"></a>
## type FileChange



```go
type FileChange struct {
    Path        string `json:"path"`
    FullContent string `json:"full_content"`
    Patch       string `json:"patch"`
}
```

<a name="KnowledgeSource"></a>
## type KnowledgeSource



```go
type KnowledgeSource struct {
    Name  string              `mapstructure:"name"`
    Type  KnowledgeSourceType `mapstructure:"type"`
    Value string              `mapstructure:"value"`
}
```

<a name="KnowledgeSourceType"></a>
## type KnowledgeSourceType



```go
type KnowledgeSourceType string
```

<a name="PRChanges"></a>
## type PRChanges



```go
type PRChanges struct {
    Files []FileChange `json:"files"`
}
```

<a name="PRCommentInfo"></a>
## type PRCommentInfo



```go
type PRCommentInfo struct {
    CommentBody string
    StartLine   int
    Line        int
}
```

<a name="PRCreation"></a>
## type PRCreation



```go
type PRCreation struct {
    // Title of the pull request
    Title string `json:"title"`
    // Body of the pull request
    Body string `json:"body"`
    // UpdatedFiles is a list of files that have been updated in the pull request
    UpdatedFiles []PRCreationFile `json:"updated_files"`
}
```

<a name="PRCreationFile"></a>
## type PRCreationFile



```go
type PRCreationFile struct {
    // Path is the file path of the file that has been updated
    Path string `json:"path"`
    // FullContent is the full content of the file that has been updated
    FullContent   string `json:"full_content"`
    CommitMessage string `json:"commit_message"`
}
```

<a name="PRReviewMap"></a>
## type PRReviewMap



```go
type PRReviewMap map[string][]PRCommentInfo
```

<a name="SnapprUserConfig"></a>
## type SnapprUserConfig



```go
type SnapprUserConfig struct {
    PromptWorkflows  []config.PromptWorkflow  `yaml:"promptWorkflows"`
    KnowledgeSources []config.KnowledgeSource `yaml:"knowledgeSources"`
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
