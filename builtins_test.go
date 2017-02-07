package main_test

import (
	"testing"

	s "github.com/kode4food/sputter"
	"github.com/stretchr/testify/assert"
)

func TestBuiltInsContext(t *testing.T) {
	a := assert.New(t)

	bg1 := s.BuiltIns.Child()
	bg2 := bg1.Child()
	bg3 := bg2.Child()

	a.Equal(s.BuiltIns, bg3.Globals())
}
