package reader_test

import (
	"fmt"
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func TestCreateReader(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer("99")
	c := a.NewContext()
	tr := r.NewReader(c, l)
	as.NotNil(tr)
}

func TestReadInteger(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer("99")
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	f, ok := v.(*big.Float)
	as.True(ok)
	as.Equal(0, f.Cmp(big.NewFloat(99)))
}

func TestReadList(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(`(99 "hello" 55.12)`)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	list, ok := v.(*a.List)
	as.True(ok)

	i := a.Iterate(list)
	val, ok := i.Next()
	as.True(ok)
	as.Equal(0, big.NewFloat(99).Cmp(val.(*big.Float)))

	val, ok = i.Next()
	as.True(ok)
	as.Equal("hello", val)

	val, ok = i.Next()
	as.True(ok)
	as.Equal(0, big.NewFloat(55.12).Cmp(val.(*big.Float)))

	val, ok = i.Next()
	as.False(ok)
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(`[99 "hello" 55.12]`)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	vector, ok := v.(a.Vector)
	as.True(ok)
	as.Equal(0, big.NewFloat(99.0).Cmp(vector.Get(0).(*big.Float)))
	as.Equal("hello", vector.Get(1))
	as.Equal(0, big.NewFloat(55.120).Cmp(vector.Get(2).(*big.Float)))
}

func TestReadMap(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(`{:name "blah" :age 99}`)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	m, ok := v.(a.ArrayMap)
	as.True(ok)
	as.Equal(2, m.Count(), "map count is correct")
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()
	list, ok := v.(*a.List)
	as.True(ok)

	i1 := a.Iterate(list)
	val, ok := i1.Next()
	as.True(ok)
	as.Equal(0, big.NewFloat(99).Cmp(val.(*big.Float)))

	// get nested list
	val, ok = i1.Next()
	as.True(ok)
	list2, ok := val.(*a.List)
	as.True(ok)

	// iterate over the rest of top-level list
	val, ok = i1.Next()
	as.True(ok)
	as.Equal(0, big.NewFloat(55.12).Cmp(val.(*big.Float)))

	val, ok = i1.Next()
	as.False(ok)

	// iterate over the nested list
	i2 := a.Iterate(list2)
	val, ok = i2.Next()
	as.True(ok)
	as.Equal("hello", val)

	val, ok = i2.Next()
	as.True(ok)
	as.Equal("there", val)

	val, ok = i2.Next()
	as.False(ok)
}

func TestSimpleData(t *testing.T) {
	as := assert.New(t)

	l := r.NewLexer(`'99`)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()

	d, ok := v.(*a.Quote)
	as.True(ok)

	value, ok := d.Value.(*big.Float)
	as.True(ok)
	as.Equal(0, big.NewFloat(99).Cmp(value))
}

func TestListData(t *testing.T) {
	as := assert.New(t)

	l := r.NewLexer(`'(symbol true)`)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	v := tr.Next()

	d, ok := v.(*a.Quote)
	as.True(ok)

	value, ok := d.Value.(*a.List)
	as.True(ok)

	if sym, ok := value.First().(*a.Symbol); ok {
		as.Equal("symbol", string(sym.Name), "symbol was literal")
	} else {
		as.Fail("first element should be symbol")
	}

	if n, ok := value.Rest().(*a.List); ok {
		b, ok := n.First().(*a.Atom)
		as.True(ok)
		as.Equal(a.True, b)

		nl, ok := n.Rest().(*a.List)
		as.True(ok)
		as.Equal(a.EmptyList, nl, "list properly terminated")
		as.False(nl.IsSequence(), "list properly terminated")
	} else {
		as.Fail("rest() elements not a list")
	}
}

func testCodeWithContext(as *assert.Assertions, code string, expect a.Value, context a.Context) {
	l := r.NewLexer(code)
	c := a.NewContext()
	tr := r.NewReader(c, l)
	as.Equal(expect, r.EvalReader(context, tr), code)
}

func evaluateToString(c a.Context, v a.Value) string {
	return fmt.Sprint(a.Eval(c, v))
}

func TestEvaluable(t *testing.T) {
	as := assert.New(t)
	c := a.NewContext()

	hello := &a.Function{
		Name: "hello",
		Apply: func(c a.Context, args a.Sequence) a.Value {
			i := a.Iterate(args)
			arg, _ := i.Next()
			v := evaluateToString(c, arg)
			return "Hello, " + v + "!"
		},
	}

	c.Put("hello", hello)
	c.Put("name", "Bob")

	testCodeWithContext(as, `(hello "World")`, "Hello, World!", c)
	testCodeWithContext(as, `(hello name)`, "Hello, Bob!", c)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	b := a.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Put("hello", &a.Function{
		Name: "hello",
		Apply: func(c a.Context, args a.Sequence) a.Value {
			return "there"
		},
	})

	l := r.NewLexer(`(hello)`)
	tr := r.NewReader(b, l)
	c := a.ChildContext(b)
	as.Equal("there", r.EvalReader(c, tr), "builtin called")
}

func TestReaderPrepare(t *testing.T) {
	as := assert.New(t)

	b := a.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Delete("hello")
	ns.Put("hello", &a.Function{
		Name: "hello",
		Prepare: func(c a.Context, l a.Sequence) a.Value {
			if _, ok := l.(*a.List); !ok {
				as.Fail("provided list is not a cons")
			}
			return a.Vector{"you"}
		},
	})

	l := r.NewLexer(`(hello)`)
	tr := r.NewReader(b, l)
	v := tr.Next()

	if rv, ok := v.(a.Vector); ok {
		as.Equal("you", rv.Get(0), "prepared transformed into vector")
	} else {
		as.Fail("prepare did not transform")
	}
}

func testReaderError(t *testing.T, src string, err string) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.Equal(err, rec, "coder errors out")
			return
		}
		as.Fail("coder doesn't error out like it should")
	}()

	c := a.NewContext()
	l := r.NewLexer(src)
	tr := r.NewReader(c, l)
	r.EvalReader(c, tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", r.ListNotClosed)
	testReaderError(t, "[99 100 ", r.VectorNotClosed)
	testReaderError(t, "{:key 99", r.MapNotClosed)

	testReaderError(t, "99 100)", r.UnmatchedListEnd)
	testReaderError(t, "99 100]", r.UnmatchedVectorEnd)
	testReaderError(t, "99}", r.UnmatchedMapEnd)
	testReaderError(t, "{99}", r.MapNotPaired)
}
