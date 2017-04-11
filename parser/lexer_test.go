package parser_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	p "github.com/kode4food/sputter/parser"
	"github.com/stretchr/testify/assert"
)

func makeToken(t p.TokenType, v a.Value) *p.Token {
	return &p.Token{Type: t, Value: v}
}

func assertToken(as *assert.Assertions, like *p.Token, value *p.Token) {
	as.Equal(like.Type, value.Type)
	switch like.Type {
	case p.Number:
		as.Equal(a.EqualTo, like.Value.(*a.Number).Cmp(value.Value.(*a.Number)))
	case p.Ratio:
		as.Equal(a.EqualTo, like.Value.(*a.Number).Cmp(value.Value.(*a.Number)))
	default:
		as.EqualValues(like.Value, value.Value)
	}
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer("hello")
	as.NotNil(l)
}

func TestWhitespace(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer("   \t ")
	assertToken(as, p.EOFToken, l.Next())
}

func TestEmptyList(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(" ( \t ) ")
	assertToken(as, makeToken(p.ListStart, "("), l.Next())
	assertToken(as, makeToken(p.ListEnd, ")"), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}

func TestNumbers(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(as, makeToken(p.Number, a.NewFloat(10)), l.Next())
	assertToken(as, makeToken(p.Number, a.NewFloat(12.8)), l.Next())
	assertToken(as, makeToken(p.Number, a.NewFloat(8E+10)), l.Next())
	assertToken(as, makeToken(p.Number, a.NewFloat(99.598e+10)), l.Next())
	assertToken(as, makeToken(p.Number, a.NewFloat(54e+12)), l.Next())
	assertToken(as, makeToken(p.Ratio, a.NewRatio(1, 2)), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}

func TestStrings(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(as, makeToken(p.String, `hello there`), l.Next())
	assertToken(as, makeToken(p.String, `how's \"life\"?`), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(as, makeToken(p.String, `hello there`), l.Next())
	assertToken(as, makeToken(p.String, `how's life?`), l.Next())
	assertToken(as, makeToken(p.Number, a.NewFloat(99)), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}
