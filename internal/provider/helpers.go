package provider

import (
	"encoding/json"
	"fmt"
)

// parseJSON parses a JSON string into a map
func parseJSON(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return make(map[string]interface{}), nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return result, nil
}

// toJSONString converts a map to a JSON string
func toJSONString(data map[string]interface{}) (string, error) {
	if len(data) == 0 {
		return "{}", nil
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(bytes), nil
}
