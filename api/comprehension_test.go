package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestMap(t *testing.T) {
	as := assert.New(t)

	l := a.NewList(s("first"), s("middle"), s("last"))
	fn := a.NewExecFunction(func(_ a.Context, args a.Vector) a.Value {
		return s("this is the " + string(args[0].(a.Str)))
	})
	w := a.Map(nil, l, fn)

	v1 := w.First()
	as.String("this is the first", v1)

	v2 := w.Rest().First()
	as.String("this is the middle", v2)

	v3 := w.Rest().Rest().First()
	as.String("this is the last", v3)

	r1 := w.Rest().Rest().Rest()
	as.False(r1.IsSequence())

	p1 := w.Prepend(s("not mapped"))
	p2 := p1.Prepend(s("also not mapped"))
	v4 := p1.First()
	r2 := p1.Rest()

	as.String("not mapped", v4)
	as.Equal(w.First(), r2.First())
	as.String("also not mapped", p2.First())
}

func TestMapParallel(t *testing.T) {
	as := assert.New(t)

	add := a.NewExecFunction(func(_ a.Context, args a.Vector) a.Value {
		return args[0].(a.Number).Add(args[1].(a.Number))
	})

	s1 := a.NewList(f(1), f(2), f(3), f(4))
	s2 := a.NewVector(f(5), f(10), f(15), f(20), f(30))

	w := a.MapParallel(nil, a.Vector{s1, s2}, add)

	as.Number(6, w.First())
	as.Number(12, w.Rest().First())
	as.Number(18, w.Rest().Rest().First())
	as.Number(24, w.Rest().Rest().Rest().First())

	s3 := w.Rest().Rest().Rest().Rest()
	as.False(s3.IsSequence())
}

func TestFilter(t *testing.T) {
	as := assert.New(t)

	l := a.NewList(s("first"), s("filtered out"), s("last"))
	fn := a.NewExecFunction(func(_ a.Context, args a.Vector) a.Value {
		return a.Bool(string(args[0].(a.Str)) != "filtered out")
	})
	w := a.Filter(nil, l, fn)

	v1 := w.First()
	as.String("first", v1)

	v2 := w.Rest().First()
	as.String("last", v2)

	r1 := w.Rest().Rest()
	as.False(r1.IsSequence())

	p := w.Prepend(s("filtered out"))
	v4 := p.First()
	r2 := p.Rest()
	as.String("filtered out", v4)
	as.Equal(w.First(), r2.First())
}

func TestFilteredAndMapped(t *testing.T) {
	as := assert.New(t)

	l := a.NewList(s("first"), s("middle"), s("last"))
	fn1 := a.NewExecFunction(func(_ a.Context, args a.Vector) a.Value {
		return a.Bool(string(args[0].(a.Str)) != "middle")
	})
	w1 := a.Filter(nil, l, fn1)

	fn2 := a.NewExecFunction(func(_ a.Context, args a.Vector) a.Value {
		return s("this is the " + string(args[0].(a.Str)))
	})
	w2 := a.Map(nil, w1, fn2)

	v1 := w2.First()
	as.String("this is the first", v1)

	v2 := w2.Rest().First()
	as.String("this is the last", v2)

	r1 := w2.Rest().Rest()
	as.False(r1.IsSequence())
}

func testNext(as *assert.Wrapper, i *a.Iterator, expected a.Value) {
	v, ok := i.Next()
	as.True(ok)
	as.Equal(expected, v)
}

func TestConcat(t *testing.T) {
	as := assert.New(t)

	l1 := a.NewList(s("first"), s("middle"), s("last"))
	l2 := a.EmptyList
	l3 := a.NewVector(f(1), f(2), f(3))
	l4 := a.NewList(s("blah1"), s("blah2"), s("blah3"))
	l5 := a.EmptyList

	w1 := a.Concat(a.NewVector(l1, l2, l3, l4, l5))
	w2 := w1.Prepend(s("I was prepended"))

	i := a.Iterate(w2)

	testNext(as, i, s("I was prepended"))
	testNext(as, i, s("first"))
	testNext(as, i, s("middle"))
	testNext(as, i, s("last"))
	testNext(as, i, f(1))
	testNext(as, i, f(2))
	testNext(as, i, f(3))
	testNext(as, i, s("blah1"))
	testNext(as, i, s("blah2"))
	testNext(as, i, s("blah3"))

	_, ok := i.Next()
	as.False(ok)

	expect := `("first" "middle" "last" 1 2 3 "blah1" "blah2" "blah3")`
	as.String(expect, a.MakeSequenceStr(w1))
}

func TestReduce(t *testing.T) {
	as := assert.New(t)

	add := a.NewExecFunction(func(_ a.Context, args a.Vector) a.Value {
		return args[0].(a.Number).Add(args[1].(a.Number))
	})

	as.Number(30, a.Reduce(nil, a.NewVector(f(10), f(20)), add))
	as.Number(60, a.Reduce(nil, a.NewVector(f(10), f(20), f(30)), add))
	as.Number(100, a.Reduce(nil, a.NewVector(f(10), f(20), f(30), f(40)), add))
}

func TestTakeDrop(t *testing.T) {
	as := assert.New(t)

	s1 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	v1 := a.NewVector()
	for _, e := range s1 {
		v1 = v1.Conjoin(s(e)).(a.Vector)
	}

	t1 := a.Take(v1, 4)
	t2 := t1.Prepend(s("0"))
	d1 := a.Drop(v1, 4)
	d2 := d1.Prepend(s("4"))
	t3 := a.Take(d1, 6)
	d3 := a.Drop(t3, 6)
	d4 := a.Drop(t3, 8)

	as.String(`("1" "2" "3" "4")`, a.MakeSequenceStr(t1))
	as.String(`("0" "1" "2" "3" "4")`, a.MakeSequenceStr(t2))
	as.String(`("5" "6" "7" "8" "9" "10")`, a.MakeSequenceStr(d1))
	as.String(`("4" "5" "6" "7" "8" "9" "10")`, a.MakeSequenceStr(d2))
	as.String(`("5" "6" "7" "8" "9" "10")`, a.MakeSequenceStr(t3))
	as.False(d3.IsSequence())
	as.False(d4.IsSequence())
}
