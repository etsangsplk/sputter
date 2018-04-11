package reader_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	r "github.com/kode4food/sputter/reader"
)

func makeToken(t r.TokenType, v a.Value) *r.Token {
	return &r.Token{
		Type:  t,
		Value: v,
	}
}

func assertToken(as *assert.Wrapper, like *r.Token, value *r.Token) {
	as.Number(float64(like.Type), float64(value.Type))
}

func assertTokenSequence(as *assert.Wrapper, s a.Sequence, tokens []*r.Token) {
	iter := a.Iterate(s)
	for _, l := range tokens {
		v, ok := iter.Next()
		as.True(ok)
		assertToken(as, l, v.(*r.Token))
	}
	v, ok := iter.Next()
	as.False(ok)
	as.Nil(v)
}

func TestCreateLexer(t *testing.T) {
	as := assert.New(t)
	l := r.Scan("hello")
	as.NotNil(l)
	as.String(`([1 "hello"])`, a.MakeSequenceStr(l))
}

func TestWhitespace(t *testing.T) {
	as := assert.New(t)
	l := r.Scan("   \t ")
	assertTokenSequence(as, l, []*r.Token{})
}

func TestEmptyList(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(" ( \t ) ")
	assertTokenSequence(as, l, []*r.Token{
		makeToken(r.ListStart, s("(")),
		makeToken(r.ListEnd, s(")")),
	})
}

func TestNumbers(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(" 10 12.8 8E+10 99.598e+10 54e+12 1/2 -0xFF  071 0xf1e9d8c7")
	assertTokenSequence(as, l, []*r.Token{
		makeToken(r.Number, f(10)),
		makeToken(r.Number, f(12.8)),
		makeToken(r.Number, f(8E+10)),
		makeToken(r.Number, f(99.598e+10)),
		makeToken(r.Number, f(54e+12)),
		makeToken(r.Ratio, a.NewRatio(1, 2)),
		makeToken(r.Number, f(-255)),
		makeToken(r.Number, f(57)),
		makeToken(r.Number, f(4058634439)),
	})
}

func TestStrings(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(` "hello there" "how's \"life\"?"  `)
	assertTokenSequence(as, l, []*r.Token{
		makeToken(r.String, s(`hello there`)),
		makeToken(r.String, s(`how's "life"?`)),
	})
}

func TestMultiLine(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(` "hello there"
  "how's life?"
99`)

	assertTokenSequence(as, l, []*r.Token{
		makeToken(r.String, s(`hello there`)),
		makeToken(r.String, s(`how's life?`)),
		makeToken(r.Number, f(99)),
	})
}

func TestComments(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(`"hello" ; (this is commented)`)
	assertTokenSequence(as, l, []*r.Token{
		makeToken(r.String, s(`hello`)),
	})
}
