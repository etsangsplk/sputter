package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func getTestMap() a.Associative {
	return a.Associative{
		a.Vector{a.NewKeyword("name"), s("Sputter")},
		a.Vector{a.NewKeyword("age"), f(99)},
		a.Vector{s("string"), s("value")},
	}
}

func TestAssociative(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	as.Equal(3, a.Count(m1))

	nameKey := a.NewKeyword("name")
	as.Equal(a.Name("name"), nameKey.Name())

	nameValue, ok := m1.Get(nameKey)
	as.True(ok)
	as.String("Sputter", nameValue)

	ageKey := a.NewKeyword("age")
	ageValue, ok := m1.Get(ageKey)
	as.True(ok)
	as.Float(99, ageValue)

	strValue, ok := m1.Get(s("string"))
	as.True(ok)
	as.String("value", strValue)

	r, ok := m1.Get(s("missing"))
	as.False(ok)
	as.Equal(a.Nil, r)

	c := a.NewContext()
	as.String("Sputter", m1.Apply(c, a.NewList(nameKey)))
}

func TestAssociativeSequence(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	first := m1.First()
	if v, ok := first.(a.Vector); ok {
		as.Equal(a.NewKeyword("name"), v[0])
		as.String("Sputter", v[1])
	} else {
		as.Fail("map.First() is not a vector")
	}

	rest := m1.Rest()
	as.String(`{:age 99, "string" "value"}`, rest)

}

func TestAssociativePrepend(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	m2 := m1.Prepend(a.Vector{a.NewKeyword("foo"), s("bar")}).(a.Associative)
	as.NotIdentical(m1, m2)

	r, ok := m2.Get(a.NewKeyword("foo"))
	as.True(ok)
	as.String("bar", r)

	if e2, ok := a.Eval(a.NewContext(), m2).(a.Associative); ok {
		as.True(&e2 != &m2)
	} else {
		as.Fail("map.Eval() didn't return an Associative")
	}

	defer expectError(as, a.ExpectedPair)
	m2.Conjoin(f(99))
}

func TestAssociativeIterate(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	i := a.Iterate(m1)
	if v, ok := i.Next(); ok {
		vec := v.(a.Vector)
		as.Equal(a.NewKeyword("name"), vec[0])
		as.String("Sputter", vec[1])
	} else {
		as.Fail("couldn't get first element")
	}

	if v, ok := i.Next(); ok {
		vec := v.(a.Vector)
		as.Equal(a.NewKeyword("age"), vec[0])
		as.Float(99, vec[1])
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
	as.String("Sputter", nameKey.Apply(c, args))

	defer expectError(as, a.Err(a.ExpectedMapped, f(99)))
	nameKey.Apply(c, a.NewList(f(99)))
}

func TestAssociativeMiss(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("miss")
	c := a.NewContext()

	defer expectError(as, a.Err(a.KeyNotFound, nameKey))
	m1.Apply(c, a.NewList(nameKey))
}

func TestKeywordMiss(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("miss")
	c := a.NewContext()

	defer expectError(as, a.Err(a.KeyNotFound, nameKey))
	nameKey.Apply(c, a.NewList(m1))
}

func TestAssertMapped(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()
	a.AssertMapped(m1)

	defer expectError(as, a.Err(a.ExpectedMapped, f(99)))
	a.AssertMapped(f(99))
}
