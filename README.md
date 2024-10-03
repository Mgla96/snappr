# SnapPR

<div align="center">
<img src="assets/logo.jpeg" alt="Logo" width="250"/>

<p>
    <a href="LICENSE">
      <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"/>
    </a>
    <a href="https://goreportcard.com/report/github.com/Mgla96/snappr">
      <img src="https://goreportcard.com/badge/github.com/Mgla96/snappr" alt="Go Report Card"/>
    </a>
    <img src="https://img.shields.io/github/v/release/mgla6/snappr?sort=semver" alt="GitHub Release (latest SemVer)"/>
  </p>
</div>

----

SnapPR is a tool for snappy PR creation and review, helping developers save time and catch bugs earlier.

PR reviews can be used to add an additional static code analysis layer to your CI pipelines.

----

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)

## Installation

The `snappr` CLI can be installed by running:

```bash
go install github.com/Mgla96/snappr@latest
```

Alternatively, pre-compiled binaries are available under the `Assets` section of a release. To programatically download the binary, follow the instructions below:


To download the latest version with the provided command, you need curl and jq installed on your system.

```bash
# For Linux
curl -Lo /usr/local/bin/snappr $(curl -s https://api.github.com/repos/Mgla96/snappr/releases/latest | \
jq -r '.assets[] | select(.name == "snappr-linux-amd64") | .browser_download_url') && \
chmod +x /usr/local/bin/snappr
```

## Usage

```bash
snappr --help
Snappr is a tool for snappy PR creation and review to increase developer velocity.

Usage:
  snappr [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage configurations, workflows, and knowledge sources
  create      Create a new pull request
  help        Help about any command
  review      Review a pull request

Flags:
      --config string   config file (default "~/.snappr/config.yaml")
  -h, --help            help for snappr
  -t, --toggle          Help message for toggle

Use "snappr [command] --help" for more information about a command.
```

### Adding custom prompt workflows to use for PR review and creation.

A prompt workflow is a chain of prompts to send to the LLM as part of the PR review or creation. 

Default prompts are provided but if you prefer, you could create custom prompt workflows to use instead.

The config file used by SnapPR can be found at `~/.snappr/config.yaml` and additional prompt workflows can be added to the `promptWorkflows` array.


### Adding knowledge sources for improved PR review and creation.

Additional knowledge sources can be added to the `knowledgeSources` section found in the same SnapPR config file located at `~/.snappr/config.yaml`.

<!-- Info about how knowledge sources work -->

<!-- TODO: example of adding
https://go.dev/wiki/CodeReviewComments
for context on go PR reviews. Also example github workflow -->

<!-- TODO: example of how to use with ollama and llama -->
<!-- 
### Using Llama3.1 hosted with ollama

1. Follow the Ollama download instructions: https://ollama.com/download

1. Run llama 3.1

    ```bash
    ollama run llama3.1
    ```

1. Verify that you can access the model

    ```bash
    curl http://localhost:11434/api/generate -d '{
      "model": "llama3.1",
      "prompt": ""
    }'
    ```

1. Use the model with SnapPR

  ```bash
  
  ``` -->


### Snappy PR review

You can have SnapPR review a github PR with your custom prompt workflow and knowledge sources.

You can use the `-p` or `--printOnly` flag to print the output to STDOUT only instead of creating PR comments.

> [!IMPORTANT]  
> Make sure to source your `LLM_TOKEN` and `GH_TOKEN` environment variables for this command to succeed.
> You will also need to change some of these arguments to match your specific github repository.

