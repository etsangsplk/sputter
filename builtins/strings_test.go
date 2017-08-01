package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestReaderStr(t *testing.T) {
	as := assert.New(t)

	readerString := getBuiltIn("str!")
	v1 := a.NewVector(a.NewKeyword("boom"), a.Str("hello"))
	s1 := readerString.Apply(a.NewContext(), v1)
	as.String(":boom \"hello\"", s1)
}
