package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestLazySeq(t *testing.T) {
	var inc a.LazyResolver
	as := assert.New(t)

	i := 0
	inc = func() (bool, a.Value, a.Sequence) {
		if i >= 10 {
			return false, a.Nil, a.EmptyList
		}
		i++
		f := a.NewFloat(float64(i))
		return true, f, a.NewLazySequence(inc)
	}

	l := a.NewLazySequence(inc).Prepend(a.NewFloat(0))
	as.True(l.IsSequence())
	as.Number(0, l.First())
	as.Number(1, l.Rest().First())
	as.Number(2, l.Rest().Rest().First())
	as.Contains(":type lazy-sequence", l.Str())
}
