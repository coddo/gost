package models

import (
	"encoding/json"
)

// Constants used for JSON serializations
const (
	jsonPrefix = ""
	jsonIndent = "  "
)

// Create the JSON representation of a model
func SerializeJson(target interface{}) ([]byte, error) {
	return json.MarshalIndent(target, jsonPrefix, jsonIndent)
}

// Deserialize a JSON representation into a model
func DeserializeJson(jsonData []byte, target interface{}) error {
	return json.Unmarshal(jsonData, target)
}

// Interface for switching between dbmodels and models
type Expander interface {
	PopConstrains()
}
