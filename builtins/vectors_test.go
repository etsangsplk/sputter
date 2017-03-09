package builtins_test

import (
	"testing"
	
	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestNewVector(t *testing.T) {
	as := assert.New(t)
	r := runCode(`(vector 1 (- 5 3) (+ 1 2))`)
	as.Equal("[1 2 3]", a.String(r), "correct vector")
}
