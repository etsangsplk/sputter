package reader_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func TestCreateReader(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer("99")
	c := s.NewContext()
	tr := r.NewReader(c, l)
	a.NotNil(tr)
}

func TestReadInteger(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer("99")
	c := s.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	f, ok := v.(*big.Float)
	a.True(ok)
	a.Equal(0, f.Cmp(big.NewFloat(99)))
}

func TestReadList(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`(99 "hello" 55.12)`)
	c := s.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
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

func TestReadVector(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`[99 "hello" 55.12]`)
	c := s.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	vector, ok := v.(s.Vector)
	a.True(ok)
	a.Equal(0, big.NewFloat(99.0).Cmp(vector.Get(0).(*big.Float)))
	a.Equal("hello", vector.Get(1))
	a.Equal(0, big.NewFloat(55.120).Cmp(vector.Get(2).(*big.Float)))
}

func TestReadNestedList(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := s.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
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
	c := s.NewContext()
	tr := r.NewReader(c, l)
	tr.Next()
}

func TestSimpleData(t *testing.T) {
	a := assert.New(t)

	l := r.NewLexer(`'99`)
	c := s.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()

	d, ok := v.(*s.Quote)
	a.True(ok)

	value, ok := d.Value.(*big.Float)
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value))
}

func TestListData(t *testing.T) {
	a := assert.New(t)

	l := r.NewLexer(`'(symbol true)`)
	c := s.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()

	d, ok := v.(*s.Quote)
	a.True(ok)

	value, ok := d.Value.(*s.Cons)
	a.True(ok)

	if sym, ok := value.Car.(*s.Symbol); ok {
		a.Equal("symbol", string(sym.Name), "symbol was literal")
	} else {
		a.Fail("first element should be symbol")
	}

	if n, ok := value.Cdr.(*s.Cons); ok {
		b, ok := n.Car.(*s.Atom)
		a.True(ok)
		a.Equal(s.True, b)

		nl, ok := n.Cdr.(*s.Cons)
		a.True(ok)
		a.Equal(s.Nil, nl, "list properly terminated")
	} else {
		a.Fail("rest() elements not a cons")
	}
}

func testCodeWithContext(a *assert.Assertions, code string, expect s.Value, context s.Context) {
	l := r.NewLexer(code)
	c := s.NewContext()
	tr := r.NewReader(c, l)
	a.Equal(expect, r.EvalReader(context, tr), code)
}

func evaluateToString(c s.Context, v s.Value) string {
	return s.String(s.Eval(c, v))
}

func TestEvaluable(t *testing.T) {
	a := assert.New(t)
	c := s.NewContext()

	hello := &s.Function{
		Name: "hello",
		Exec: func(c s.Context, args s.Sequence) s.Value {
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

func TestBuiltIns(t *testing.T) {
	a := assert.New(t)

	b := s.NewContext()
	s.PutFunction(b, &s.Function{
		Name: "hello",
		Exec: func(c s.Context, args s.Sequence) s.Value {
			return "there"
		},
	})

	l := r.NewLexer(`(hello)`)
	tr := r.NewReader(b, l)
	c := s.ChildContext(b)
	a.Equal("there", r.EvalReader(c, tr), "builtin called")
}

func TestReaderPrepare(t *testing.T) {
	a := assert.New(t)

	b := s.NewContext()
	s.PutFunction(b, &s.Function{
		Name: "hello",
		Prepare: func(c s.Context, l s.Sequence) s.Value {
			if _, ok := l.(*s.Cons); !ok {
				a.Fail("provided list is not a cons")
			}
			return s.Vector{"you"}
		},
	})

	l := r.NewLexer(`(hello)`)
	tr := r.NewReader(b, l)
	v := tr.Next()

	if rv, ok := v.(s.Vector); ok {
		a.Equal("you", rv.Get(0), "prepared transformed into vector")
	} else {
		a.Fail("prepare did not transform")
	}
}

func testReaderError(t *testing.T, src string, err string) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(err, rec, "coder errors out")
			return
		}
		a.Fail("coder doesn't error out like it should")
	}()

	c := s.NewContext()
	l := r.NewLexer(src)
	tr := r.NewReader(c, l)
	r.EvalReader(c, tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", r.ListNotClosed)
	testReaderError(t, "[99 100 ", r.VectorNotClosed)

	testReaderError(t, "99 100)", r.UnmatchedListEnd)
	testReaderError(t, "99 100]", r.UnmatchedVectorEnd)
}
