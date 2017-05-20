package evaluator

import a "github.com/kode4food/sputter/api"

// NewEvaluator creates a new Evaluator using the given Context and raw source
func NewEvaluator(c a.Context, src a.Str) a.Sequence {
	l := NewLexer(src)
	r := NewReader(c, l)
	e := Expand(c, r).(a.Sequence)
	return e
}

// Read converts the raw source into unexpanded data structures
func Read(c a.Context, src a.Str) a.Sequence {
	l := NewLexer(src)
	return NewReader(c, l)
}

// Eval evaluates data structures that have not yet been expanded
func Eval(c a.Context, s a.Sequence) a.Value {
	ex := Expand(c, s).(a.Sequence)
	return a.EvalBlock(c, ex)
}

// EvalStr evaluates the specified raw Source
func EvalStr(c a.Context, src a.Str) a.Value {
	r := Read(c, src)
	return Eval(c, r)
}

// NewEvalContext creates a new Context instance that
// chains up to the UserDomain Context for special forms
func NewEvalContext() a.Context {
	ns := a.GetNamespace(a.UserDomain)
	c := a.ChildContext(ns)
	c.Put(a.ContextDomain, ns)
	return c
}
