package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

var helloThere = a.NewFunction(func(_ a.Context, _ a.Sequence) a.Value {
	return s("there")
}).WithMetadata(a.Metadata{
	a.MetaName: a.Name("hello"),
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

	c := a.NewContext()
	as.Equal(f(12), l2.Apply(c, a.NewList(f(1))))
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
		fv := v.(*a.Number)
		sum = sum.Add(fv)
	}

	val, exact := sum.Float64()
	as.Number(32.5, val)
	as.True(exact)
}

func TestListEval(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put(helloThere.Name(), helloThere)

	fl := a.NewList(helloThere)
	as.String("there", a.Eval(c, fl))

	sl := a.NewList(a.NewLocalSymbol("hello"))
	as.String("there", a.Eval(c, sl))

	as.Equal(a.EmptyList, a.Eval(c, a.EmptyList))
}

func testBrokenEval(t *testing.T, seq a.Sequence, err string) {
	as := assert.New(t)

	defer expectError(as, err)
	c := a.NewContext()
	a.Eval(c, seq)
}

func TestNonFunction(t *testing.T) {
	err := a.Err(a.UnknownSymbol, "unknown")
	seq := a.NewList(s("foo")).Prepend(a.NewLocalSymbol("unknown"))
	testBrokenEval(t, seq, err)
}

func TestListExplosion(t *testing.T) {
	as := assert.New(t)

	seq := a.NewList(s("foo"))
	idx := f(3)
	err := a.Err(a.IndexNotFound, idx)

	v := seq.Apply(a.NewContext(), a.Vector{idx, s("default")})
	as.String("default", v)

	defer expectError(as, err)
	seq.Apply(a.NewContext(), a.Vector{idx})
}
