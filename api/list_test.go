package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

var helloThere = a.NewExecFunction(func(_ a.Context, _ a.Vector) a.Value {
	return s("there")
}).WithMetadata(a.Properties{
	a.NameKey: a.Name("hello"),
}).(a.Function)

func TestSimpleList(t *testing.T) {
	as := assert.New(t)
	n := f(12)
	l := a.NewList(n)
	as.Equal(n, l.First())
	as.Equal(a.EmptyList, l.Rest())
}

func TestList(t *testing.T) {
	as := assert.New(t)
	n1 := f(12)
	l1 := a.NewList(n1)

	as.Equal(n1, l1.First())
	as.Equal(a.EmptyList, l1.Rest())
	as.False(l1.Rest().IsSequence())

	n2 := f(20.5)
	l2 := l1.Prepend(n2).(*a.List)

	as.String("()", a.EmptyList)
	as.String("(20.5 12)", l2)
	as.Equal(n2, l2.First())
	as.Identical(l1, l2.Rest())
	as.Number(2, l2.Count())

	r, ok := l2.ElementAt(1)
	as.True(ok)
	as.Equal(f(12), r)
	as.Number(2, a.Count(l2))

	r, ok = a.EmptyList.ElementAt(1)
	as.False(ok)
	as.Equal(a.Nil, r)

	c := a.Variables{}
	as.Equal(f(12), a.Apply(c, l2, a.NewList(f(1))))
}

func TestIterator(t *testing.T) {
	as := assert.New(t)
	n1 := f(12)
	l1 := a.NewList(n1)
	as.Equal(n1, l1.First())
	as.Identical(a.EmptyList, l1.Rest())
	as.False(l1.Rest().IsSequence())

	n2 := f(20.5)
	l2 := l1.Conjoin(n2)
	as.Equal(n2, l2.First())
	as.Identical(l1, l2.Rest())

	sum := f(0.0)
	i := a.Iterate(l2)
	for {
		v, ok := i.Next()
		if !ok {
			break
		}
		fv := v.(a.Number)
		sum = sum.Add(fv)
	}

	val, exact := sum.Float64()
	as.Number(32.5, val)
	as.True(exact)
}

func TestListEval(t *testing.T) {
	as := assert.New(t)

	c := a.Variables{}
	n := helloThere.Metadata().MustGet(a.NameKey).(a.Name)
	c.Put(n, helloThere)

	fl := a.NewList(helloThere)
	as.String("there", a.Eval(c, fl))

	sym := a.NewLocalSymbol("hello")
	sl := a.NewList(sym)
	as.String("there", a.Eval(c, sl))

	as.Equal(a.EmptyList, a.Eval(c, a.EmptyList))
}

func testBrokenEval(t *testing.T, val a.Value, err error) {
	as := assert.New(t)

	defer as.ExpectError(err)
	c := a.Variables{}
	a.Eval(c, val)
}

func TestNonFunction(t *testing.T) {
	e := a.ErrStr(a.UnknownSymbol, "unknown")
	sym := a.NewLocalSymbol("unknown")
	seq := a.NewList(sym, s("foo"))
	testBrokenEval(t, seq, e)

	e = cvtErr("*api.dec", "api.Applicable", "Apply")
	testBrokenEval(t, a.NewList(f(99)), e)
}

func TestListExplosion(t *testing.T) {
	as := assert.New(t)

	seq := a.NewList(s("foo"))
	idx := f(3)
	err := a.ErrStr(a.IndexNotFound, idx)

	v := seq.Apply(a.Variables{}, a.Vector{idx, s("default")})
	as.String("default", v)

	defer as.ExpectError(err)
	seq.Apply(a.Variables{}, a.Vector{idx})
}
