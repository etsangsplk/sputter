package sputter

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	a := assert.New(t)
	num := big.NewFloat(12)
	list := NewList(num)
	a.Equal(num, list.value, "head is populated correctly")
	a.Equal(EmptyList, list.rest, "list terminated properly")
}

func TestCons(t *testing.T) {
	a := assert.New(t)
	num1 := big.NewFloat(12)
	list1 := NewList(num1)
	a.Equal(num1, list1.value, "1st head is populated correctly")
	a.Equal(EmptyList, list1.rest, "list terminated properly")

	num2 := big.NewFloat(20.5)
	list2 := list1.Cons(num2)
	a.Equal(num2, list2.value, "2nd head is populated correctly")
	a.Equal(list1, list2.rest, "2nd tail is populated correctly")
}

func TestIterator(t *testing.T) {
	a := assert.New(t)
	num1 := big.NewFloat(12)
	list1 := NewList(num1)
	a.Equal(num1, list1.value, "1st head is populated correctly")
	a.Equal(EmptyList, list1.rest, "list terminated properly")

	num2 := big.NewFloat(20.5)
	list2 := list1.Cons(num2)
	a.Equal(num2, list2.value, "2nd head is populated correctly")
	a.Equal(list1, list2.rest, "2nd tail is populated correctly")

	sum := big.NewFloat(0.0)
	next := list2.Iterate()
	for {
		value, found := next()
		if !found {
			break
		}
		fv := value.(*big.Float)
		sum.Add(sum, fv)
	}

	val, accuracy := sum.Float64()
	a.Equal(32.5, val, "values are summed correctly")
	a.EqualValues(0, accuracy, "should be no loss of accuracy")
}

func TestDuplicateEmpty(t *testing.T) {
	a := assert.New(t)

	l1 := EmptyList
	l2, _ := l1.duplicate()

	a.Exactly(l1, l2, "Empty Lists are same")
}
