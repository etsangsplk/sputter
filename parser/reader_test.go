package parser_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	p "github.com/kode4food/sputter/parser"
	"github.com/stretchr/testify/assert"
)

func TestCreateReader(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer("99")
	c := a.NewContext()
	tr := p.NewReader(c, l)
	as.NotNil(tr)
}

func TestReadInteger(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer("99")
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()
	f, ok := v.(*a.Number)
	as.True(ok)
	as.Equal(a.EqualTo, f.Cmp(a.NewFloat(99)))
}

func TestReadList(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`(99 "hello" 55.12)`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()
	list, ok := v.(*a.List)
	as.True(ok)

	i := a.Iterate(list)
	val, ok := i.Next()
	as.True(ok)
	as.Equal(a.EqualTo, a.NewFloat(99).Cmp(val.(*a.Number)))

	val, ok = i.Next()
	as.True(ok)
	as.Equal("hello", val)

	val, ok = i.Next()
	as.True(ok)
	as.Equal(a.EqualTo, a.NewFloat(55.12).Cmp(val.(*a.Number)))

	val, ok = i.Next()
	as.False(ok)
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`[99 "hello" 55.12]`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()
	vector, ok := v.(a.Vector)
	as.True(ok)

	res, ok := vector.Get(0)
	as.True(ok, "get by index")
	as.Equal(a.EqualTo, a.NewFloat(99.0).Cmp(res.(*a.Number)))

	res, ok = vector.Get(1)
	as.True(ok, "get by index")
	as.Equal("hello", res)

	res, ok = vector.Get(2)
	as.True(ok, "get by index")
	as.Equal(a.EqualTo, a.NewFloat(55.120).Cmp(res.(*a.Number)))
}

func TestReadMap(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`{:name "blah" :age 99}`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()
	m, ok := v.(a.Associative)
	as.True(ok)
	as.Equal(2, m.Count(), "map count is correct")
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()
	list, ok := v.(*a.List)
	as.True(ok)

	i1 := a.Iterate(list)
	val, ok := i1.Next()
	as.True(ok)
	as.Equal(a.EqualTo, a.NewFloat(99).Cmp(val.(*a.Number)))

	// get nested list
	val, ok = i1.Next()
	as.True(ok)
	list2, ok := val.(*a.List)
	as.True(ok)

	// iterate over the rest of top-level list
	val, ok = i1.Next()
	as.True(ok)
	as.Equal(a.EqualTo, a.NewFloat(55.12).Cmp(val.(*a.Number)))

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

	l := p.NewLexer(`'99`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()

	d, ok := v.(a.Quoted)
	as.True(ok)

	value, ok := a.Eval(c, d).(*a.Number)
	as.True(ok)
	as.Equal(a.EqualTo, a.NewFloat(99).Cmp(value))
}

func TestListData(t *testing.T) {
	as := assert.New(t)

	l := p.NewLexer(`'(symbol true)`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.Next()

	d, ok := v.(a.Quoted)
	as.True(ok)

	value, ok := d.Eval(c).(*a.List)
	as.True(ok)

	if sym, ok := value.First().(a.Symbol); ok {
		as.Equal("symbol", string(sym.Name()), "symbol was literal")
	} else {
		as.Fail("first element should be symbol")
	}

	if n, ok := value.Rest().(*a.List); ok {
		b, ok := n.First().(a.Value)
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

func testCodeWithContext(
	as *assert.Assertions, code string, expect a.Value, c a.Context) {
	l := p.NewLexer(code)
	tr := p.NewReader(a.NewContext(), l)
	as.Equal(expect, p.EvalReader(c, tr), code)
}

func TestEvaluable(t *testing.T) {
	as := assert.New(t)
	c := a.NewContext()

	hello := a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		i := a.Iterate(args)
		arg, _ := i.Next()
		v := a.Eval(c, arg)
		return "Hello, " + v.(string) + "!"
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	})

	c.Put("hello", hello)
	c.Put("name", "Bob")

	testCodeWithContext(as, `(hello "World")`, "Hello, World!", c)
	testCodeWithContext(as, `(hello name)`, "Hello, Bob!", c)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	b := a.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Put("hello", a.NewFunction(func(_ a.Context, _ a.Sequence) a.Value {
		return "there"
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	}))

	l := p.NewLexer(`(hello)`)
	tr := p.NewReader(b, l)
	c := a.ChildContext(b)
	as.Equal("there", p.EvalReader(c, tr), "builtin called")
}

func TestReaderPrepare(t *testing.T) {
	as := assert.New(t)

	b := a.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Delete("hello")
	ns.Put("hello", a.NewMacro(func(_ a.Context, l a.Sequence) a.Value {
		if _, ok := l.(*a.List); !ok {
			as.Fail("provided list is not a cons")
		}
		return a.Vector{"you"}
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	}))

	l := p.NewLexer(`(hello)`)
	tr := p.NewReader(b, l)
	v := tr.Next()

	if rv, ok := v.(a.Vector); ok {
		v1, ok := rv.Get(0)
		as.True(ok, "prepared transformed into vector")
		as.Equal("you", v1, "prepared transformed into vector")
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
	l := p.NewLexer(src)
	tr := p.NewReader(c, l)
	p.EvalReader(c, tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", p.ListNotClosed)
	testReaderError(t, "[99 100 ", p.VectorNotClosed)
	testReaderError(t, "{:key 99", p.MapNotClosed)

	testReaderError(t, "99 100)", p.UnmatchedListEnd)
	testReaderError(t, "99 100]", p.UnmatchedVectorEnd)
	testReaderError(t, "99}", p.UnmatchedMapEnd)
	testReaderError(t, "{99}", p.MapNotPaired)
}
