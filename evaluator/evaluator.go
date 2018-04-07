package evaluator

import (
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/compiler"
	"github.com/kode4food/sputter/reader"
)

type evaluatorFunc struct{}

var evaluator = &evaluatorFunc{}

// Evaluate a Sequence of Values that a call to Read might produce
func Evaluate(c a.Context, s a.Sequence) a.Sequence {
	return a.Map(c, compiler.Compile(c, s), evaluator)
}

// EvalStr evaluates the specified raw Source
func EvalStr(c a.Context, src a.Str) a.Value {
	r := reader.ReadStr(src)
	if v, ok := a.Last(Evaluate(c, r)); ok {
		return v
	}
	return a.Nil
}

// NewEvalContext creates a new Context instance that
// chains up to the UserDomain Context for special forms
func NewEvalContext() a.Context {
	ns := a.GetNamespace(a.UserDomain)
	return a.ChildLocals(ns)
}

func (*evaluatorFunc) Apply(c a.Context, args a.Vector) a.Value {
	return a.Eval(c, args[0])
}
