package api

import "fmt"

const errName = Name("error")

type (
	// Error represents a Go error as an Object
	Error interface {
		error
		Object
	}

	err struct {
		Object
	}
)

var (
	// MessageKey identifies the message property of an error
	MessageKey = NewKeyword("message")

	errPrototype = Properties{
		TypeKey:    errName,
		MessageKey: Str("an unknown error occurred"),
	}
)

// Err generates a standard interpreter error
func Err(p Properties) Error {
	return &err{
		Object: errPrototype.Child(p),
	}
}

// ErrStr generates a standard interpreter error from a Str
func ErrStr(s string, args ...interface{}) Error {
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

func (e *err) Error() string {
	return string(e.MustGet(MessageKey).(Str))
}
