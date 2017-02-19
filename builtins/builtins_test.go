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
	a.Equal(expect, r.EvalCoder(s.NewContext(), c), src)
}

func testBadCode(t *testing.T, src string, err string) {
	a := assert.New(t)

	defer func() {
		if rec := recover(); rec != nil {
			a.Equal(err, rec, "bad code panics properly")
			return
		}
		a.Fail("bad code should panic")
	}()

	l := r.NewLexer(src)
	c := r.NewCoder(b.Context, l)
	r.EvalCoder(s.NewContext(), c)
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
