package dbmodels

// Objecter is an interface for defining a method for comparing two entities
type Objecter interface {
	Equal(obj Objecter) bool
}
