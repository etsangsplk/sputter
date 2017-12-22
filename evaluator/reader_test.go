package evaluator_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}

func TestCreateReader(t *testing.T) {
	as := assert.New(t)
	l := e.Scan("99")
	c := a.Variables{}
	tr := e.Read(c, l)
	as.NotNil(tr)
}

func TestReadInteger(t *testing.T) {
	as := assert.New(t)
	l := e.Scan("99")
	c := a.Variables{}
	tr := e.Read(c, l)
	v := tr.First()
	n, ok := v.(a.Number)
	as.True(ok)
	as.Equal(f(99), n)
}

func TestReadList(t *testing.T) {
	as := assert.New(t)
	l := e.Scan(`(99 "hello" 55.12)`)
	c := a.Variables{}
	tr := e.Read(c, l)
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
	l := e.Scan(`[99 "hello" 55.12]`)
	c := a.Variables{}
	tr := e.Read(c, l)
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
	l := e.Scan(`{:name "blah" :age 99}`)
	c := a.Variables{}
	tr := e.Read(c, l)
	v := tr.First()
	m, ok := v.(a.Associative)
	as.True(ok)
	as.Number(2, m.Count())
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := e.Scan(`(99 ("hello" "there") 55.12)`)
	c := a.Variables{}
	tr := e.Read(c, l)
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

func testCodeWithContext(
	as *assert.Wrapper, code string, expect a.Value, c a.Context) {
	as.Equal(expect, e.EvalStr(c, s(code)))
}

func TestEvaluable(t *testing.T) {
	as := assert.New(t)

	hello := a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		i := a.Iterate(args)
		arg, _ := i.Next()
		v := a.Eval(c, arg)
		return s("Hello, " + string(v.(a.Str)) + "!")
	}).WithMetadata(a.Properties{
		a.NameKey: a.Name("hello"),
	}).(a.Function)

	c := e.NewEvalContext()
	c.Put("hello", hello)
	c.Put("name", s("Bob"))

	testCodeWithContext(as, `(hello "World")`, s("Hello, World!"), c)
	testCodeWithContext(as, `(hello name)`, s("Hello, Bob!"), c)
}

func TestBuiltIns(t *testing.T) {
	as := assert.New(t)

	b := e.NewEvalContext()
	ns := a.GetContextNamespace(b)
	ns.Delete("hello")

	ns.Put("hello", a.NewExecFunction(func(_ a.Context, _ a.Sequence) a.Value {
		return s("there")
	}).WithMetadata(a.Properties{
		a.NameKey: a.Name("hello"),
	}).(a.Function))

	l := e.Scan(`(hello)`)
	tr := e.Read(b, l)
	c := a.ChildLocals(b)
	as.String("there", a.Eval(c, tr))
}

func testReaderError(t *testing.T, src string, err error) {
	as := assert.New(t)

	defer as.ExpectError(err)

	c := a.Variables{}
	l := e.Scan(s(src))
	tr := e.Read(c, l)
	a.Eval(c, tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", a.ErrStr(e.ListNotClosed))
	testReaderError(t, "[99 100 ", a.ErrStr(e.VectorNotClosed))
	testReaderError(t, "{:key 99", a.ErrStr(e.MapNotClosed))

	testReaderError(t, "99 100)", a.ErrStr(e.UnmatchedListEnd))
	testReaderError(t, "99 100]", a.ErrStr(e.UnmatchedVectorEnd))
	testReaderError(t, "99}", a.ErrStr(e.UnmatchedMapEnd))
	testReaderError(t, "{99}", a.ErrStr(e.MapNotPaired))

	testReaderError(t, "(", a.ErrStr(e.ListNotClosed))
	testReaderError(t, "'", a.ErrStr(e.PrefixedNotPaired, "sputter:quote"))
	testReaderError(t, "~@", a.ErrStr(e.PrefixedNotPaired, "sputter:unquote-splicing"))
	testReaderError(t, "~", a.ErrStr(e.PrefixedNotPaired, "sputter:unquote"))
}
