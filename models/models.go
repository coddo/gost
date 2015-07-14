package models

import "go-server-template/dbmodels"

// Constants used for serializations
const (
	JsonPrefix = ""
	JsonIndent = "  "
)

// Object that can be serialized and deserialized using JSON data type
type Serializable interface {
	SerializeJson() ([]byte, error)
	DeserializeJson(obj []byte) error
}

// Interface for switching between dbmodels and models
type Expandable interface {
	PopConstraints()
	Expand(obj *dbmodels.Object)
	Collapse() *dbmodels.Object
}
