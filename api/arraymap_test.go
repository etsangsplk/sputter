package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func getTestMap() a.ArrayMap {
	return a.ArrayMap{
		a.Vector{a.NewKeyword("name"), "Sputter"},
		a.Vector{a.NewKeyword("age"), a.NewFloat(99)},
		a.Vector{"string", "value"},
	}
}

func TestArrayMap(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	as.Equal(3, a.Count(m1), "count works")

	nameKey := a.NewKeyword("name")
	as.Equal(a.Name("name"), nameKey.Name(), "Name() works")
	as.Equal("Sputter", m1.Get(nameKey), "get works")

	ageKey := a.NewKeyword("age")
	ageValue := m1.Get(ageKey)
	as.Equal(a.EqualTo, a.NewFloat(99).Cmp(ageValue.(*a.Number)), "get works")

	strValue := m1.Get("string")
	as.Equal("value", strValue, "get works")

	as.Equal(a.Nil, m1.Get("missing"), "miss works")
}

func TestArrayMapSequence(t *testing.T) {
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

func TestArrayMapPrepend(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	m2 := m1.Prepend(a.Vector{a.NewKeyword("foo"), "bar"}).(a.ArrayMap)
	as.NotEqual(m1, m2, "prepended map not the same")
	as.Equal("bar", m2.Get(a.NewKeyword("foo")), "prepended get works")

	if e2, ok := a.Eval(a.NewContext(), m2).(a.ArrayMap); ok {
		as.True(&e2 != &m2, "evaluated map not the same")
	} else {
		as.Fail("map.Eval() didn't return an ArrayMap")
	}

	defer expectError(as, a.ExpectedPair)
	m2.Prepend(99)
}

func TestArrayMapIterate(t *testing.T) {
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

func TestArrayMapLookup(t *testing.T) {
	as := assert.New(t)
	m1 := getTestMap()

	nameKey := a.NewKeyword("name")
	c := a.NewContext()
	args := a.NewList(m1)
	as.Equal("Sputter", nameKey.Apply(c, args), "get works")

	defer expectError(as, a.ExpectedMapped)
	nameKey.Apply(c, a.NewList(99))
}
