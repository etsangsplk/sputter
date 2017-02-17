package api

import "fmt"

var (
	// True is literal value that represents any value other than False
	True = &Data{Value: true}

	// False is an alias for Nil or Nil
	False = Nil
)

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Sequence interfaces expose a one dimensional set of Values
type Sequence interface {
	Iterate() Iterator
}

// Finite interfaces allow a Sequence item to be retrieved by index
type Finite interface {
	Count() int
	Get(index int) Value
}

// Iterator interfaces are stateful iteration interfaces
type Iterator interface {
	Next() (Value, bool)
	Slice() Sequence
}

// Truthy evaluates whether or not a Value is Truthy
func Truthy(v Value) bool {
	switch {
	case v == Nil || v == False || v == nil || v == false:
		return false
	default:
		return true
	}
}

// CountSequence will either use Finite.Count() or iterate over the Sequence
func CountSequence(s Sequence) int {
	if f, ok := s.(Finite); ok {
		return f.Count()
	}

	i := s.Iterate()
	var r = 0
	for _, ok := i.Next(); ok; _, ok = i.Next() {
		r++
	}
	return r
}

// ValueToString either calls the String() method or tries to convert
func ValueToString(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return v.(string)
}
