package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestLazyList(t *testing.T) {
	as := assert.New(t)

	c := 0
	l := a.NewList("last").Prepend("middle").Prepend("first")
	w := a.NewLazySequence(l, func(v a.Value) a.Value {
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
