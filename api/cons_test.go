package api_test

import (
	"fmt"
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

var helloThere = &s.Function{
	Name: "hello",
	Exec: func(c s.Context, args s.Sequence) s.Value {
		return "there"
	},
}

func TestList(t *testing.T) {
	a := assert.New(t)
	n := big.NewFloat(12)
	l := s.NewList(n)
	a.Equal(n, l.Car, "head is populated correctly")
	a.Equal(s.Nil, l.Cdr, "list terminated properly")
}

func TestCons(t *testing.T) {
	a := assert.New(t)
	n1 := big.NewFloat(12)
	l1 := s.NewList(n1)
	a.Equal(n1, l1.Car, "1st head is populated correctly")
	a.Equal(s.Nil, l1.Cdr, "list terminated properly")

	n2 := big.NewFloat(20.5)
	l2 := &s.Cons{Car: n2, Cdr: l1}

	a.Equal("(20.5 12)", l2.String())
	a.Equal(n2, l2.Car, "2nd head is populated correctly")
	a.Equal(l1, l2.Cdr, "2nd tail is populated correctly")
	a.Equal(2, l2.Count(), "2nd list count is correct")
	a.Equal(2, s.Count(l2), "2nd list general count is correct")

	a.Equal(s.Nil, s.Nil.Get(1), "get from empty list")
}

func TestIterator(t *testing.T) {
	a := assert.New(t)
	n1 := big.NewFloat(12)
	l1 := s.NewList(n1)
	a.Equal(n1, l1.Car, "1st head is populated correctly")
	a.Equal(s.Nil, l1.Cdr, "list terminated properly")

	n2 := big.NewFloat(20.5)
	l2 := &s.Cons{Car: n2, Cdr: l1}
	a.Equal(n2, l2.Car, "2nd head is populated correctly")
	a.Equal(l1, l2.Cdr, "2nd tail is populated correctly")

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

func TestDotCons(t *testing.T) {
	a := assert.New(t)

	c1 := &s.Cons{Car: "are", Cdr: "You?"}
	c2 := &s.Cons{Car: "how", Cdr: c1}
	c3 := &s.Cons{Car: "hello", Cdr: c2}

	a.Equal("are", c3.Get(2), "list ending in Cons gets Car, not Cons")
	a.Equal(3, c3.Count(), "list ending in Cons has right count")
	a.Equal(3, s.Count(c3), "list ending in Cons has right general count")

	i := c3.Iterate()
	v1, _ := i.Next()
	s1 := i.Rest()
	v2, _ := i.Next()
	s2 := i.Rest()
	v3, _ := i.Next()
	v4, ok := i.Next()

	a.Equal("(are . You?)", c1.String())
	a.Equal("(hello how are . You?)", c3.String())

	a.Equal(2, s1.(s.Finite).Count(), "slicing correctly")
	a.Equal(1, s2.(s.Finite).Count(), "slicing correctly")

	a.Equal("hello", v1, "first iteration correct")
	a.Equal("how", v2, "second iteration correct")
	a.Equal(c1, v3, "third iteration correct")
	a.Equal(s.Nil, v4, "fourth iteration correct")
	a.False(ok, "fourth iteration failed properly")

	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Sprintf(s.IndexNotCons, 3)
			a.Equal(err, rec, "bad Get index")
			return
		}
		a.Fail("bad Get index didn't panic")
	}()

	c3.Get(3) // will explode
}

func TestConsEval(t *testing.T) {
	a := assert.New(t)

	c := s.NewContext()
	c.Put(helloThere.Name, helloThere)

	fl := s.NewList(helloThere)
	a.Equal("there", fl.Eval(c), "function-based list eval")

	sl := s.NewList(&s.Symbol{Name: "hello"})
	a.Equal("there", sl.Eval(c), "symbol-based list eval")

	a.Equal(s.Nil, s.Nil.Eval(c), "empty list eval")
}

func testBrokenEval(t *testing.T, cons *s.Cons, err string) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(err, rec, "eval panics properly")
			return
		}
		a.Fail("eval should panic")
	}()

	c := s.NewContext()
	cons.Eval(c)
}

func TestNonList(t *testing.T) {
	cons := &s.Cons{Car: helloThere, Cdr: "uh-oh"}
	testBrokenEval(t, cons, s.NonList)
}

func TestNonFunction(t *testing.T) {
	cons := &s.Cons{
		Car: &s.Symbol{Name: "unknown"},
		Cdr: s.NewList("foo"),
	}
	testBrokenEval(t, cons, s.NonFunction)
}
