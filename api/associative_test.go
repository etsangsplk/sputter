package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func getTestMap() a.Associative {
	r := a.NewAssociative(
		a.NewVector(a.NewKeyword("name"), s("Sputter")),
		a.NewVector(a.NewKeyword("age"), f(99)),
		a.NewVector(s("string"), s("value")),
	)
	return a.Eval(a.Variables{}, r).(a.Associative)
}

func TestAssociative(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	as.Number(3, a.Count(m1))

	nameKey := a.NewKeyword("name")
	as.Equal(a.Name("name"), nameKey.Name())

	nameValue, ok := m1.Get(nameKey)
	as.True(ok)
	as.String("Sputter", nameValue)

	ageKey := a.NewKeyword("age")
	ageValue, ok := m1.Get(ageKey)
	as.True(ok)
	as.Number(99, ageValue)

	strValue, ok := m1.Get(s("string"))
	as.True(ok)
	as.String("value", strValue)

	r, ok := m1.Get(s("missing"))
	as.False(ok)
	as.Equal(a.Nil, r)

	c := a.Variables{}
	as.String("Sputter", m1.Apply(c, a.NewList(nameKey)))
}

func TestAssociativeSequence(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	first := m1.First()
	if e, ok := first.(a.Vector); ok {
		k, _ := e.ElementAt(0)
		v, _ := e.ElementAt(1)
		as.Equal(a.NewKeyword("name"), k)
		as.String("Sputter", v)
	} else {
		as.Fail("map.First() is not a vector")
	}

	rest := m1.Rest()
	as.String(`{:age 99, "string" "value"}`, rest)

}

func TestAssociativePrepend(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	m2 := m1.Prepend(a.NewVector(a.NewKeyword("foo"), s("bar"))).(a.Associative)
	as.NotIdentical(m1, m2)

	r, ok := m2.Get(a.NewKeyword("foo"))
	as.True(ok)
	as.String("bar", r)

	if e2, ok := a.Eval(a.Variables{}, m2).(a.Associative); ok {
		as.True(&e2 != &m2)
	} else {
		as.Fail("map.Eval() didn't return an Associative")
	}

	defer as.ExpectError(a.ErrStr(a.ExpectedPair))
	m2.Conjoin(f(99))
}

func TestAssociativeIterate(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	i := a.Iterate(m1)
	if v, ok := i.Next(); ok {
		vec := v.(a.Vector)
		key, _ := vec.ElementAt(0)
		val, _ := vec.ElementAt(1)
		as.Equal(a.NewKeyword("name"), key)
		as.String("Sputter", val)
	} else {
		as.Fail("couldn't get first element")
	}

	if v, ok := i.Next(); ok {
		vec := v.(a.Vector)
		key, _ := vec.ElementAt(0)
		val, _ := vec.ElementAt(1)
		as.Equal(a.NewKeyword("age"), key)
		as.Number(99, val)
	} else {
		as.Fail("couldn't get second element")
	}
}

func TestAssociativeApply(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	c := a.Variables{}
	nameKey := a.NewKeyword("name")
	args := a.NewList(nameKey)
	as.String("Sputter", m1.Apply(c, args))

	missKey := a.NewKeyword("miss")
	args = a.NewList(missKey, s("you missed"))
	as.String("you missed", m1.Apply(c, args))

	defer as.ExpectError(a.ErrStr(a.KeyNotFound, missKey))
	args = a.NewList(missKey)
	m1.Apply(c, args)
}

func TestAssociativeLookup(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("name")
	c := a.Variables{}
	args := a.NewList(m1)
	as.String("Sputter", nameKey.Apply(c, args))

	defer as.ExpectError(cvtErr("*api.dec", "api.Mapped", "Get"))
	nameKey.Apply(c, a.NewList(f(99)))
}

func TestAssociativeMiss(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("miss")
	c := a.Variables{}

	defer as.ExpectError(a.ErrStr(a.KeyNotFound, nameKey))
	m1.Apply(c, a.NewList(nameKey))
}

func TestKeywordMiss(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("miss")
	c := a.Variables{}

	defer as.ExpectError(a.ErrStr(a.KeyNotFound, nameKey))
	nameKey.Apply(c, a.NewList(m1))
}
