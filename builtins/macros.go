package builtins

import a "github.com/kode4food/sputter/api"

var macroMetadata = a.Properties{
	a.MetaMacro: a.True,
}

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	ap := makeArgProcessor(closure, d.args)
	md := d.meta.Child(macroMetadata)
	ex := a.MacroExpandAll(closure, d.body).(a.Sequence)
	db := a.NewBlock(ex)

	return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		return a.Eval(ap(c, args), db)
	}).WithMetadata(md).(a.Function)
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
	RegisterBuiltIn("defmacro", defmacro)
	RegisterBuiltIn("macroexpand1", macroexpand1)
	RegisterBuiltIn("macroexpand", macroexpand)
	RegisterBuiltIn("macroexpand-all", macroexpandAll)
}
