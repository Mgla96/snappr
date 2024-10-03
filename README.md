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
snappr review --prNumber 3 --commitSHA 637ca803061762711dc7064378ed69e96b4fca5e -p --repository snappr --repositoryOwner Mgla96 --workflowName codeReview
2024-09-08T14:03:31-07:00 INF {
  "/cmd/add.go": [
    {
      "CommentBody": "Consider validating inputs in 'Run' functions to ensure robustness.",
      "StartLine": 14,
      "Line": 14
    },
    {
      "CommentBody": "This change aligns with the project's new configuration management handling. Ensure corresponding documentation reflects these adjustments.",
      "StartLine": 22,
      "Line": 22
    },
    {
      "CommentBody": "Use structured logging instead of 'fmt.Println' for better production traceability.",
      "StartLine": 15,
      "Line": 15
    }
  ],
  "/cmd/root.go": [
    {
      "CommentBody": "Ensure that 'config' module properly abstracts configuration handling and isn't tightly coupled with other application components.",
      "StartLine": 8,
      "Line": 8
    },
    {
      "CommentBody": "This change correctly addresses previous misuse of import. Verify that all relevant tests are updated to reflect this change in module dependency.",
      "StartLine": 12,
      "Line": 12
    }
  ],
  "/internal/adapters/clients/github.go": [
    {
      "CommentBody": "Ensure new method adheres to the interface commitments and error handling conventions used throughout the application.",
      "StartLine": 38,
      "Line": 38
    },
    {
      "CommentBody": "Adding method comments is a good practice. Make sure the comments are comprehensive and cover all aspects of the function behavior.",
      "StartLine": 119,
      "Line": 121
    },
    {
      "CommentBody": "Verify that 'CodeFilter' and 'GetCommitCode' methods do not introduce any security implications or data handling issues.",
      "StartLine": 253,
      "Line": 256
    }
  ],
  "/internal/app/app.go": [
    {
      "CommentBody": "Removing 'InputConfig' and relying on the separate configuration approach aligns with best practices of separation of concerns.",
      "StartLine": 15,
      "Line": 15
    }
  ],
  "/internal/app/workflow.go": [
    {
      "CommentBody": "Good encapsulation practice by abstracting workflow logic into the 'config' module. Ensure comprehensive tests cover these changes to prevent regression.",
      "StartLine": 1,
      "Line": 1
    }
  ],
  "/internal/config/input.go": [
    {
      "CommentBody": "This file initializes new configuration structures well. Next steps include implementing validation and proper integration points.",
      "StartLine": 3,
      "Line": 31
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
snappr create --branch super-cool-pr --commitSHA 6700d0448ade1695d1f38b432f8c72cc1c7bb54b -p --repository snappr --repositoryOwner Mgla96 --workflowName createPR
{"level":"info","time":"2024-09-14T00:15:43-07:00","message":"create called"}
2024-09-14T00:16:05-07:00 INF {
  "title": "Refactor Go Code for Better Practices and Performance",
  "body": "This pull request includes code refactoring for better performance, readability, and adherence to best practices in the provided Go modules. Changes target missing implementations, best practices in error handling, and improved structuring.",
  "updated_files": [
    {
      "path": "cmd/validate.go",
      "full_content": "package cmd\n\nimport (\n\t\"fmt\"\n\n\t\"github.com/spf13/cobra\"\n)\n\nvar validateCmd = &cobra.Command{\n\tUse:   \"validate\",\n\tShort: \"Validate all configured workflows and knowledge sources\",\n\tRun: func(cmd *cobra.Command, args []string) {\n\t\tvalidateWorkflows()\n\t\tvalidateKnowledgeSources()\n\t},\n}\n\nfunc validateWorkflows() {\n\tif inputConfig.PromptWorkflows == nil {\n\t\tfmt.Println(\"No workflows configured\")\n\t\treturn\n\t}\n\tfmt.Println(\"Validation not yet implemented\")\n\n}\n\nfunc validateKnowledgeSources() {\n\tif len(inputConfig.KnowledgeSources) == 0 {\n\t\tfmt.Println(\"No knowledge sources configured\")\n\t\treturn\n\t}\n\tfmt.Println(\"Validation not yet implemented\")\n}\n",
      "commit_message": "Implement proper checks in validation commands and refactor for readability."
    },
    {
      "path": "internal/app/knowledge.go",
      "full_content": "package app\n\nimport (\n\t\"fmt\"\n\t\"os\"\n)\n\ntype KnowledgeSourceType string\n\nconst (\n\tKnowledgeSourceTypeFile KnowledgeSourceType = \"file\"\n\tKnowledgeSourceTypeURL  KnowledgeSourceType = \"url\"\n\tKnowledgeSourceTypeAPI  KnowledgeSourceType = \"api\"\n\tKnowledgeSourceTypeText KnowledgeSourceType = \"text\"\n)\n\ntype KnowledgeSource struct {\n\tName  string              `mapstructure:\"name\"`\n\tType  KnowledgeSourceType `mapstructure:\"type\"`\n\tValue string              `mapstructure:\"value\"`\n}\n\nvar knowledgeSources []KnowledgeSource\n\nfunc GetAllKnowledgeSources() []KnowledgeSource {\n\treturn knowledgeSources\n}\n\nfunc RetrieveKnowledge(sourceName string) string {\n\tfor _, source := range knowledgeSources {\n\t\tif source.Name == sourceName {\n\t\t\tswitch source.Type {\n\t\t\tcase KnowledgeSourceTypeFile:\n\t\t\t\tdata, err := os.ReadFile(source.Value)\n\t\t\t\tif err != nil {\n\t\t\t\t\tfmt.Printf(\"Error reading file: %s\\n\", err)\n\t\t\t\t\treturn \"\"\n\t\t\t\t}\n\t\t\t\treturn string(data)\n\t\t\tcase KnowledgeSourceTypeURL,\n\t\t\t\t KnowledgeSourceTypeAPI,\n\t\t\t\t KnowledgeSourceTypeText:\n\t\t\t\treturn fmt.Sprintf(\"Retrieval not implemented for type: %s\", source.Type)\n\t\t\t}\n\t\t}\n\t}\n\treturn \"\"\n}\n\nfunc AddKnowledgeSource(source KnowledgeSource) {\n\tknowledgeSources = append(knowledgeSources, source)\n}\n\nfunc RemoveKnowledgeSource(sourceName string) {\n\tfor i, source := range knowledgeSources {\n\t\tif source.Name == sourceName {\n\t\t\tknowledgeSources = append(knowledgeSources[:i], knowledgeSources[i+1:]...)\n\t\t\treturn\n\t\t}\n\t}\n}\n",
      "commit_message": "Handle errors gracefully and add messages for unimplemented types in RetrieveKnowledge."
    }
  ]
}
```
