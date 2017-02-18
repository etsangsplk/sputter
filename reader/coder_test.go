package interpreter_test

import (
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
	list, ok := v.(*s.Cons)
	a.True(ok)

	i := list.Iterate()
	val, ok := i.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(val.(*big.Float)))

	val, ok = i.Next()
	a.True(ok)
	a.Equal("hello", val)

	val, ok = i.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(55.12).Cmp(val.(*big.Float)))

	val, ok = i.Next()
	a.False(ok)
}

func TestCodeVector(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`[99 "hello" 55.12]`)
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()
	vector, ok := v.(s.Vector)
	a.True(ok)
	a.Equal(0, big.NewFloat(99.0).Cmp(vector.Get(0).(*big.Float)))
	a.Equal("hello", vector.Get(1))
	a.Equal(0, big.NewFloat(55.120).Cmp(vector.Get(2).(*big.Float)))
}

func TestCodeNestedList(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()
	list, ok := v.(*s.Cons)
	a.True(ok)

	i1 := list.Iterate()
	val, ok := i1.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(val.(*big.Float)))

	// get nested list
	val, ok = i1.Next()
	a.True(ok)
	list2, ok := val.(*s.Cons)
	a.True(ok)

	// iterate over the rest of top-level list
	val, ok = i1.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(55.12).Cmp(val.(*big.Float)))

	val, ok = i1.Next()
	a.False(ok)

	// iterate over the nested list
	i2 := list2.Iterate()
	val, ok = i2.Next()
	a.True(ok)
	a.Equal("hello", val)

	val, ok = i2.Next()
	a.True(ok)
	a.Equal("there", val)

	val, ok = i2.Next()
	a.False(ok)
}

func TestUnclosedList(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(r.ListNotClosed, rec, "unclosed list")
			return
		}
		a.Fail("unclosed list didn't panic")
	}()

	l := r.NewLexer(`(99 ("hello" "there") 55.12`)
	c := r.NewCoder(s.NewContext(), l)
	c.Next()
}

func TestData(t *testing.T) {
	a := assert.New(t)

	l := r.NewLexer(`'99`)
	c := r.NewCoder(s.NewContext(), l)
	v := c.Next()

	d, ok := v.(*s.Data)
	a.True(ok)

	value, ok := d.Value.(*big.Float)
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value))
}

func testCodeWithContext(a *assert.Assertions, code string, expect s.Value, context *s.Context) {
	l := r.NewLexer(code)
	c := r.NewCoder(s.NewContext(), l)
	a.Equal(expect, r.EvalCoder(context, c), code)
}

func evaluateToString(c *s.Context, v s.Value) string {
	return s.String(s.Eval(c, v))
}

func TestEvaluable(t *testing.T) {
	a := assert.New(t)
	c := s.NewContext()

	hello := &s.Function{
		Name: "hello",
		Exec: func(c *s.Context, args s.Sequence) s.Value {
			i := args.Iterate()
			arg, _ := i.Next()
			v := evaluateToString(c, arg)
			return "Hello, " + v + "!"
		},
	}

	c.Put("hello", hello)
	c.Put("name", "Bob")

	testCodeWithContext(a, `(hello "World")`, "Hello, World!", c)
	testCodeWithContext(a, `(hello name)`, "Hello, Bob!", c)
}

func TestCoderBuiltins(t *testing.T) {
	a := assert.New(t)

	b := s.NewContext()
	b.PutFunction(&s.Function{
		Name: "hello",
		Exec: func(c *s.Context, args s.Sequence) s.Value {
			return "there"
		},
	})

	l := r.NewLexer(`(hello)`)
	c := r.NewCoder(b, l)
	a.Equal("there", r.EvalCoder(s.NewContext(), c), "builtin called")
}

func testCoderError(t *testing.T, src string, err string) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(err, rec, "coder errors out")
			return
		}
		a.Fail("coder doesn't error out like it should")
	}()

	l := r.NewLexer(src)
	c := r.NewCoder(s.NewContext(), l)
	for n := c.Next(); n != r.EndOfCoder; n = c.Next() {
	}
}

func TestCoderErrors(t *testing.T) {
	testCoderError(t, "(99 100 ", r.ListNotClosed)
	testCoderError(t, "[99 100 ", r.VectorNotClosed)

	testCoderError(t, "99 100)", r.UnmatchedListEnd)
	testCoderError(t, "99 100]", r.UnmatchedVectorEnd)
}
