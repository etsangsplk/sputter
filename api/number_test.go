package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestParseNumber(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := f(12.8)
	n3 := f(12.8)

	as.Equal(n1, n2)
	as.Equal(n2, n3)

	defer as.ExpectError(a.ErrStr(a.ExpectedNumber, s(`'splosion!`)))
	a.ParseNumber("'splosion!")
}

func testExact(as *assert.Wrapper, n a.Number, expect float64) {
	val, exact := n.Float64()
	as.True(exact)
	as.Number(expect, val)
}

func TestConvertNumber(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := f(12.9)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := f(20)

	testExact(as, n1, 12.8)
	testExact(as, n2, 12.9)
	testExact(as, n3, 25)
	testExact(as, n4, 20)
	testExact(as, n5, 20)
}

func TestCompareNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := f(12.9)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := f(20)

	as.Compare(a.LessThan, n1, n2)
	as.Compare(a.LessThan, n1, n3)
	as.Compare(a.LessThan, n1, n4)
	as.Compare(a.GreaterThan, n3, n4)
	as.Compare(a.LessThan, n1, n5)
	as.Compare(a.EqualTo, n4, n5)
}

func TestStringifyNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := f(12.9)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := f(20)

	as.String("12.8", n1)
	as.String("12.9", n2)
	as.String("25/1", n3)
	as.String("20/1", n4)
	as.String("20", n5)
}

func testResult(as *assert.Wrapper, n a.Number, expect string) {
	expectNum := a.ParseNumber(s(expect))

	f, _ := expectNum.Float64()
	testExact(as, n, f)
	as.Equal(expectNum, n)
}

func TestNumberMath(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := f(13)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := f(20)
	n6 := a.ParseNumber("1/2")
	n7 := a.ParseNumber("12.9")

	testResult(as, n1.Add(n2), "25.8")
	testResult(as, n5.Sub(n2), "7")
	testResult(as, n4.Sub(n2), "7")
	testResult(as, n1.Mul(n3), "320")
	testResult(as, n3.Mul(n1), "320")
	testResult(as, n1.Add(n4), "32.8")
	testResult(as, n4.Add(n1), "32.8")
	testResult(as, n4.Div(n1), "1.5625")
	testResult(as, n5.Div(n1), "1.5625")
	testResult(as, n3.Sub(n4), "5")
	testResult(as, n3.Sub(n5), "5")
	testResult(as, n4.Add(n5), "40")
	testResult(as, n3.Add(n4), "45")
	testResult(as, n6.Mul(n4), "10")
	testResult(as, n4.Div(n6), "40")
	testResult(as, n1.Add(n7), "25.7")
}

func TestAssertNumber(t *testing.T) {
	as := assert.New(t)
	a.AssertNumber(f(99))

	defer as.ExpectError(a.ErrStr(a.ExpectedNumber, s("not a number")))
	a.AssertNumber(s("not a number"))
}

func TestAssertInteger(t *testing.T) {
	as := assert.New(t)
	a.AssertInteger(f(99))

	defer as.ExpectError(a.ErrStr(a.ExpectedInteger, f(99.5)))
	a.AssertInteger(f(99.5))
}
