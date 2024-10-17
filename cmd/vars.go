package cmd

import "os"

var (
	repository       = os.Getenv("REPOSITORY")
	repositoryOwner  = os.Getenv("REPOSITORY_OWNER")
	commitSHA        = os.Getenv("COMMIT_SHA")
	prNumber         int
	printOnly        bool
	llmEndpoint      = os.Getenv("LLM_ENDPOINT")
	workflowName     = os.Getenv("WORKFLOW_NAME")
	branch           = os.Getenv("BRANCH_NAME")
	fileRegexPattern = os.Getenv("FILE_REGEX_PATTERN")
	llmRetries       int
	llmModel         = os.Getenv("LLM_MODEL")
	llmAPI           = os.Getenv("LLM_API")
	knowledgeSources = os.Getenv("KNOWLEDGE_SOURCES")
)