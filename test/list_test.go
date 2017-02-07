package sputter_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	a := assert.New(t)
	num := big.NewFloat(12)
	list := s.NewList(num)
	a.Equal(num, list.Value, "head is populated correctly")
	a.Equal(s.EmptyList, list.Rest, "list terminated properly")
}

func TestCons(t *testing.T) {
	a := assert.New(t)
	num1 := big.NewFloat(12)
	list1 := s.NewList(num1)
	a.Equal(num1, list1.Value, "1st head is populated correctly")
	a.Equal(s.EmptyList, list1.Rest, "list terminated properly")

	num2 := big.NewFloat(20.5)
	list2 := list1.Cons(num2)
	a.Equal(num2, list2.Value, "2nd head is populated correctly")
	a.Equal(list1, list2.Rest, "2nd tail is populated correctly")
}

func TestIterator(t *testing.T) {
	a := assert.New(t)
	num1 := big.NewFloat(12)
	list1 := s.NewList(num1)
	a.Equal(num1, list1.Value, "1st head is populated correctly")
	a.Equal(s.EmptyList, list1.Rest, "list terminated properly")

	num2 := big.NewFloat(20.5)
	list2 := list1.Cons(num2)
	a.Equal(num2, list2.Value, "2nd head is populated correctly")
	a.Equal(list1, list2.Rest, "2nd tail is populated correctly")

	sum := big.NewFloat(0.0)
	iter := list2.Iterate()
	for {
		value, ok := iter.Next()
		if !ok {
			break
		}
		fv := value.(*big.Float)
		sum.Add(sum, fv)
	}

	val, accuracy := sum.Float64()
	a.Equal(32.5, val, "values are summed correctly")
	a.EqualValues(0, accuracy, "should be no loss of accuracy")
}
