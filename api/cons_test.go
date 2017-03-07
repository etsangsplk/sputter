package api_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

var helloThere = &a.Function{
	Name: "hello",
	Apply: func(c a.Context, args a.Sequence) a.Value {
		return "there"
	},
}

func TestSimpleList(t *testing.T) {
	as := assert.New(t)
	n := big.NewFloat(12)
	l := a.NewList(n)
	as.Equal(n, l.Car, "head is populated correctly")
	as.Equal(a.Nil, l.Cdr, "list terminated properly")
}

func TestCons(t *testing.T) {
	as := assert.New(t)
	n1 := big.NewFloat(12)
	l1 := a.NewList(n1)
	as.Equal(n1, l1.Car, "1st head is populated correctly")
	as.Equal(a.Nil, l1.Cdr, "list terminated properly")

	n2 := big.NewFloat(20.5)
	l2 := &a.Cons{Car: n2, Cdr: l1}

	as.Equal("()", a.Nil.String())
	as.Equal("(20.5 12)", l2.String())
	as.Equal(n2, l2.Car, "2nd head is populated correctly")
	as.Equal(l1, l2.Cdr, "2nd tail is populated correctly")
	as.Equal(2, l2.Count(), "2nd list count is correct")
	as.Equal(2, a.Count(l2), "2nd list general count is correct")

	as.Equal(a.Nil, a.Nil.Get(1), "get from empty list")
}

func TestIterator(t *testing.T) {
	as := assert.New(t)
	n1 := big.NewFloat(12)
	l1 := a.NewList(n1)
	as.Equal(n1, l1.Car, "1st head is populated correctly")
	as.Equal(a.Nil, l1.Cdr, "list terminated properly")

	n2 := big.NewFloat(20.5)
	l2 := &a.Cons{Car: n2, Cdr: l1}
	as.Equal(n2, l2.Car, "2nd head is populated correctly")
	as.Equal(l1, l2.Cdr, "2nd tail is populated correctly")

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
	as.Equal(32.5, val, "values are summed correctly")
	as.EqualValues(0, acc, "should be no loss of accuracy")
}

func TestDotCons(t *testing.T) {
	as := assert.New(t)

	c1 := &a.Cons{Car: "are", Cdr: "You?"}
	c2 := &a.Cons{Car: "how", Cdr: c1}
	c3 := &a.Cons{Car: "hello", Cdr: c2}

	as.Equal("are", c3.Get(2), "list ending in Cons gets Car, not Cons")
	as.Equal(3, c3.Count(), "list ending in Cons has right count")
	as.Equal(3, a.Count(c3), "list ending in Cons has right general count")

	i := c3.Iterate()
	v1, _ := i.Next()
	s1 := i.Rest()
	v2, _ := i.Next()
	s2 := i.Rest()
	v3, _ := i.Next()
	v4, ok := i.Next()

	as.Equal("(are . You?)", c1.String())
	as.Equal("(hello how are . You?)", c3.String())

	as.Equal(2, s1.(a.Finite).Count(), "slicing correctly")
	as.Equal(1, s2.(a.Finite).Count(), "slicing correctly")

	as.Equal("hello", v1, "first iteration correct")
	as.Equal("how", v2, "second iteration correct")
	as.Equal(c1, v3, "third iteration correct")
	as.Equal(a.Nil, v4, "fourth iteration correct")
	as.False(ok, "fourth iteration failed properly")

	defer func() {
		if rec := recover(); rec != nil {
			as.Equal(a.ExpectedCons, rec, "bad Get index")
			return
		}
		as.Fail("bad Get index didn't panic")
	}()

	c3.Get(3)
}

func TestConsEval(t *testing.T) {
	as := assert.New(t)

	c := a.NewContext()
	c.Put(helloThere.Name, helloThere)

	fl := a.NewList(helloThere)
	as.Equal("there", fl.Eval(c), "function-based list eval")

	sl := a.NewList(&a.Symbol{Name: "hello"})
	as.Equal("there", sl.Eval(c), "symbol-based list eval")

	as.Equal(a.Nil, a.Nil.Eval(c), "empty list eval")
}

func testBrokenEval(t *testing.T, cons *a.Cons, err string) {
	as := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			as.Equal(err, rec, "eval panics properly")
			return
		}
		as.Fail("eval should panic")
	}()

	c := a.NewContext()
	cons.Eval(c)
}

func TestNonList(t *testing.T) {
	cons := &a.Cons{Car: helloThere, Cdr: "uh-oh"}
	testBrokenEval(t, cons, a.ExpectedList)
}

func TestNonFunction(t *testing.T) {
	cons := &a.Cons{
		Car: &a.Symbol{Name: "unknown"},
		Cdr: a.NewList("foo"),
	}
	testBrokenEval(t, cons, a.ExpectedFunction)
}

func TestAssertCons(t *testing.T) {
	as := assert.New(t)
	a.AssertCons(&a.Cons{Car: "hello", Cdr: "there"})

	defer expectError(as, a.ExpectedCons)
	a.AssertCons(big.NewFloat(99))
}

func TestAssertSequence(t *testing.T) {
	as := assert.New(t)
	a.AssertSequence(a.NewList("hello"))

	defer expectError(as, a.ExpectedSequence)
	a.AssertSequence(big.NewFloat(99))
}
