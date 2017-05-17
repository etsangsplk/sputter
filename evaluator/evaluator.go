package evaluator

import a "github.com/kode4food/sputter/api"

// NewEvaluator creates a new Evaluator using the given Context and raw source
func NewEvaluator(c a.Context, src a.Str) a.Sequence {
	l := NewLexer(src)
	r := NewReader(c, l)
	e := NewExpander(c, r)
	return e
}

// Eval evaluates a previously expanded Reader sequence
func Eval(c a.Context, s a.Sequence) a.Value {
	return a.EvalSequence(c, s)
}

// EvalStr evaluates the specified raw Source
func EvalStr(c a.Context, src a.Str) a.Value {
	e := NewEvaluator(c, src)
	return Eval(c, e)
}

// NewEvalContext creates a new Context instance that
// chains up to the UserDomain Context for special forms
func NewEvalContext() a.Context {
	ns := a.GetNamespace(a.UserDomain)
	c := a.ChildContext(ns)
	c.Put(a.ContextDomain, ns)
	return c
}
