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

// Value is the generic interface for all 'Values' in the VM
type Value interface {
}

// Sequence interfaces expose a one dimensional set of Values
type Sequence interface {
	Iterate() Iterator
	Count() int	
}

// Iterator interfaces are stateful iteration interfaces
type Iterator interface {
	Next() (Value, bool)
	Iterable() Sequence
}

// Indexable interfaces can return an item by index
type Indexable interface {
	Get(index int) Value
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

// ValueToString either calls the String() method or tries to convert
func ValueToString(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return v.(string)
}
