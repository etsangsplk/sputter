package assert

import (
	"strings"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

// InvalidTestExpression is thrown if an unsupported expression type is used
const InvalidTestExpression = "invalid test expression: %s"

type (
	// Any is the friendly name for a generic interface
	Any interface{}

	// Wrapper wraps the testify assertions module in order to perform
	// checking and conversion that is Sputter-specific
	Wrapper struct {
		as *assert.Assertions
	}
)

// New instantiates a new Wrapper instance from the specified test
func New(t *testing.T) *Wrapper {
	return &Wrapper{
		as: assert.New(t),
	}
}

// String tests a Value for string equality
func (w *Wrapper) String(expect string, expr Any) {
	if s, ok := expr.(string); ok {
		w.as.Equal(expect, s)
		return
	}
	if v, ok := expr.(a.Value); ok {
		w.as.Equal(expect, string(a.MakeStr(v)))
		return
	}
	panic(a.ErrStr(InvalidTestExpression, expr))
}

// Number tests a Value for numeric equality
func (w *Wrapper) Number(expect float64, expr Any) {
	if f, ok := expr.(float64); ok {
		w.as.Equal(expect, f)
		return
	}
	if i, ok := expr.(int); ok {
		w.as.Equal(expect, float64(i))
		return
	}
	if n, ok := expr.(a.Number); ok {
		w.as.Equal(a.EqualTo, a.NewFloat(expect).Cmp(n))
		return
	}
	panic(a.ErrStr(InvalidTestExpression, expr))
}

// Equal tests a Value for some kind of equality. Performs checks to do so
func (w *Wrapper) Equal(expect a.Value, expr Any) {
	if s, ok := expect.(a.Str); ok {
		w.String(string(s), expr)
		return
	}
	if n, ok := expect.(a.Number); ok {
		f, _ := n.Float64()
		w.Number(f, expr)
		//w.as.Equal(a.EqualTo, n.Cmp(expr.(a.Number)))
		return
	}
	w.String(string(expect.Str()), expr)
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
	val := string(a.MakeStr(expr))
	w.as.True(strings.Contains(val, expect))
}

// NotContains checks if the expected string is not in the provided Value
func (w *Wrapper) NotContains(expect string, expr a.Value) {
	val := string(a.MakeStr(expr))
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
	if expr == a.Nil {
		w.as.Nil(nil)
		return
	}
	w.as.Nil(expr)
}

// NotNil tests if a Value is not nil
func (w *Wrapper) NotNil(expr Any) {
	if expr == a.Nil {
		w.as.NotNil(nil)
		return
	}
	w.as.NotNil(expr)
}

// Compare tests if the Comparison of two Numbers is correct
func (w *Wrapper) Compare(c a.Comparison, l a.Number, r a.Number) {
	w.as.Equal(c, l.Cmp(r))
}

// ExpectError is used with a defer to make sure an error was triggered
func (w *Wrapper) ExpectError(err a.Object) {
	if rec := recover(); rec != nil {
		if a.IsErr(rec) {
			errStr := string(err.MustGet(a.MessageKey).Str())
			recStr := rec.(a.Object).MustGet(a.MessageKey).Str()
			w.String(errStr, recStr)
			return
		}
	}
	w.Fail("error not raised")
}
