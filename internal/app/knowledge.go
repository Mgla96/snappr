package app

import (
	"fmt"
	"os"
)

type KnowledgeSourceType string

const (
	KnowledgeSourceTypeFile KnowledgeSourceType = "file"
	KnowledgeSourceTypeURL  KnowledgeSourceType = "url"
	KnowledgeSourceTypeAPI  KnowledgeSourceType = "api"
	KnowledgeSourceTypeText KnowledgeSourceType = "text"
	NotImplementedMessage   string              = "Not implemented"
)

type KnowledgeSource struct {
	Name  string              `mapstructure:"name"`
	Type  KnowledgeSourceType `mapstructure:"type"`
	Value string              `mapstructure:"value"`
}

var knowledgeSources []KnowledgeSource

func RetrieveKnowledge(sourceName string, knowledgeSources []KnowledgeSource) (string, error) {
	for _, source := range knowledgeSources {
		if source.Name == sourceName {
			switch source.Type {
			case KnowledgeSourceTypeFile:
				data, err := os.ReadFile(source.Value)
				if err != nil {
					return "", fmt.Errorf("could not read file: %w", err)
				}
				return string(data), nil
			case KnowledgeSourceTypeURL:
				return NotImplementedMessage, nil
			case KnowledgeSourceTypeAPI:
				return NotImplementedMessage, nil
			case KnowledgeSourceTypeText:
				return source.Value, nil
			}
		}
	}
	return "", nil
}
