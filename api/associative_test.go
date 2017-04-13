package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func getTestMap() a.Associative {
	return a.Associative{
		a.Vector{a.NewKeyword("name"), "Sputter"},
		a.Vector{a.NewKeyword("age"), a.NewFloat(99)},
		a.Vector{"string", "value"},
	}
}

func TestAssociative(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	as.Equal(3, a.Count(m1), "count works")

	nameKey := a.NewKeyword("name")
	as.Equal(a.Name("name"), nameKey.Name(), "Name() works")

	nameValue, ok := m1.Get(nameKey)
	as.True(ok, "get works")
	as.Equal("Sputter", nameValue, "get works")

	ageKey := a.NewKeyword("age")
	ageValue, ok := m1.Get(ageKey)
	as.True(ok, "get works")
	as.Equal(a.EqualTo, a.NewFloat(99).Cmp(ageValue.(*a.Number)), "get works")

	strValue, ok := m1.Get("string")
	as.True(ok, "get works")
	as.Equal("value", strValue, "get works")

	r, ok := m1.Get("missing")
	as.False(ok, "miss works")
	as.Equal(a.Nil, r, "miss works")

	c := a.NewContext()
	as.Equal("Sputter", m1.Apply(c, a.NewList(nameKey)))
}

func TestAssociativeSequence(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	first := m1.First()
	if v, ok := first.(a.Vector); ok {
		as.Equal(a.NewKeyword("name"), v[0], "pair is good")
		as.Equal("Sputter", v[1], "pair is good")
	} else {
		as.Fail("map.First() is not a vector")
	}

	rest := m1.Rest()
	as.Equal(`{:age 99, "string" "value"}`, a.String(rest), "string works")

}

func TestAssociativePrepend(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	m2 := m1.Prepend(a.Vector{a.NewKeyword("foo"), "bar"}).(a.Associative)
	as.NotEqual(m1, m2, "prepended map not the same")

	r, ok := m2.Get(a.NewKeyword("foo"))
	as.True(ok, "prepended get works")
	as.Equal("bar", r, "prepended get works")

	if e2, ok := a.Eval(a.NewContext(), m2).(a.Associative); ok {
		as.True(&e2 != &m2, "evaluated map not the same")
	} else {
		as.Fail("map.Eval() didn't return an Associative")
	}

	defer expectError(as, a.ExpectedPair)
	m2.Conjoin(99)
}

func TestAssociativeIterate(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	i := a.Iterate(m1)
	if v, ok := i.Next(); ok {
		vec := v.(a.Vector)
		as.Equal(a.NewKeyword("name"), vec[0], "key")
		as.Equal("Sputter", vec[1], "value")
	} else {
		as.Fail("couldn't get first element")
	}

	if v, ok := i.Next(); ok {
		vec := v.(a.Vector)
		as.Equal(a.NewKeyword("age"), vec[0], "key")
		as.Equal(a.EqualTo, a.NewFloat(99).Cmp(vec[1].(*a.Number)), "value")
	} else {
		as.Fail("couldn't get second element")
	}
}

func TestAssociativeLookup(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("name")
	c := a.NewContext()
	args := a.NewList(m1)
	as.Equal("Sputter", nameKey.Apply(c, args), "get works")

	defer expectError(as, a.Err(a.ExpectedMapped, "99"))
	nameKey.Apply(c, a.NewList(a.NewFloat(99)))
}

func TestAssociativeMiss(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("miss")
	c := a.NewContext()

	defer expectError(as, a.Err(a.KeyNotFound, a.String(nameKey)))
	m1.Apply(c, a.NewList(nameKey))
}

func TestKeywordMiss(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("miss")
	c := a.NewContext()

	defer expectError(as, a.Err(a.KeyNotFound, a.String(nameKey)))
	nameKey.Apply(c, a.NewList(m1))
}

func TestAssertMapped(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()
	a.AssertMapped(m1)

	defer expectError(as, a.Err(a.ExpectedMapped, "99"))
	a.AssertMapped(a.NewFloat(99))
}
