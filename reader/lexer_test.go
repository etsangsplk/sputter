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
	l := r.NewLexer("hello")
	a.NotNil(l)
}

func TestWhitespace(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer("   \t ")
	assertToken(a, r.EOFToken, l.Next())
}

func TestEmptyList(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(" ( \t ) ")
	assertToken(a, &r.Token{r.ListStart, "("}, l.Next())
	assertToken(a, &r.Token{r.ListEnd, ")"}, l.Next())
	assertToken(a, r.EOFToken, l.Next())
}

func TestNumbers(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(a, &r.Token{r.Number, big.NewFloat(10)}, l.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(12.8)}, l.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(8E+10)}, l.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(99.598e+10)}, l.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(54e+12)}, l.Next())
	assertToken(a, &r.Token{r.Ratio, big.NewRat(1, 2)}, l.Next())
	assertToken(a, r.EOFToken, l.Next())
}

func TestStrings(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(a, &r.Token{r.String, `hello there`}, l.Next())
	assertToken(a, &r.Token{r.String, `how's \"life\"?`}, l.Next())
	assertToken(a, r.EOFToken, l.Next())
}

func TestMultiline(t *testing.T) {
	a := assert.New(t)
	l := r.NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(a, &r.Token{r.String, `hello there`}, l.Next())
	assertToken(a, &r.Token{r.String, `how's life?`}, l.Next())
	assertToken(a, &r.Token{r.Number, big.NewFloat(99)}, l.Next())
	assertToken(a, r.EOFToken, l.Next())
}
