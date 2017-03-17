package api

import "fmt"

var (
	// True represents the boolean value of True
	True = &Atom{Label: "true"}

	// False represents the boolean value of false
	False = &Atom{Label: "false"}

	// Nil is a value that represents the absence of a Value
	Nil = &Atom{Label: "nil"}
)

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	switch {
	case v == Nil || v == False || v == nil || v == false:
		return false
	default:
		return true
	}
}

// String either calls the String() method or tries to convert
func String(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return `"` + v.(string) + `"`
}
