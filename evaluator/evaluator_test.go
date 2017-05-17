package evaluator_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	e "github.com/kode4food/sputter/evaluator"
)

func TestEvalContext(t *testing.T) {
	as := assert.New(t)

	uc := a.GetNamespace(a.UserDomain)
	uc.Delete("foo")
	uc.Put("foo", f(99))

	ec := e.NewEvalContext()
	v, _ := ec.Get("foo")
	as.Number(99, v)
}
