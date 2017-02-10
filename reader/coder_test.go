package interpreter_test

import (
	"fmt"
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func TestCreateCoder(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer("99")
	c := r.NewCoder(s.NewContext(), l)
	a.NotNil(c)
}

func TestCodeInteger(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer("99")
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()
	f, ok := v.(*big.Float)
	a.True(ok)
	a.Equal(0, f.Cmp(big.NewFloat(99)))
}

func TestCodeList(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`(99 "hello" 55.12)`)
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()
	list, ok := v.(*s.List)
	a.True(ok)

	iter := list.Iterate()
	value, ok := iter.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value.(*big.Float)))

	value, ok = iter.Next()
	a.True(ok)
	a.Equal("hello", value)

	value, ok = iter.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(55.12).Cmp(value.(*big.Float)))

	value, ok = iter.Next()
	a.False(ok)
}

func TestCodeNestedList(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()
	list, ok := v.(*s.List)
	a.True(ok)

	iter1 := list.Iterate()
	value, ok := iter1.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value.(*big.Float)))

	// get nested list
	value, ok = iter1.Next()
	a.True(ok)
	list2, ok := value.(*s.List)
	a.True(ok)

	// iterate over the rest of top-level list
	value, ok = iter1.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(55.12).Cmp(value.(*big.Float)))

	value, ok = iter1.Next()
	a.False(ok)

	// iterate over the nested list
	iter2 := list2.Iterate()
	value, ok = iter2.Next()
	a.True(ok)
	a.Equal("hello", value)

	value, ok = iter2.Next()
	a.True(ok)
	a.Equal("there", value)

	value, ok = iter2.Next()
	a.False(ok)
}

func TestUnclosedList(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(rec, r.ListNotClosed, "unclosed list")
			return
		}
		a.Fail("unclosed list didn't panic")
	}()

	l := r.NewLexer(`(99 ("hello" "there") 55.12`)
	c := r.NewCoder(s.NewContext(), l)
	c.Next()
}

func TestLiteral(t *testing.T) {
	a := assert.New(t)

	l := r.NewLexer(`'99`)
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()

	literal, ok := v.(*s.Data)
	a.True(ok)

	value, ok := literal.Value.(*big.Float)
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value))
}

func testCodeWithContext(a *assert.Assertions, code string, expect s.Value, context *s.Context) {
	l := r.NewLexer(code)
	c := r.NewCoder(s.NewContext(), l)
	a.Equal(expect, r.EvaluateCoder(context, c), code)
}

func evaluateToString(c *s.Context, v s.Value) string {
	result := s.Evaluate(c, v)
	if str, ok := result.(fmt.Stringer); ok {
		return str.String()
	}
	return result.(string)
}

func TestEvaluable(t *testing.T) {
	a := assert.New(t)
	c := s.NewContext()

	hello := &s.Function{
		Name: "hello",
		Exec: func(c *s.Context, args s.Iterable) s.Value {
			iter := args.Iterate()
			arg, _ := iter.Next()
			value := evaluateToString(c, arg)
			return "Hello, " + value + "!"
		},
	}

	c.Put("hello", hello)
	c.Put("name", "Bob")

	testCodeWithContext(a, `(hello "World")`, "Hello, World!", c)
	testCodeWithContext(a, `(hello name)`, "Hello, Bob!", c)
}
