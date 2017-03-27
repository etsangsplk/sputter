package api

import "fmt"

var (
	// True represents the boolean value of True
	True Value = true

	// False represents the boolean value of false
	False Value = false

	// Nil is f value that represents the absence of f Value
	Nil Value
)

const (
	// MetaName is the Metadata key for a Value's Name
	MetaName = Name("name")

	// MetaDoc is the Metadata key for Documentation Strings
	MetaDoc = Name("docstring")
)

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Annotated is implemented if a Value is Annotated with Metadata
type Annotated interface {
	Metadata() Variables
	WithMetadata(md Variables) Annotated
}

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	switch {
	case v == False || v == Nil:
		return false
	default:
		return true
	}
}

// String either calls the String() method or tries to convert
func String(v Value) string {
	if v == Nil {
		return "nil"
	}
	if v == True {
		return "true"
	}
	if v == False {
		return "false"
	}

	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	if n, ok := v.(Name); ok {
		return string(n)
	}
	if s, ok := v.(string); ok {
		return fmt.Sprintf("%q", s)
	}
	return fmt.Sprintf("(<anon> :instance %p)", &v)
}
