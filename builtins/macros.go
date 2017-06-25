package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	db := a.NewBlock(d.body)

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args, a.Identity)
		return a.Eval(l, db)
	}).WithMetadata(d.meta).(a.Function)
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	m := defineMacro(c, fd)
	a.GetContextNamespace(c).Put(m.Name(), m)
	return m
}

func macroexpand1(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	r, _ := a.MacroExpand1(c, args.First())
	return r
}

func macroexpand(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	r, _ := a.MacroExpand(c, args.First())
	return r
}

func init() {
	registerAnnotated(
		a.NewFunction(defmacro).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defmacro"),
			a.MetaDoc:  d.Get("defmacro"),
		}),
	)

	registerAnnotated(
		a.NewFunction(macroexpand1).WithMetadata(a.Metadata{
			a.MetaName: a.Name("macroexpand1"),
		}),
	)

	registerAnnotated(
		a.NewFunction(macroexpand).WithMetadata(a.Metadata{
			a.MetaName: a.Name("macroexpand"),
		}),
	)
}
