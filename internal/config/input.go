package config

type KnowledgeSourceType string

const (
	KnowledgeSourceTypeFile KnowledgeSourceType = "file"
	KnowledgeSourceTypeURL  KnowledgeSourceType = "url"
	KnowledgeSourceTypeAPI  KnowledgeSourceType = "api"
	KnowledgeSourceTypeText KnowledgeSourceType = "text"
)

type InputConfig struct {
	KnowledgeSources []KnowledgeSource `mapstructure:"knowledgeSources"`
	PromptWorkflows  []PromptWorkflow  `mapstructure:"promptWorkflows"`
}

type KnowledgeSource struct {
	Name  string              `mapstructure:"name"`
	Type  KnowledgeSourceType `mapstructure:"type"`
	Value string              `mapstructure:"value"`
}

type PromptWorkflowStep struct {
	Prompt      string `mapstructure:"prompt"`
	InputSource string `mapstructure:"inputSource"`
}

type PromptWorkflow struct {
	Name  string               `mapstructure:"name"`
	Steps []PromptWorkflowStep `mapstructure:"steps"`
}
