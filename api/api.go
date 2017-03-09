package api

import (
	"fmt"
	"math/big"
)

const (
	// ExpectedFinite is thrown if when taking count of a non-finite sequence
	ExpectedFinite = "sequence is not finite and can't be counted"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a list or vector"

	// ExpectedNumeric is thrown when a Value is not a Number
	ExpectedNumeric = "value is not numeric"
)

var (
	// True is a value that represents any value other than False
	True = &Atom{Label: "true"}

	// False is a value that represents either itself or nil
	False = &Atom{Label: "false"}
)

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Sequence interfaces expose a one dimensional set of Values
type Sequence interface {
	Iterate() Iterator
}

// Finite interfaces allow a Sequence item to be retrieved by index
type Finite interface {
	Count() int
	Get(index int) Value
}

// Mapped interfaces allow a Sequence item to be retrieved by Name
type Mapped interface {
	Count() int
	Get(key Value) Value
}

// Iterator interfaces are stateful iteration interfaces
type Iterator interface {
	Next() (Value, bool)
	Rest() Sequence
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

// Count will either use Finite.Count() or iterate over the Sequence
func Count(s Sequence) int {
	if f, ok := s.(Finite); ok {
		return f.Count()
	}
	panic(ExpectedFinite)
}

// String either calls the String() method or tries to convert
func String(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return v.(string)
}

// AssertSequence will cast a Value into a Sequence or explode violently
func AssertSequence(v Value) Sequence {
	if r, ok := v.(Sequence); ok {
		return r
	}
	panic(ExpectedSequence)
}

// AssertNumeric will cast a Value into a Numeric or explode violently
func AssertNumeric(v Value) *big.Float {
	if r, ok := v.(*big.Float); ok {
		return r
	}
	panic(ExpectedNumeric)
}
