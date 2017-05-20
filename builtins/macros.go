package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	db := e.ExpandSequence(closure, d.body)

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args)
		ev := a.EvalBlock(l, db)
		ex := e.Expand(l, ev)
		return ex
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
		}),
	)
}