```bash
snappr review --prNumber 1 --commitSHA 73151df1fdfe7a46c4aa57aaf25fcd698783442e -p --repository snappr --repositoryOwner Mgla96 --workflowName codeReview
{
    "/cmd/review.go": [
        {
            "CommentBody": "Consider using context passed from caller instead of `context.TODO()` for better control over context and its cancellation from caller functions.",
            "StartLine": 27,
            "Line": 27
        },
        {
            "CommentBody": "The command description declared as `Long:  'Create a new pull request'` might benefit from a more descriptive or detailed string explaining more about the command's functionality and its impact.",
            "StartLine": 17,
            "Line": 19
        },
        {
            "CommentBody": "Error strings should not be capitalized or end with punctuation as per Go conventions. Modify error message accordingly in `fmt.Errorf(\"error no github token: %w\", err)`.",
            "StartLine": 21,
            "Line": 23
        },
        {
            "CommentBody": "Logging fatal error directly in the app command setup could be harsh for an API consumer of this command. Consider returning an error and letting the caller handle it, potentially still exiting but with more control.",
            "StartLine": 38,
            "Line": 38
        },
        {
            "CommentBody": "Variable `repositoryOwner` is required but not used in the direct function body nor passed to any other interfaces, check if it needs to be part of this scope.",
            "StartLine": 11,
            "Line": 35
        },
        {
            "CommentBody": "It's good to provide context for regex but consider adding more examples or edge cases that this regex would fit to or reject, in variable `fileRegexPattern`, for a clearer understanding.",
            "StartLine": 58,
            "Line": 58
        },
        {
            "CommentBody": "Instead of directly calling os.Getenv(\"GH_TOKEN\") inside the function, consider passing it as a configuration or parameter, which helps in testing and managing configurations more effectively.",
            "StartLine": 34,
            "Line": 34
        },
        {
            "CommentBody": "Use explicit struct initialization to improve readability and avoid future bugs related to added fields in the `config.Config{...}` struct.",
            "StartLine": 30,
            "Line": 35
        },
        {
            "CommentBody": "Ensure functional options or configurations are validated or defaulted properly before use, particularly for the logger setup in the init function.",
            "StartLine": 66,
            "Line": 66
        }
    ],
    "/cmd/create.go": [
        {
            "CommentBody": "Ensure to handle possible empty strings or unexpected inputs for flags like 'repository', which could lead to runtime issues if unchecked.",
            "StartLine": 59,
            "Line": 59
        },
        {
            "CommentBody": "Review the necessity of using `context.TODO()` here; if the context is supposed to be cancellable or have a timeout, use the proper context.",
            "StartLine": 17,
            "Line": 17
        },
        {
            "CommentBody": "Check the consistency and necessity of `llmRetries`. If the retry logic can be encapsulated elsewhere or made more configurable, it could improve the command's robustness.",
            "StartLine": 11,
            "Line": 11
        },
        {
            "CommentBody": "The string description for `workflowName` seems redundant if it's marked as 'required'. Consider relocating detailed descriptions or utility explanations to a documentation section rather than inline.",
            "StartLine": 61,
            "Line": 61
        },
        {
            "CommentBody": "Follow up on the handling of environment variables: os.Getenv could return an empty string if not set, ensure there's a fallback or error handling mechanism.",
            "StartLine": 34,
            "Line": 34
        }
    ]
}
```

### Snappy PR creation

You can have SnapPR create a github PR with your custom prompt workflow and knowledge sources.

You can use the `-p` or `--printOnly` flag to print the output to STDOUT only instead of creating the Pull Request.

> [!IMPORTANT]  
> Make sure to source your `LLM_TOKEN` and `GH_TOKEN` environment variables for this command to succeed.
> You will also need to change some of these arguments to match your specific github repository.

```bash
snappr create --branch super-cool-pr --commitSHA a29e1c417e2e54076d5e88ba2e935989fb93fe1e --repository snappr --repositoryOwner Mgla96 --workflowName createPR
2024-10-02T22:18:26-07:00 INF Created pull request: https://github.com/Mgla96/snappr/pull/1

