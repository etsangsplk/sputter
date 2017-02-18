package builtins_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
	"github.com/stretchr/testify/assert"
)

func testCode(t *testing.T, src string, expect s.Value) {
	a := assert.New(t)
	l := r.NewLexer(src)
	c := r.NewCoder(b.Context, l)
	ctx := b.Context.Child()
	a.Equal(expect, r.EvalCoder(ctx, c), src)
}

func TestBuiltInsContext(t *testing.T) {
	a := assert.New(t)

	bg1 := b.Context.Child()
	bg2 := bg1.Child()
	bg3 := bg2.Child()

	a.Equal(b.Context, bg3.Globals())

	tv, ok := bg3.Get("true")
	a.True(ok)
	a.Equal(s.True, tv)
}

func TestTrueFalse(t *testing.T) {
	testCode(t, `true`, s.True)
	testCode(t, `false`, s.False)
	testCode(t, `nil`, s.Nil)
}
