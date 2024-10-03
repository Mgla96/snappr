package app

import (
	"encoding/json"
	"fmt"
	"strings"
)

func unmarshalTo[T any](data []byte) (T, error) {
	var m T
	err := json.Unmarshal(data, &m)
	if err != nil {
		return m, fmt.Errorf("failed to unmarshal to %T: %w", m, err)
	}
	return m, nil
}

func extractJSON(response string) string {
	start := strings.Index(response, "{")
	if start == -1 {
		return "" // No JSON found if there is no "{" character
	}

	end := strings.LastIndex(response, "}")
	if end == -1 || end <= start {
		return "" // Improper JSON if no "}" or it's before "{"
	}

	return response[start : end+1] // Include the last "}" in the substring
}
