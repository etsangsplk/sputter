package api

import "fmt"

var (
	// True represents the boolean value of True
	True = Bool(true)

	// False represents the boolean value of false
	False = Bool(false)

	// Nil is a value that represents the absence of a Value
	Nil = Atom("nil")
)

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	switch {
	case v == False || v == Nil:
		return false
	default:
		return true
	}
}

// Err generates a standard interpreter error
func Err(s string, args ...interface{}) string {
	sargs := make([]interface{}, len(args))
	for i, a := range args {
		if v, ok := a.(Value); ok {
			sargs[i] = string(v.Str())
		} else {
			sargs[i] = a
		}
	}
	return fmt.Sprintf(s, sargs...)
}
