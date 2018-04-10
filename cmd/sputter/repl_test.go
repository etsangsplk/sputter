package main_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/builtins"
	c "github.com/kode4food/sputter/cmd/sputter"
	e "github.com/kode4food/sputter/evaluator"
)

func TestREPL(t *testing.T) {
	as := assert.New(t)

	r := c.NewREPL()
	as.NotNil(r)
}

func asApplicable(as *assert.Wrapper, v a.Value) a.Applicable {
	if r, ok := v.(a.Applicable); ok {
		return r
	}
	as.Fail("value is not Applicable")
	return nil
}

func TestBuiltInUse(t *testing.T) {
	as := assert.New(t)

	ec := e.NewEvalContext()
	v, ok := ec.Get("use")
	as.True(ok)
	as.NotNil(v)
	ap := asApplicable(as, v)
	nsName := a.NewLocalSymbol(a.Name("test-ns"))
	ns, ok := ap.Apply(ec, a.Vector{nsName}).(a.Namespace)
	as.True(ok)
	as.String("test-ns", ns.Domain())

	ns = a.GetContextNamespace(ec)
	as.String("test-ns", ns.Domain())

	ec.Delete(a.ContextDomain)
}
