package reader_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func makeToken(t r.TokenType, v a.Value) *r.Token {
	return &r.Token{Type: t, Value: v}
}

func assertToken(as *assert.Assertions, like *r.Token, value *r.Token) {
	as.Equal(like.Type, value.Type)
	switch like.Type {
	case r.Number:
		as.Equal(a.EqualTo, like.Value.(*a.Number).Cmp(value.Value.(*a.Number)))
	case r.Ratio:
		as.Equal(a.EqualTo, like.Value.(*a.Number).Cmp(value.Value.(*a.Number)))
	default:
		as.EqualValues(like.Value, value.Value)
	}
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer("hello")
	as.NotNil(l)
}

func TestWhitespace(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer("   \t ")
	assertToken(as, r.EOFToken, l.Next())
}

func TestEmptyList(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(" ( \t ) ")
	assertToken(as, makeToken(r.ListStart, "("), l.Next())
	assertToken(as, makeToken(r.ListEnd, ")"), l.Next())
	assertToken(as, r.EOFToken, l.Next())
}

func TestNumbers(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(" 10 12.8 8E+10 99.598e+10 54e+12 1/2")
	assertToken(as, makeToken(r.Number, a.NewFloat(10)), l.Next())
	assertToken(as, makeToken(r.Number, a.NewFloat(12.8)), l.Next())
	assertToken(as, makeToken(r.Number, a.NewFloat(8E+10)), l.Next())
	assertToken(as, makeToken(r.Number, a.NewFloat(99.598e+10)), l.Next())
	assertToken(as, makeToken(r.Number, a.NewFloat(54e+12)), l.Next())
	assertToken(as, makeToken(r.Ratio, a.NewRatio(1, 2)), l.Next())
	assertToken(as, r.EOFToken, l.Next())
}

func TestStrings(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(` "hello there" "how's \"life\"?"  `)
	assertToken(as, makeToken(r.String, `hello there`), l.Next())
	assertToken(as, makeToken(r.String, `how's \"life\"?`), l.Next())
	assertToken(as, r.EOFToken, l.Next())
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)
	l := r.NewLexer(` "hello there"
  "how's life?"
99`)

	assertToken(as, makeToken(r.String, `hello there`), l.Next())
	assertToken(as, makeToken(r.String, `how's life?`), l.Next())
	assertToken(as, makeToken(r.Number, a.NewFloat(99)), l.Next())
	assertToken(as, r.EOFToken, l.Next())
}
