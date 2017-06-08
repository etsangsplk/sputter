package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	e "github.com/kode4food/sputter/evaluator"
)

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args)
		db := e.ExpandSequence(l, d.body)
		return e.EvalExpand(l, db)
	}).WithMetadata(d.meta).(a.Function)
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	m := defineMacro(c, fd)
	a.GetContextNamespace(c).Put(m.Name(), m)
	return m
}

func macroexpand(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	ev := a.NewBlock(args).Eval(c)
	return e.Expand(c, ev).Str()
}

func init() {
	registerAnnotated(
		a.NewMacro(defmacro).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defmacro"),
			a.MetaDoc:  d.Get("defmacro"),
		}),
	)

	registerAnnotated(
		a.NewMacro(macroexpand).WithMetadata(a.Metadata{
			a.MetaName: a.Name("macroexpand"),
		}),
	)
}
