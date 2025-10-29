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
	if data == nil || len(data) == 0 {
		return "{}", nil
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(bytes), nil
}

// Action constants for resource operations
const (
	readAction   = "read"
	createAction = "create"
	deleteAction = "delete"
	updateAction = "update"
)

// resourceActionError returns a formatted error message for resource actions
func resourceActionError(action, resourceType, err string) (string, string) {
	summary := fmt.Sprintf("Failed to %s %s", action, resourceType)
	detail := fmt.Sprintf("An error occurred while performing %s operation on %s: %s", action, resourceType, err)
	return summary, detail
}

// configResourceError returns a formatted error message for configuration errors
func configResourceError(resourceType string) (string, string) {
	summary := fmt.Sprintf("%s Resource Configuration Error", resourceType)
	detail := fmt.Sprintf("Expected configured %s client. Please report this issue to the provider developers.", resourceType)
	return summary, detail
}

// convertResourceError returns a formatted error message for conversion errors
func convertResourceError(resourceType, err string) (string, string) {
	summary := fmt.Sprintf("Failed to convert %s data", resourceType)
	detail := fmt.Sprintf("An error occurred while converting %s data: %s", resourceType, err)
	return summary, detail
}
