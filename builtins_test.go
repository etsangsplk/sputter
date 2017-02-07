package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltInsContext(t *testing.T) {
	a := assert.New(t)

	bg1 := BuiltIns.Child()
	bg2 := bg1.Child()
	bg3 := bg2.Child()

	a.Equal(BuiltIns, bg3.Globals())
}
