// Package api provides all of the interfaces, primitives and data
// structures that comprise the core of the Sputter interpreter.
package api

import "fmt"

var (
	// True represents the boolean value of True
	True Value = true

	// False represents the boolean value of false
	False Value = false

	// Nil is a value that represents the absence of a Value
	Nil Value
)

var Native = Name("<native>")

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
	switch {
	case v == Nil:
		return "nil"
	case v == True:
		return "true"
	case v == False:
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
	return stringDump(v)
}

func stringDump(v Value) string {
	md := Metadata{}
	if n, ok := v.(Named); ok {
		md = md.Merge(Metadata{MetaName: n.Name()})
	}
	if t, ok := v.(Typed); ok {
		md = md.Merge(Metadata{MetaType: t.Type()})
	} else {
		md = md.Merge(Metadata{MetaType: Native})
	}
	ptr := fmt.Sprintf("%p", &v)
	md = md.Merge(Metadata{MetaInstance: ptr})
	if a, ok := v.(Annotated); ok {
		md = md.Merge(Metadata{MetaMeta: a.Metadata()})
	}
	return md.String()
}

// Err generates a standard interpreter error
func Err(s string, args ...interface{}) string {
	return fmt.Sprintf(s, args...)
}
