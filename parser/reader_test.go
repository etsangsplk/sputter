package parser_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/builtins"
	p "github.com/kode4food/sputter/parser"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) *a.Number {
	return a.NewFloat(f)
}

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
	v := tr.First()
	n, ok := v.(*a.Number)
	as.True(ok)
	as.Equal(f(99), n)
}

func TestReadList(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`(99 "hello" 55.12)`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.First()
	list, ok := v.(a.List)
	as.True(ok)

	i := a.Iterate(list)
	val, ok := i.Next()
	as.True(ok)
	as.Number(99, val)

	val, ok = i.Next()
	as.True(ok)
	as.String("hello", val)

	val, ok = i.Next()
	as.True(ok)
	as.Number(55.12, val)

	_, ok = i.Next()
	as.False(ok)
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`[99 "hello" 55.12]`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.First()
	vector, ok := v.(a.Vector)
	as.True(ok)

	res, ok := vector.ElementAt(0)
	as.True(ok)
	as.Number(99, res)

	res, ok = vector.ElementAt(1)
	as.True(ok)
	as.String("hello", res)

	res, ok = vector.ElementAt(2)
	as.True(ok)
	as.Number(55.120, res)
}

func TestReadMap(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`{:name "blah" :age 99}`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.First()
	m, ok := v.(a.Associative)
	as.True(ok)
	as.Number(2, m.Count())
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := tr.First()
	list, ok := v.(a.List)
	as.True(ok)

	i1 := a.Iterate(list)
	val, ok := i1.Next()
	as.True(ok)
	as.Number(99, val)

	// get nested list
	val, ok = i1.Next()
	as.True(ok)
	list2, ok := val.(a.List)
	as.True(ok)

	// iterate over the rest of top-level list
	val, ok = i1.Next()
	as.True(ok)
	as.Number(55.12, val)

	_, ok = i1.Next()
	as.False(ok)

	// iterate over the nested list
	i2 := a.Iterate(list2)
	val, ok = i2.Next()
	as.True(ok)
	as.String("hello", val)

	val, ok = i2.Next()
	as.True(ok)
	as.String("there", val)

	_, ok = i2.Next()
	as.False(ok)
}

func TestSimpleData(t *testing.T) {
	as := assert.New(t)

	l := p.NewLexer(`'99`)
	c := a.NewEvalContext()
	tr := p.NewReader(c, l)
	v1, ok := a.EvalSequence(c, tr).(*a.Number)
	as.True(ok)
	as.Number(99, v1)

	l = p.NewLexer(`'[1 2 3]`)
	c = a.NewEvalContext()
	tr = p.NewReader(c, l)
	v2, ok := a.EvalSequence(c, tr).(a.Vector)
	as.True(ok)
	as.String("[1 2 3]", v2)
	_, ok = v2.(a.Expression)
	as.False(ok)

	l = p.NewLexer(`'{:name "bob"}`)
	c = a.NewEvalContext()
	tr = p.NewReader(c, l)
	v3, ok := a.EvalSequence(c, tr).(a.Associative)
	as.True(ok)
	as.String(`{:name "bob"}`, v3)
	_, ok = v3.(a.Expression)
	as.False(ok)
}

func TestNestedData(t *testing.T) {
	as := assert.New(t)

	l := p.NewLexer(`''99`)
	c := a.NewEvalContext()
	tr := p.NewReader(c, l)
	v := a.EvalSequence(c, tr)
	as.String("(sputter:quote 99)", v)
}

func TestListData(t *testing.T) {
	as := assert.New(t)

	l := p.NewLexer(`'(symbol true)`)
	c := a.NewContext()
	tr := p.NewReader(c, l)
	v := a.EvalSequence(c, tr)
	value, ok := v.(a.List)
	as.True(ok)

	if sym, ok := value.First().(a.Symbol); ok {
		as.String("symbol", sym.Name())
	} else {
		as.Fail("first element should be symbol")
	}

	if n, ok := value.Rest().(a.List); ok {
		b, ok := n.First().(a.Value)
		as.True(ok)
		as.Equal(a.True, b)

		nl, ok := n.Rest().(a.List)
		as.True(ok)
		as.Equal(a.EmptyList, nl)
		as.False(nl.IsSequence())
	} else {
		as.Fail("rest() elements not a list")
	}
}

func testCodeWithContext(
	as *assert.Wrapper, code string, expect a.Value, c a.Context) {
	l := p.NewLexer(s(code))
	tr := p.NewReader(a.NewEvalContext(), l)
	as.Equal(expect, a.EvalSequence(c, tr))
}

func TestEvaluable(t *testing.T) {
	as := assert.New(t)
	c := a.NewContext()

	hello := a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		i := a.Iterate(args)
		arg, _ := i.Next()
		v := arg.Eval(c)
		return s("Hello, " + string(v.(a.Str)) + "!")
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	}).(a.Function)

	c.Put("hello", hello)
	c.Put("name", s("Bob"))

	testCodeWithContext(as, `(hello "World")`, s("Hello, World!"), c)
	testCodeWithContext(as, `(hello name)`, s("Hello, Bob!"), c)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	b := a.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Put("hello", a.NewFunction(func(_ a.Context, _ a.Sequence) a.Value {
		return s("there")
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	}).(a.Function))

	l := p.NewLexer(`(hello)`)
	tr := p.NewReader(b, l)
	c := a.ChildContext(b)
	as.String("there", a.EvalSequence(c, tr))
}

func TestReaderPrepare(t *testing.T) {
	as := assert.New(t)

	b := a.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Delete("hello")
	ns.Put("hello", a.NewMacro(func(_ a.Context, l a.Sequence) a.Value {
		if _, ok := l.(a.List); !ok {
			as.Fail("provided sequence is not a list")
		}
		return a.NewVector(s("you"))
	}).WithMetadata(a.Metadata{
		a.MetaName: a.Name("hello"),
	}).(a.Function))

	l := p.NewLexer(`(hello)`)
	tr := p.NewReader(b, l)
	v := a.EvalSequence(b, tr)

	if rv, ok := v.(a.Vector); ok {
		v1, ok := rv.ElementAt(0)
		as.True(ok)
		as.String("you", v1)
	} else {
		as.Fail("prepare did not transform")
	}
}

func testReaderError(t *testing.T, src string, err string) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.String(err, rec)
			return
		}
		as.Fail("coder doesn't error out like it should")
	}()

	c := a.NewContext()
	l := p.NewLexer(s(src))
	tr := p.NewReader(c, l)
	a.EvalSequence(c, tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", p.ListNotClosed)
	testReaderError(t, "[99 100 ", p.VectorNotClosed)
	testReaderError(t, "{:key 99", p.MapNotClosed)

	testReaderError(t, "99 100)", p.UnmatchedListEnd)
	testReaderError(t, "99 100]", p.UnmatchedVectorEnd)
	testReaderError(t, "99}", p.UnmatchedMapEnd)
	testReaderError(t, "{99}", p.MapNotPaired)

	testReaderError(t, "(", p.ListNotClosed)
	testReaderError(t, "'", p.QuoteNotPaired)
}
