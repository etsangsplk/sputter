package integration_test_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestStr(t *testing.T) {
	testCode(t, `
	  (str "hello" nil [1 2 3 4])
	`, s("hello[1 2 3 4]"))

	testCode(t, `
	  (str? "hello" "there")
	`, a.True)

	testCode(t, `
	  (str? "hello" 99)
	`, a.False)
}

func TestReadableStr(t *testing.T) {
	testCode(t, "(str! \"hello\nyou\")", s("\"hello\nyou\""))
	testCode(t, `(str! "hello" "you")`, s(`"hello" "you"`))
}
