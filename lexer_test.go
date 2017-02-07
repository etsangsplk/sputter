package main

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertToken(a *assert.Assertions, like *Token, value *Token) {
	a.Equal(like.Type, value.Type)
	switch like.Type {
	case Number:
		a.Equal(0, like.Value.(*big.Float).Cmp(value.Value.(*big.Float)))
	case Ratio:
		a.Equal(0, like.Value.(*big.Rat).Cmp(value.Value.(*big.Rat)))
	default:
		a.EqualValues(like.Value, value.Value)
	}
}

func TestCreateLexer(t *testing.T) {
	a := assert.New(t)
	lexer := NewLexer("hello")
	a.NotNil(lexer)
}

func TestWhitespace(t *testing.T) {
	a := assert.New(t)
	lexer := NewLexer("   \t ")
	assertToken(a, EOFToken, lexer.Next())
}

func TestEmptyList(t *testing.T) {
	a := assert.New(t)
	lexer := NewLexer(" ( \t ) ")
	assertToken(a, &Token{ListStart, "("}, lexer.Next())
	assertToken(a, &Token{ListEnd, ")"}, lexer.Next())
	assertToken(a, EOFToken, lexer.Next())
}

func TestNumbers(t *testing.T) {
	a := assert.New(t)
	lexer := NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(a, &Token{Number, big.NewFloat(10)}, lexer.Next())
	assertToken(a, &Token{Number, big.NewFloat(12.8)}, lexer.Next())
	assertToken(a, &Token{Number, big.NewFloat(8E+10)}, lexer.Next())
	assertToken(a, &Token{Number, big.NewFloat(99.598e+10)}, lexer.Next())
	assertToken(a, &Token{Number, big.NewFloat(54e+12)}, lexer.Next())
	assertToken(a, &Token{Ratio, big.NewRat(1, 2)}, lexer.Next())
	assertToken(a, EOFToken, lexer.Next())
}

func TestStrings(t *testing.T) {
	a := assert.New(t)
	lexer := NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(a, &Token{String, `hello there`}, lexer.Next())
	assertToken(a, &Token{String, `how's \"life\"?`}, lexer.Next())
	assertToken(a, EOFToken, lexer.Next())
}

func TestMultiline(t *testing.T) {
	a := assert.New(t)
	lexer := NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(a, &Token{String, `hello there`}, lexer.Next())
	assertToken(a, &Token{String, `how's life?`}, lexer.Next())
	assertToken(a, &Token{Number, big.NewFloat(99)}, lexer.Next())
	assertToken(a, EOFToken, lexer.Next())
}
