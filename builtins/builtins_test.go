package builtins_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	"github.com/stretchr/testify/assert"
)

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
