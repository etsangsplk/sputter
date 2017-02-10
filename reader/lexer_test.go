package interpreter_test

import (
	"math/big"
	"testing"

	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func assertToken(a *assert.Assertions, like *r.Token, value *r.Token) {
	a.Equal(like.Type, value.Type)
	switch like.Type {
	case r.Number:
		a.Equal(0, like.Value.(*big.Float).Cmp(value.Value.(*big.Float)))
	case r.Ratio:
		a.Equal(0, like.Value.(*big.Rat).Cmp(value.Value.(*big.Rat)))
	default:
		a.EqualValues(like.Value, value.Value)
	}
}

func TestCreateLexer(t *testing.T) {
	a := assert.New(t)
	lexer := r.NewLexer("hello")
	a.NotNil(lexer)
}

func TestWhitespace(t *testing.T) {
	a := assert.New(t)
	lexer := r.NewLexer("   \t ")
	assertToken(a, r.EOFToken, lexer.Next())
}

func TestEmptyList(t *testing.T) {
	a := assert.New(t)
	lexer := r.NewLexer(" ( \t ) ")
	assertToken(a, &r.Token{r.ListStart, "("}, lexer.Next())
	assertToken(a, &r.Token{r.ListEnd, ")"}, lexer.Next())
	assertToken(a, r.EOFToken, lexer.Next())
}

func TestNumbers(t *testing.T) {
	a := assert.New(t)
	lexer := r.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(a, &r.Token{r.Number, big.NewFloat(10)}, lexer.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(12.8)}, lexer.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(8E+10)}, lexer.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(99.598e+10)}, lexer.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(54e+12)}, lexer.Next())
	assertToken(a, &r.Token{r.Ratio, big.NewRat(1, 2)}, lexer.Next())
	assertToken(a, r.EOFToken, lexer.Next())
}

func TestStrings(t *testing.T) {
	a := assert.New(t)
	lexer := r.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(a, &r.Token{r.String, `hello there`}, lexer.Next())
	assertToken(a, &r.Token{r.String, `how's \"life\"?`}, lexer.Next())
	assertToken(a, r.EOFToken, lexer.Next())
}

func TestMultiline(t *testing.T) {
	a := assert.New(t)
	lexer := r.NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(a, &r.Token{r.String, `hello there`}, lexer.Next())
	assertToken(a, &r.Token{r.String, `how's life?`}, lexer.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(99)}, lexer.Next())
	assertToken(a, r.EOFToken, lexer.Next())
}
