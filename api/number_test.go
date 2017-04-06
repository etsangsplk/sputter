package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestParseNumber(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := a.NewFloat(12.8)
	n3 := a.NewFloat(12.8)

	as.Equal(a.EqualTo, n1.Cmp(n2))
	as.Equal(a.EqualTo, n2.Cmp(n3))

	defer expectError(as, a.Err(a.ExpectedNumber, `"'splosion!"`))
	a.ParseNumber("'splosion!")
}

func testExact(as *assert.Assertions, n *a.Number, expect float64) {
	val, exact := n.Float64()
	as.True(exact, "non-exact precision returned")
	as.Equal(expect, val)
}

func TestConvertNumber(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := a.NewFloat(12.9)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := a.NewFloat(20)

	testExact(as, n1, 12.8)
	testExact(as, n2, 12.9)
	testExact(as, n3, 25)
	testExact(as, n4, 20)
	testExact(as, n5, 20)
}

func TestCompareNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := a.NewFloat(12.9)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := a.NewFloat(20)

	as.Equal(a.LessThan, n1.Cmp(n2))
	as.Equal(a.LessThan, n1.Cmp(n3))
	as.Equal(a.LessThan, n1.Cmp(n4))
	as.Equal(a.GreaterThan, n3.Cmp(n4))
	as.Equal(a.LessThan, n1.Cmp(n5))
	as.Equal(a.EqualTo, n4.Cmp(n5))
}

func TestStringifyNumbers(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := a.NewFloat(12.9)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := a.NewFloat(20)

	as.Equal("12.8", a.String(n1))
	as.Equal("12.9", a.String(n2))
	as.Equal("25/1", a.String(n3))
	as.Equal("20/1", a.String(n4))
	as.Equal("20", a.String(n5))
}

func testResult(as *assert.Assertions, n *a.Number, expect string) {
	expectNum := a.ParseNumber(expect)

	f, _ := expectNum.Float64()
	testExact(as, n, f)
	as.Equal(a.EqualTo, expectNum.Cmp(n))
}

func TestNumberMath(t *testing.T) {
	as := assert.New(t)
	n1 := a.ParseNumber("12.8")
	n2 := a.NewFloat(13)
	n3 := a.ParseNumber("50/2")
	n4 := a.NewRatio(40, 2)
	n5 := a.NewFloat(20)
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
	a.AssertNumber(a.NewFloat(99))

	defer expectError(as, a.Err(a.ExpectedNumber, `"not a number"`))
	a.AssertNumber("not a number")
}

func TestAssertInteger(t *testing.T) {
	as := assert.New(t)
	a.AssertInteger(a.NewFloat(99))

	defer expectError(as, a.Err(a.ExpectedInteger, `99.5`))
	a.AssertInteger(a.NewFloat(99.5))
}
