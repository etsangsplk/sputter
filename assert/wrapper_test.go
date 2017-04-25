package assert_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func expectError(as *assert.Wrapper, err string) {
	if r := recover(); r == err {
		return
	}
	as.Fail("should have failed with: " + err)
}

func TestTheStringTests(t *testing.T) {
	as := assert.New(t)

	as.String("hello", "hello")
	as.String("hello", a.Str("hello"))
	as.String(":hello", a.NewKeyword("hello"))

	defer expectError(as, a.Err(assert.InvalidTestExpression, 10))
	as.String("10", 10)
}

func TestTheFloatTests(t *testing.T) {
	as := assert.New(t)

	as.Number(10.5, a.NewFloat(10.5))
	as.Number(10, a.NewFloat(10.0))
	as.Number(10.5, 10.5)
	as.Number(10, 10)

	defer expectError(as, a.Err(assert.InvalidTestExpression, "10"))
	as.Number(10, "10")
}

func TestTheBoolTests(t *testing.T) {
	as := assert.New(t)

	as.True(a.True)
	as.True(true)
	as.False(a.False)
	as.False(false)
	as.Truthy(a.Str("hello"))
	as.Falsey(a.Nil)
}

func TestTheContainsTests(t *testing.T) {
	as := assert.New(t)

	as.Contains("there", a.Str("hello there!"))
	as.Contains("key", a.NewKeyword("iamkeyword"))
	as.NotContains("there", a.Str("hello nobody!"))
	as.NotContains("key", a.NewKeyword("iamnot"))
}

func TestTheIdenticalTests(t *testing.T) {
	as := assert.New(t)

	k1 := a.NewKeyword("yes")
	k2 := a.NewKeyword("no")

	as.Identical(k1, k1)
	as.NotIdentical(k1, k2)
}

func TestTheNilTests(t *testing.T) {
	as := assert.New(t)

	as.Nil(a.Nil)
	as.Nil(nil)
	as.NotNil("hello there")
}

func TestTheValueTests(t *testing.T) {
	as := assert.New(t)

	as.Equal(a.Str("hi"), "hi")
	as.Equal(a.NewKeyword("hello"), a.NewKeyword("hello"))
	as.Equal(a.NewFloat(10), a.NewFloat(10))
}
