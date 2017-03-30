package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

var helloThere = a.NewFunction(func(_ a.Context, _ a.Sequence) a.Value {
	return "there"
}).WithMetadata(a.Metadata{
	a.MetaName: a.Name("hello"),
}).(a.Function)

func TestSimpleList(t *testing.T) {
	as := assert.New(t)
	n := a.NewFloat(12)
	l := a.NewList(n)
	as.Equal(n, l.First(), "head is populated correctly")
	as.Equal(a.EmptyList, l.Rest(), "list terminated properly")
}

func TestList(t *testing.T) {
	as := assert.New(t)
	n1 := a.NewFloat(12)
	l1 := a.NewList(n1)
	as.Equal(n1, l1.First(), "1st head is populated correctly")
	as.Equal(a.EmptyList, l1.Rest(), "list terminated properly")
	as.False(l1.Rest().IsSequence(), "list terminated properly")

	n2 := a.NewFloat(20.5)
	l2 := l1.Prepend(n2).(*a.List)

	as.Equal("()", a.EmptyList.String())
	as.Equal("(20.5 12)", l2.String())
	as.Equal(n2, l2.First(), "2nd head is populated correctly")
	as.Equal(l1, l2.Rest(), "2nd tail is populated correctly")
	as.Equal(2, l2.Count(), "2nd list count is correct")
	as.Equal(a.NewFloat(12), l2.Get(1), "get(int) works")
	as.Equal(2, a.Count(l2), "2nd list general count is correct")

	as.Equal(a.Nil, a.EmptyList.Get(1), "get from empty list")
}

func TestIterator(t *testing.T) {
	as := assert.New(t)
	n1 := a.NewFloat(12)
	l1 := a.NewList(n1)
	as.Equal(n1, l1.First(), "1st head is populated correctly")
	as.Equal(a.EmptyList, l1.Rest(), "list terminated properly")
	as.False(l1.Rest().IsSequence(), "list terminated properly")

	n2 := a.NewFloat(20.5)
	l2 := l1.Prepend(n2)
	as.Equal(n2, l2.First(), "2nd head is populated correctly")
	as.Equal(l1, l2.Rest(), "2nd tail is populated correctly")

	sum := a.NewFloat(0.0)
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
	as.Equal(32.5, val, "values are summed correctly")
	as.EqualValues(true, exact, "should be no loss of accuracy")
}

func TestListEval(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put(helloThere.Name(), helloThere)

	fl := a.NewList(helloThere)
	as.Equal("there", a.Eval(c, fl), "function-based list eval")

	sl := a.NewList(a.NewLocalSymbol("hello"))
	as.Equal("there", a.Eval(c, sl), "symbol-based list eval")

	as.Equal(a.EmptyList, a.Eval(c, a.EmptyList), "empty list eval")
}

func testBrokenEval(t *testing.T, seq a.Sequence, err string) {
	as := assert.New(t)

	defer expectError(as, err)
	c := a.NewContext()
	a.Eval(c, seq)
}

func TestNonFunction(t *testing.T) {
	seq := a.NewList("foo").Prepend(a.NewLocalSymbol("unknown"))
	testBrokenEval(t, seq, a.ExpectedApplicable)
}
