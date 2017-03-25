package api_test

import (
	"fmt"
	"strings"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestLazyMapper(t *testing.T) {
	as := assert.New(t)

	c := 0
	l := a.NewList("last").Prepend("middle").Prepend("first")
	w := a.NewMapper(l, func(v a.Value) a.Value {
		c++
		return "this is the " + v.(string)
	})

	as.Equal(0, c, "nothing has been processed")

	v1 := w.First()
	as.Equal(1, c, "first element has been processed")
	as.Equal("this is the first", v1, "first element mapped")

	v2 := w.Rest().First()
	as.Equal(2, c, "second element has been processed")
	as.Equal("this is the middle", v2, "second element mapped")

	v3 := w.Rest().Rest().First()
	as.Equal(3, c, "third element has been processed")
	as.Equal("this is the last", v3, "third element mapped")

	r1 := w.Rest().Rest().Rest()
	as.False(r1.IsSequence(), "finished")

	p := w.Prepend("not mapped")
	v4 := p.First()
	r2 := p.Rest()
	as.Equal(3, c, "prepend doesn't trigger mapper")
	as.Equal("not mapped", v4, "prepended element retrieved")
	as.Equal(w, r2, "prepended rest is the original sequence")
}

func TestLazyFilter(t *testing.T) {
	as := assert.New(t)

	c := 0
	l := a.NewList("last").Prepend("filtered out").Prepend("first")
	w := a.NewFilter(l, func(v a.Value) bool {
		c++
		return v.(string) != "filtered out"
	})

	as.Equal(0, c, "nothing has been processed")

	v1 := w.First()
	as.Equal(1, c, "first element has been processed")
	as.Equal("first", v1, "first element passed")

	v2 := w.Rest().First()
	as.Equal(3, c, "second element has been skipped")
	as.Equal("last", v2, "third element passed")

	r1 := w.Rest().Rest()
	as.False(r1.IsSequence(), "finished")

	p := w.Prepend("filtered out")
	v4 := p.First()
	r2 := p.Rest()
	as.Equal(3, c, "prepend doesn't trigger filter")
	as.Equal("filtered out", v4, "prepended element retrieved")
	as.Equal(w, r2, "prepended rest is the original sequence")
}

func TestLazyFilteredAndMapped(t *testing.T) {
	as := assert.New(t)

	c1 := 0
	c2 := 0
	l := a.NewList("last").Prepend("middle").Prepend("first")
	w1 := a.NewFilter(l, func(v a.Value) bool {
		c1++
		return v.(string) != "middle"
	})
	w2 := a.NewMapper(w1, func(v a.Value) a.Value {
		c2++
		return "this is the " + v.(string)
	})

	as.Equal(0, c1, "nothing has been processed")
	as.Equal(0, c2, "nothing has been processed")

	v1 := w2.First()
	as.Equal(1, c1, "first element has been processed")
	as.Equal(1, c2, "first element has been processed")
	as.Equal("this is the first", v1, "first element mapped")

	v2 := w2.Rest().First()
	as.Equal(3, c1, "second element has been skipped")
	as.Equal(2, c2, "second element is first mapper sees")
	as.Equal("this is the last", v2, "last element mapped")

	r1 := w2.Rest().Rest()
	as.False(r1.IsSequence(), "finished")
}

func testNext(as *assert.Assertions, i *a.Iterator, expected a.Value) {
	v, ok := i.Next()
	as.True(ok, "iterator has more values")
	as.Equal(expected, v, fmt.Sprintf("%s expected", expected))
}

func TestLazyConcat(t *testing.T) {
	as := assert.New(t)

	l1 := a.NewList("last").Prepend("middle").Prepend("first")
	l2 := a.EmptyList
	l3 := a.Vector{1, 2, 3}
	l4 := a.NewList("blah3").Prepend("blah2").Prepend("blah1")
	l5 := a.EmptyList

	w1 := a.NewConcat(a.Vector{l1, l2, l3, l4, l5})
	w2 := w1.Prepend("I was prepended")

	i := a.Iterate(w2)

	testNext(as, i, "I was prepended")
	testNext(as, i, "first")
	testNext(as, i, "middle")
	testNext(as, i, "last")
	testNext(as, i, 1)
	testNext(as, i, 2)
	testNext(as, i, 3)
	testNext(as, i, "blah1")
	testNext(as, i, "blah2")
	testNext(as, i, "blah3")

	_, ok := i.Next()
	as.False(ok, "end of sequence")
}

func TestLazyToString(t *testing.T) {
	as := assert.New(t)

	id := func(v a.Value) a.Value { return v }
	all := func(v a.Value) bool { return true }

	v := a.Vector{}
	as.True(strings.HasPrefix(a.String(a.NewConcat(v)), "(concat :instance"))
	as.True(strings.HasPrefix(a.String(a.NewFilter(v, all)), "(filter :instance"))
	as.True(strings.HasPrefix(a.String(a.NewMapper(v, id)), "(map :instance"))
}
