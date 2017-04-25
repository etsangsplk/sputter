package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestLazyMapper(t *testing.T) {
	as := assert.New(t)

	c := 0
	l := a.NewList(s("last")).Prepend(s("middle")).Prepend(s("first"))
	w := a.Map(l, func(v a.Value) a.Value {
		c++
		return s("this is the " + string(v.(a.Str)))
	})

	as.Number(0, c)

	v1 := w.First()
	as.Number(1, c)
	as.String("this is the first", v1)

	v2 := w.Rest().First()
	as.Number(2, c)
	as.String("this is the middle", v2)

	v3 := w.Rest().Rest().First()
	as.Number(3, c)
	as.String("this is the last", v3)

	r1 := w.Rest().Rest().Rest()
	as.False(r1.IsSequence())

	p1 := w.Prepend(s("not mapped"))
	p2 := p1.Prepend(s("also not mapped"))
	v4 := p1.First()
	r2 := p1.Rest()
	as.Number(3, c)
	as.String("not mapped", v4)
	as.Equal(w, r2)
	as.String("also not mapped", p2.First())
}

func TestLazyFilter(t *testing.T) {
	as := assert.New(t)

	c := 0
	l := a.NewList(s("last")).Prepend(s("filtered out")).Prepend(s("first"))
	w := a.Filter(l, func(v a.Value) bool {
		c++
		return string(v.(a.Str)) != "filtered out"
	})

	as.Number(0, c)

	v1 := w.First()
	as.Number(1, c)
	as.String("first", v1)

	v2 := w.Rest().First()
	as.Number(3, c)
	as.String("last", v2)

	r1 := w.Rest().Rest()
	as.False(r1.IsSequence())

	p := w.Prepend(s("filtered out"))
	v4 := p.First()
	r2 := p.Rest()
	as.Number(3, c)
	as.String("filtered out", v4)
	as.Equal(w, r2)
}

func TestLazyFilteredAndMapped(t *testing.T) {
	as := assert.New(t)

	c1 := 0
	c2 := 0
	l := a.NewList(s("last")).Prepend(s("middle")).Prepend(s("first"))
	w1 := a.Filter(l, func(v a.Value) bool {
		c1++
		return string(v.(a.Str)) != "middle"
	})
	w2 := a.Map(w1, func(v a.Value) a.Value {
		c2++
		return s("this is the " + string(v.(a.Str)))
	})

	as.Number(0, c1)
	as.Number(0, c2)

	v1 := w2.First()
	as.Number(1, c1)
	as.Number(1, c2)
	as.String("this is the first", v1)

	v2 := w2.Rest().First()
	as.Number(3, c1)
	as.Number(2, c2)
	as.String("this is the last", v2)

	r1 := w2.Rest().Rest()
	as.False(r1.IsSequence())
}

func testNext(as *assert.Wrapper, i *a.Iterator, expected a.Value) {
	v, ok := i.Next()
	as.True(ok)
	as.Equal(expected, v)
}

func TestLazyConcat(t *testing.T) {
	as := assert.New(t)

	l1 := a.NewList(s("last")).Prepend(s("middle")).Prepend(s("first"))
	l2 := a.EmptyList
	l3 := a.Vector{f(1), f(2), f(3)}
	l4 := a.NewList(s("blah3")).Prepend(s("blah2")).Prepend(s("blah1"))
	l5 := a.EmptyList

	w1 := a.Concat(a.Vector{l1, l2, l3, l4, l5})
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

	s := `("first" "middle" "last" 1 2 3 "blah1" "blah2" "blah3")`
	as.String(s, w1)
}

func TestReduce(t *testing.T) {
	as := assert.New(t)

	add := a.NewFunction(
		func(_ a.Context, args a.Sequence) a.Value {
			v := args.(a.Vector)
			return v[0].(*a.Number).Add(v[1].(*a.Number))
		},
	)

	as.Number(30, a.Reduce(nil, a.Vector{f(10), f(20)}, add))
	as.Number(60, a.Reduce(nil, a.Vector{f(10), f(20), f(30)}, add))
	as.Number(100, a.Reduce(nil, a.Vector{f(10), f(20), f(30), f(40)}, add))

	err := a.Err(a.BadMinimumArity, 2, 1)
	defer expectError(as, err)
	a.Reduce(nil, a.Vector{f(10)}, add)
}

func TestTakeDrop(t *testing.T) {
	as := assert.New(t)

	s1 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	v1 := make(a.Vector, len(s1))
	for i, e := range s1 {
		v1[i] = s(e)
	}

	t1 := a.Take(v1, 4)
	t2 := t1.Prepend(s("0"))
	d1 := a.Drop(v1, 4)
	d2 := d1.Prepend(s("4"))
	t3 := a.Take(d1, 6)
	d3 := a.Drop(t3, 6)
	d4 := a.Drop(t3, 8)

	as.String(`("1" "2" "3" "4")`, t1)
	as.String(`("0" "1" "2" "3" "4")`, t2)
	as.String(`("5" "6" "7" "8" "9" "10")`, d1)
	as.String(`("4" "5" "6" "7" "8" "9" "10")`, d2)
	as.String(`("5" "6" "7" "8" "9" "10")`, t3)
	as.String(`()`, d3)
	as.String(`()`, d4)
}
