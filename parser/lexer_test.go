package parser_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	p "github.com/kode4food/sputter/parser"
)

func makeToken(t p.TokenType, v a.Value) *p.Token {
	return &p.Token{Type: t, Value: v}
}

func assertToken(as *assert.Wrapper, like *p.Token, value *p.Token) {
	as.Equal(like.Type, value.Type)
	if like.Type != p.EndOfFile {
		as.Equal(like.Value, value.Value)
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
	assertToken(as, makeToken(p.ListStart, a.Atom("(")), l.Next())
	assertToken(as, makeToken(p.ListEnd, a.Atom(")")), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}

func TestNumbers(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(as, makeToken(p.Number, f(10)), l.Next())
	assertToken(as, makeToken(p.Number, f(12.8)), l.Next())
	assertToken(as, makeToken(p.Number, f(8E+10)), l.Next())
	assertToken(as, makeToken(p.Number, f(99.598e+10)), l.Next())
	assertToken(as, makeToken(p.Number, f(54e+12)), l.Next())
	assertToken(as, makeToken(p.Ratio, a.NewRatio(1, 2)), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}

func TestStrings(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(as, makeToken(p.String, a.Atom(`hello there`)), l.Next())
	assertToken(as, makeToken(p.String, a.Atom(`how's "life"?`)), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(as, makeToken(p.String, a.Atom(`hello there`)), l.Next())
	assertToken(as, makeToken(p.String, a.Atom(`how's life?`)), l.Next())
	assertToken(as, makeToken(p.Number, f(99)), l.Next())
	assertToken(as, p.EOFToken, l.Next())
}
