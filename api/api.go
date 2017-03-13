package api

import (
	"fmt"
	"math/big"
)

const (
	// ExpectedCountable is thrown if taking count of a non-countable sequence
	ExpectedCountable = "sequence is not countable"

	// ExpectedSequence is thrown when a Value is not a Sequence
	ExpectedSequence = "value is not a list or vector"

	// ExpectedNumeric is thrown when a Value is not a Number
	ExpectedNumeric = "value is not numeric"
)

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

// Count will return the Count from a Countable Sequence or explode
func Count(s Sequence) int {
	if f, ok := s.(Countable); ok {
		return f.Count()
	}
	panic(ExpectedCountable)
}

// String either calls the String() method or tries to convert
func String(v Value) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return `"` + v.(string) + `"`
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
