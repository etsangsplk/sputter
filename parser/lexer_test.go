package parser_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	p "github.com/kode4food/sputter/parser"
)

func makeToken(t p.TokenType, v a.Value) *p.Token {
	return &p.Token{
		Type:  t,
		Value: v,
	}
}

func assertToken(as *assert.Wrapper, like *p.Token, value *p.Token) {
	as.Number(float64(like.Type), float64(value.Type))
	if like.Type != p.EndOfFile {
		as.Equal(like.Value, value.Value)
	}
}

func assertTokenSequence(as *assert.Wrapper, s a.Sequence, tokens []*p.Token) {
	iter := a.Iterate(s)
	for _, l := range tokens {
		v, ok := iter.Next()
		as.True(ok)
		assertToken(as, l, v.(*p.Token))
	}
	v, ok := iter.Next()
	as.False(ok)
	as.Nil(v)
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer("hello")
	as.NotNil(l)
}

func TestWhitespace(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer("   \t ")
	assertTokenSequence(as, l, []*p.Token{})
}

func TestEmptyList(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(" ( \t ) ")
	assertTokenSequence(as, l, []*p.Token{
		makeToken(p.ListStart, s("(")),
		makeToken(p.ListEnd, s(")")),
	})
}

func TestNumbers(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertTokenSequence(as, l, []*p.Token{
		makeToken(p.Number, f(10)),
		makeToken(p.Number, f(12.8)),
		makeToken(p.Number, f(8E+10)),
		makeToken(p.Number, f(99.598e+10)),
		makeToken(p.Number, f(54e+12)),
		makeToken(p.Ratio, a.NewRatio(1, 2)),
	})
}

func TestStrings(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertTokenSequence(as, l, []*p.Token{
		makeToken(p.String, s(`hello there`)),
		makeToken(p.String, s(`how's "life"?`)),
	})
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)
	l := p.NewLexer(` "hello there"
  "how's life?"
99`)

	assertTokenSequence(as, l, []*p.Token{
		makeToken(p.String, s(`hello there`)),
		makeToken(p.String, s(`how's life?`)),
		makeToken(p.Number, f(99)),
	})
}
