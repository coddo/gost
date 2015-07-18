package models

import (
	"encoding/json"
)

// Constants used for JSON serializations
const (
	JsonPrefix = ""
	JsonIndent = "  "
)

// Create the JSON representation of a model
func SerializeJson(model Expander) ([]byte, error) {
	return json.MarshalIndent(model, JsonPrefix, JsonIndent)
}

// Deserialize a JSON representation into a model
func DeserializeJson(jsonData []byte, model Expander) error {
	return json.Unmarshal(jsonData, model)
}

// Interface for switching between dbmodels and models
type Expander interface {
	PopConstrains()
}
