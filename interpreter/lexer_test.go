package interpreter_test

import (
	"math/big"
	"testing"

	i "github.com/kode4food/sputter/interpreter"
	"github.com/stretchr/testify/assert"
)

func assertToken(a *assert.Assertions, like *i.Token, value *i.Token) {
	a.Equal(like.Type, value.Type)
	switch like.Type {
	case i.Number:
		a.Equal(0, like.Value.(*big.Float).Cmp(value.Value.(*big.Float)))
	case i.Ratio:
		a.Equal(0, like.Value.(*big.Rat).Cmp(value.Value.(*big.Rat)))
	default:
		a.EqualValues(like.Value, value.Value)
	}
}

func TestCreateLexer(t *testing.T) {
	a := assert.New(t)
	lexer := i.NewLexer("hello")
	a.NotNil(lexer)
}

func TestWhitespace(t *testing.T) {
	a := assert.New(t)
	lexer := i.NewLexer("   \t ")
	assertToken(a, i.EOFToken, lexer.Next())
}

func TestEmptyList(t *testing.T) {
	a := assert.New(t)
	lexer := i.NewLexer(" ( \t ) ")
	assertToken(a, &i.Token{i.ListStart, "("}, lexer.Next())
	assertToken(a, &i.Token{i.ListEnd, ")"}, lexer.Next())
	assertToken(a, i.EOFToken, lexer.Next())
}

func TestNumbers(t *testing.T) {
	a := assert.New(t)
	lexer := i.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(a, &i.Token{i.Number, big.NewFloat(10)}, lexer.Next())
	assertToken(a, &i.Token{i.Number, big.NewFloat(12.8)}, lexer.Next())
	assertToken(a, &i.Token{i.Number, big.NewFloat(8E+10)}, lexer.Next())
	assertToken(a, &i.Token{i.Number, big.NewFloat(99.598e+10)}, lexer.Next())
	assertToken(a, &i.Token{i.Number, big.NewFloat(54e+12)}, lexer.Next())
	assertToken(a, &i.Token{i.Ratio, big.NewRat(1, 2)}, lexer.Next())
	assertToken(a, i.EOFToken, lexer.Next())
}

func TestStrings(t *testing.T) {
	a := assert.New(t)
	lexer := i.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(a, &i.Token{i.String, `hello there`}, lexer.Next())
	assertToken(a, &i.Token{i.String, `how's \"life\"?`}, lexer.Next())
	assertToken(a, i.EOFToken, lexer.Next())
}

func TestMultiline(t *testing.T) {
	a := assert.New(t)
	lexer := i.NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(a, &i.Token{i.String, `hello there`}, lexer.Next())
	assertToken(a, &i.Token{i.String, `how's life?`}, lexer.Next())
	assertToken(a, &i.Token{i.Number, big.NewFloat(99)}, lexer.Next())
	assertToken(a, i.EOFToken, lexer.Next())
}
