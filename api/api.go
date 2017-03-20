package api

import "fmt"

var (
	// True represents the boolean value of True
	True = &Atom{Label: "true"}

	// False represents the boolean value of false
	False = &Atom{Label: "false"}

	// Nil is a value that represents the absence of a Value
	Nil Value
)

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Documented is implemented if a Value is Documented
type Documented interface {
	Docstring() string
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	switch {
	case v == Nil || v == False || v == false:
		return false
	default:
		return true
	}
}

// String either calls the String() method or tries to convert
func String(v Value) string {
	if v == nil {
		return "nil"
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
