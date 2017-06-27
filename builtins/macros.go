package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	ex := a.MacroExpandAll(closure, d.body).(a.Sequence)
	db := a.NewBlock(ex)

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		l := ap(c, args)
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

func macroexpandAll(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	return a.MacroExpandAll(c, args.First())
}

func init() {
	registerAnnotated(
		a.NewFunction(defmacro).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("defmacro"),
			a.MetaDoc:     d.Get("defmacro"),
			a.MetaSpecial: a.True,
		}),
	)

	registerAnnotated(
		a.NewFunction(macroexpand1).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("macroexpand1"),
			a.MetaSpecial: a.True,
		}),
	)

	registerAnnotated(
		a.NewFunction(macroexpand).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("macroexpand"),
			a.MetaSpecial: a.True,
		}),
	)

	registerAnnotated(
		a.NewFunction(macroexpandAll).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("macroexpand-all"),
			a.MetaSpecial: a.True,
		}),
	)
}
