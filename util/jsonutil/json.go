package jsonutil

import "encoding/json"

// Constants used for JSON serializations
const (
	jsonPrefix = ""
	jsonIndent = "  "
)

// SerializeJSON creates the JSON representation of a model
func SerializeJSON(target interface{}) ([]byte, error) {
	return json.MarshalIndent(target, jsonPrefix, jsonIndent)
}

// DeserializeJSON deserializes a JSON representation into a model
func DeserializeJSON(jsonData []byte, target interface{}) error {
	return json.Unmarshal(jsonData, target)
}
