package app

import (
	"fmt"
	"os"

	"github.com/Mgla96/snappr/internal/config"
)

const (
	NotImplementedMessage string = "Not implemented"
)

// RetrieveKnowledge retrieves knowledge from a knowledge source.
func RetrieveKnowledge(sourceName string, knowledgeSources []config.KnowledgeSource) (string, error) {
	for _, source := range knowledgeSources {
		if source.Name == sourceName {
			switch source.Type {
			case config.KnowledgeSourceTypeFile:
				data, err := os.ReadFile(source.Value)
				if err != nil {
					return "", fmt.Errorf("could not read file: %w", err)
				}
				return string(data), nil
			case config.KnowledgeSourceTypeURL:
				return NotImplementedMessage, nil
			case config.KnowledgeSourceTypeAPI:
				return NotImplementedMessage, nil
			case config.KnowledgeSourceTypeText:
				return source.Value, nil
			}
		}
	}
	return "", nil
}
