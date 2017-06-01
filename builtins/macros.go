package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	e "github.com/kode4food/sputter/evaluator"
)

func evalExpandBlock(c a.Context, s a.Sequence) a.Value {
	ev := a.EvalBlock(c, s)
	return e.Expand(c, ev)
}

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	db := e.ExpandSequence(closure, d.body)

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args)
		return evalExpandBlock(l, db)
	}).WithMetadata(d.meta).(a.Function)
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	m := defineMacro(c, fd)
	a.GetContextNamespace(c).Put(m.Name(), m)
	return m
}

func init() {
	registerAnnotated(
		a.NewMacro(defmacro).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defmacro"),
			a.MetaDoc:  d.Get("defmacro"),
		}),
	)
}
