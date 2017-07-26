package builtins

import a "github.com/kode4food/sputter/api"

var macroMetadata = a.Properties{
	a.MacroKey: a.True,
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := parseNamedFunction(args)
	n := a.NewLocalSymbol(fd.name)
	f := makeFunction(c, fd)
	r := f.WithMetadata(macroMetadata)
	return def(c, a.NewVector(n, r))
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

func isMacro(v a.Value) bool {
	if ap, ok := v.(a.Applicable); ok {
		m, _ := a.IsMacro(ap)
		return m
	}
	return false
}

func init() {
	RegisterBuiltIn("defmacro", defmacro)
	RegisterBuiltIn("macroexpand1", macroexpand1)
	RegisterBuiltIn("macroexpand", macroexpand)
	RegisterBuiltIn("macroexpand-all", macroexpandAll)

	RegisterSequencePredicate("macro?", isMacro)
	RegisterSequencePredicate("special-form?", isSpecialForm)
}