snappr create --branch super-cool-pr --commitSHA a29e1c417e2e54076d5e88ba2e935989fb93fe1e -p --repository snappr --repositoryOwner Mgla96 --workflowName createPR
2024-10-02T22:16:06-07:00 INF {
  "title": "Refactor Go Code for Better Performance and Readability",
  "body": "This pull request includes updates to enhance performance, readability, and adherence to best practices for several Go files in our codebase.",
  "updated_files": [
    {
      "path": "cmd/create.go",
      "full_content": "package cmd\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"os\"\n\n\t\"github.com/Mgla96/snappr/internal/app\"\n\t\"github.com/Mgla96/snappr/internal/config\"\n\t\"github.com/rs/zerolog\"\n\t\"github.com/spf13/cobra\"\n)\n\n// Variables to store command line flags\nvar (\n\tbranch, fileRegexPattern, repository, repositoryOwner, commitSHA, workflowName, llmEndpoint string\n\tllmRetries int\n\tprintOnly bool\n)\n\n// createCmd represents the create command\nvar createCmd = \u0026cobra.Command{\n\tUse:   \"create\",\n\tShort: \"Create a new pull request\",\n\tLong:  `Create a new pull request`,\n\tPreRunE: func(cmd *cobra.Command, args []string) error {\n\t\terr := app.CheckGHToken()\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"error no github token: %w\", err)\n\t\t}\n\t\terr = app.CheckLLMToken()\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"error no llm token: %w\", err)\n\t\t}\n\t\treturn nil\n\t},\n\tRun: func(cmd *cobra.Command, args []string) {\n\t\tlogger := zerolog.New(os.Stderr).With().Timestamp().Logger()\n\t\tconfig := \u0026config.Config{\n\t\t\tLog: config.Log{Level: zerolog.InfoLevel},\n\t\t\tGithub: config.Github{Token: os.Getenv(\"GH_TOKEN\"), Owner: repositoryOwner, Repo: repository},\n\t\t\tLLM: config.LLM{Token: os.Getenv(\"LLM_TOKEN\"), DefaultModel: \"gpt-4-turbo\", Endpoint: llmEndpoint},\n\t\t}\n\t\tapplication := app.SetupNoEnv(config)\n\t\terr := application.ExecuteCreatePR(context.TODO(), commitSHA, branch, workflowName, fileRegexPattern, printOnly)\n\t\tif err != nil {\n\t\t\tlogger.Fatal().Err(err).Msg(\"Error executing Create PR\")\n\t\t}\n\t},\n}\n\nfunc init() {\n\trootCmd.AddCommand(createCmd)\n\tflagSet := createCmd.Flags()\n\tflagSet.StringVar(\u0026repository, \"repository\", \"\", \"Github repository to review such as snappr (required)\")\n\tflagSet.StringVar(\u0026repositoryOwner, \"repositoryOwner\", \"\", \"The account owner of the repository. The name is not case sensitive. (required)\")\n\tflagSet.StringVar(\u0026commitSHA, \"commitSHA\", \"\", \"Commit SHA to create PR from (required)\")\n\tflagSet.StringVar(\u0026branch, \"branch\", \"\", \"Branch name to create PR from (required)\")\n\tflagSet.StringVar(\u0026workflowName, \"workflowName\", \"\", \"Prompt workflow to use (required)\")\n\tflagSet.StringVar(\u0026fileRegexPattern, \"fileRegexPattern\", `.*\\.go$`, \"Define a regex pattern to filter files to use as context for PR creation\")\n\tflagSet.StringVar(\u0026llmEndpoint, \"llmEndpoint\", \"\", \"Endpoint for the LLM service (defaults to OpenAI)\")\n\tflagSet.BoolVarP(\u0026printOnly, \"printOnly\", \"p\", false, \"Print the created PR only\")\n\tflagSet.IntVarP(\u0026llmRetries, \"llmRetries\", \"r\", 3, \"Number of retries for LLM API calls when failing to get a valid response\")\n\tmandatoryFlags := []string{\"repository\", \"repositoryOwner\", \"commitSHA\", \"branch\", \"workflowName\"}\n\tfor _, flag := range mandatoryFlags {\n\t\terr := createCmd.MarkFlagRequired(flag)\n\t\tif err != nil {\n\t\t\tlogger.Fatal().Err(err).Msg(\"Error marking \" + flag + \" as required\")\n\t\t}\n\t}\n}",
      "commit_message": "Refactor create.go: Improve flag definition and initialization readability"
    },
    {
      "path": "internal/app/convert.go",
      "full_content": "package app\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"strings\"\n)\n\n// unmarshalTo unmarshals JSON data into the specified type T, returning an error with detailed information upon failure.\nfunc unmarshalTo[T any](data []byte) (T, error) {\n\tvar result T\n\terr := json.Unmarshal(data, \u0026result)\n\tif err != nil {\n\t\treturn result, fmt.Errorf(\"failed to unmarshal to %T: %w\", result, err)\n\t}\n\treturn result, nil\n}\n\n// extractJSON searches a string for the first JSON object it contains and returns the JSON string if found.\n// If no valid JSON is found, it returns an empty string.\nfunc extractJSON(response string) string {\n\tstart := strings.Index(response, \"{\")\n\tif start == -1 {\n\t\treturn \"\" // No JSON found if there's no '{' character\n\t}\n\n\tend := strings.LastIndex(response, \"}\")\n\tif end == -1 || end \u003c= start {\n\t\treturn \"\" // No valid JSON present\n\t}\n\n\treturn response[start : end+1] // Include the last '}' in the substring\n}",
      "commit_message": "Update convert.go: Enhance comments and variable naming"
    }
  ]
}
```
