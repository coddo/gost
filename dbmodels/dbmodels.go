package dbmodels

// Interface for defining the models as Objects which
// Object if it can be compared with other objects
type Object interface {
    Equal(obj Object) bool
}
