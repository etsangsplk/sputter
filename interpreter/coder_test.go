package interpreter_test

import (
	"math/big"
	"testing"

	s "github.com/kode4food/sputter/api"
	i "github.com/kode4food/sputter/interpreter"
	"github.com/stretchr/testify/assert"
)

func TestCreateCoder(t *testing.T) {
	a := assert.New(t)
	l := i.NewLexer("99")
	c := i.NewCoder(s.NewContext(), l)
	a.NotNil(c)
}

func TestCodeInteger(t *testing.T) {
	a := assert.New(t)
	l := i.NewLexer("99")
	c := i.NewCoder(s.NewContext(), l)
	v := c.Next()
	f, ok := v.(*big.Float)
	a.True(ok)
	a.Equal(0, f.Cmp(big.NewFloat(99)))
}

func TestCodeList(t *testing.T) {
	a := assert.New(t)
	l := i.NewLexer(`(99 "hello" 55.12)`)
	c := i.NewCoder(s.NewContext(), l)
	v := c.Next()
	list, ok := v.(*s.List)
	a.True(ok)

	iter := list.Iterate()
	value, ok := iter.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value.(*big.Float)))

	value, ok = iter.Next()
	a.True(ok)
	a.Equal("hello", value)

	value, ok = iter.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(55.12).Cmp(value.(*big.Float)))

	value, ok = iter.Next()
	a.False(ok)
}

func TestCodeNestedList(t *testing.T) {
	a := assert.New(t)
	l := i.NewLexer(`(99 ("hello" "there") 55.12)`)
	c := i.NewCoder(s.NewContext(), l)
	v := c.Next()
	list, ok := v.(*s.List)
	a.True(ok)

	iter1 := list.Iterate()
	value, ok := iter1.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value.(*big.Float)))

	// get nested list
	value, ok = iter1.Next()
	a.True(ok)
	list2, ok := value.(*s.List)
	a.True(ok)

	// iterate over the rest of top-level list
	value, ok = iter1.Next()
	a.True(ok)
	a.Equal(0, big.NewFloat(55.12).Cmp(value.(*big.Float)))

	value, ok = iter1.Next()
	a.False(ok)

	// iterate over the nested list
	iter2 := list2.Iterate()
	value, ok = iter2.Next()
	a.True(ok)
	a.Equal("hello", value)

	value, ok = iter2.Next()
	a.True(ok)
	a.Equal("there", value)

	value, ok = iter2.Next()
	a.False(ok)
}

func TestUnclosedList(t *testing.T) {
	a := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			a.Equal(r, i.ListNotClosed, "unclosed list")
			return
		}
		a.Fail("unclosed list didn't panic")
	}()

	l := i.NewLexer(`(99 ("hello" "there") 55.12`)
	c := i.NewCoder(s.NewContext(), l)
	c.Next()
}

func TestLiteral(t *testing.T) {
	a := assert.New(t)

	l := i.NewLexer(`'99`)
	c := i.NewCoder(s.NewContext(), l)
	v := c.Next()

	literal, ok := v.(*s.Literal)
	a.True(ok)

	value, ok := literal.Value.(*big.Float)
	a.True(ok)
	a.Equal(0, big.NewFloat(99).Cmp(value))
}
