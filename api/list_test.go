package api_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	a := assert.New(t)
	n := big.NewFloat(12)
	l := s.NewList(n)
	a.Equal(n, l.Value, "head is populated correctly")
	a.Equal(s.EmptyList, l.Rest, "list terminated properly")
}

func TestCons(t *testing.T) {
	a := assert.New(t)
	n1 := big.NewFloat(12)
	l1 := s.NewList(n1)
	a.Equal(n1, l1.Value, "1st head is populated correctly")
	a.Equal(s.EmptyList, l1.Rest, "list terminated properly")

	n2 := big.NewFloat(20.5)
	l2 := l1.Cons(n2)
	a.Equal(n2, l2.Value, "2nd head is populated correctly")
	a.Equal(l1, l2.Rest, "2nd tail is populated correctly")
}

func TestIterator(t *testing.T) {
	a := assert.New(t)
	n1 := big.NewFloat(12)
	l1 := s.NewList(n1)
	a.Equal(n1, l1.Value, "1st head is populated correctly")
	a.Equal(s.EmptyList, l1.Rest, "list terminated properly")

	n2 := big.NewFloat(20.5)
	l2 := l1.Cons(n2)
	a.Equal(n2, l2.Value, "2nd head is populated correctly")
	a.Equal(l1, l2.Rest, "2nd tail is populated correctly")

	sum := big.NewFloat(0.0)
	i := l2.Iterate()
	for {
		v, ok := i.Next()
		if !ok {
			break
		}
		fv := v.(*big.Float)
		sum.Add(sum, fv)
	}

	val, acc := sum.Float64()
	a.Equal(32.5, val, "values are summed correctly")
	a.EqualValues(0, acc, "should be no loss of accuracy")
}
