package builtins_test

import (
	"testing"

	b "github.com/kode4food/sputter/builtins"
	"github.com/stretchr/testify/assert"
)

func TestBuiltInsContext(t *testing.T) {
	a := assert.New(t)

	bg1 := b.BuiltIns.Child()
	bg2 := bg1.Child()
	bg3 := bg2.Child()

	a.Equal(b.BuiltIns, bg3.Globals())
}
