package assert

import (
	"strings"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

// Any is the friendly name for a generic interface
type Any interface{}

// Wrapper wraps the testify assertions module in order to perform
// checking and conversion that is Sputter-specific
type Wrapper struct {
	as *assert.Assertions
}

// New instantiates a new Wrapper instance from the specified test
func New(t *testing.T) *Wrapper {
	return &Wrapper{
		as: assert.New(t),
	}
}

// String tests a Value for string equality
func (w *Wrapper) String(expect string, expr Any) {
	if s, ok := expr.(a.Str); ok {
		w.as.Equal(expect, string(s))
		return
	}
	if s, ok := expr.(a.Value); ok {
		w.as.Equal(expect, string(s.Str()))
		return
	}
	w.as.Equal(expect, expr)
}

// Float tests a Value for floating point equality
func (w *Wrapper) Float(expect float64, expr Any) {
	if f, ok := expr.(float64); ok {
		w.as.Equal(expect, f)
		return
	}
	if i, ok := expr.(int); ok {
		w.as.Equal(expect, float64(i))
		return
	}
	if n, ok := expr.(*a.Number); ok {
		w.as.Equal(a.EqualTo, a.NewFloat(expect).Cmp(n))
		return
	}
	w.as.Equal(expect, expr)
}

// Equal tests a Value for some kind of equality. Performs checks to do so
func (w *Wrapper) Equal(expect Any, expr Any) {
	if s, ok := expect.(string); ok {
		w.String(s, expr)
		return
	}
	if f, ok := expect.(float64); ok {
		w.Float(f, expr)
		return
	}
	if i, ok := expect.(int); ok {
		w.Float(float64(i), expr)
		return
	}
	if n, ok := expect.(*a.Number); ok {
		w.as.Equal(a.EqualTo, n.Cmp(expr.(*a.Number)))
		return
	}
	if s, ok := expect.(a.Str); ok {
		w.String(string(s), expr)
		return
	}
	if v, ok := expect.(a.Value); ok {
		w.String(string(v.Str()), expr)
		return
	}
	w.as.Equal(expect, expr)
}

// True tests a Value for boolean true
func (w *Wrapper) True(expr Any) {
	if b, ok := expr.(a.Bool); ok {
		w.as.True(bool(b))
		return
	}
	w.as.True(expr.(bool))
}

// Truthy tests a Value for Sputter-specific Truthy
func (w *Wrapper) Truthy(expr a.Value) {
	w.as.True(a.Truthy(expr))
}

// False tests a Value for boolean false
func (w *Wrapper) False(expr Any) {
	if b, ok := expr.(a.Bool); ok {
		w.as.False(bool(b))
		return
	}
	w.as.False(expr.(bool))
}

// Falsey tests a Value for Sputter-specific Falsey
func (w *Wrapper) Falsey(expr a.Value) {
	w.as.False(a.Truthy(expr))
}

// Contains check if the expected string is in the provided Value
func (w *Wrapper) Contains(expect string, expr a.Value) {
	var val string
	if s, ok := expr.(a.Str); ok {
		val = string(s)
	} else {
		val = string(expr.Str())
	}
	w.as.True(strings.Contains(val, expect))
}

// NotContains checks if the expected string is not in the provided Value
func (w *Wrapper) NotContains(expect string, expr a.Value) {
	var val string
	if s, ok := expr.(a.Str); ok {
		val = string(s)
	} else {
		val = string(s.Str())
	}
	w.as.False(strings.Contains(val, expect))
}

// Fail triggers an immediate test failure
func (w *Wrapper) Fail(msg string) {
	w.as.Fail(msg)
}

// Identical tests that two Values are referentially identical
func (w *Wrapper) Identical(expect Any, expr Any) {
	w.as.Equal(expect, expr)
}

// NotIdentical tests that two Values are not referentially identical
func (w *Wrapper) NotIdentical(expect Any, expr Any) {
	w.as.NotEqual(expect, expr)
}

// Nil tests if a Value is nil
func (w *Wrapper) Nil(expr Any) {
	if expr != nil && expr != a.Nil {
		w.Fail("value should be nil")
	}
}

// NotNil tests if a Value is not nil
func (w *Wrapper) NotNil(expr Any) {
	if expr == nil || expr == a.Nil {
		w.Fail("value should not be nil")
	}
}
