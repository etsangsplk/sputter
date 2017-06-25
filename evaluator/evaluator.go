package evaluator

import a "github.com/kode4food/sputter/api"

// ReadStr converts the raw source into unexpanded data structures
func ReadStr(c a.Context, src a.Str) a.Sequence {
	l := Scan(src)
	return Read(c, l)
}

// EvalStr evaluates the specified raw Source
func EvalStr(c a.Context, src a.Str) a.Value {
	r := ReadStr(c, src)
	return a.Eval(c, r)
}

// NewEvalContext creates a new Context instance that
// chains up to the UserDomain Context for special forms
func NewEvalContext() a.Context {
	ns := a.GetNamespace(a.UserDomain)
	c := a.ChildContext(ns)
	c.Put(a.ContextDomain, ns)
	return c
}
