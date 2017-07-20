package api

import "fmt"

var (
	// MessageKey identifies the message property of an error
	MessageKey = NewKeyword("message")

	errName = Name("error")

	errPrototype = Properties{
		TypeKey:    errName,
		MessageKey: Str("an unknown error occurred"),
	}
)

// Err generates a standard interpreter error
func Err(p Properties) Object {
	return errPrototype.Child(p)
}

// ErrStr generates a standard interpreter error from a Str
func ErrStr(s string, args ...interface{}) Object {
	sargs := make([]interface{}, len(args))
	for i, a := range args {
		if v, ok := a.(Value); ok {
			sargs[i] = string(v.Str())
		} else {
			sargs[i] = a
		}
	}
	return Err(Properties{
		MessageKey: Str(fmt.Sprintf(s, sargs...)),
	})
}

// IsErr tests whether or not a Value is an error Object
func IsErr(v interface{}) bool {
	if e, ok := v.(Object); ok {
		if t, ok := e.Get(TypeKey); ok && t == errName {
			return true
		}
	}
	return false
}
