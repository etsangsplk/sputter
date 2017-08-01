package builtins

import a "github.com/kode4food/sputter/api"

var macroMetadata = a.Properties{
	a.MacroKey: a.True,
}

type (
	defMacroFunction  struct{ a.ReflectedFunction }
	expand1Function   struct{ a.ReflectedFunction }
	expandFunction    struct{ a.ReflectedFunction }
	expandAllFunction struct{ a.ReflectedFunction }
)

func (f *defMacroFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fd := parseNamedFunction(args)
	n := a.NewLocalSymbol(fd.name)
	fn := makeFunction(c, fd)
	r := fn.WithMetadata(macroMetadata)
	return new(defFunction).Apply(c, a.NewVector(n, r))
}

func (f *expand1Function) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	r, _ := a.MacroExpand1(c, args.First())
	return r
}

func (f *expandFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	r, _ := a.MacroExpand(c, args.First())
	return r
}

func (f *expandAllFunction) Apply(c a.Context, args a.Sequence) a.Value {
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
	var defMacro *defMacroFunction
	var macroExpand1 *expand1Function
	var macroExpand *expandFunction
	var macroExpandAll *expandAllFunction

	RegisterBaseFunction("defmacro", defMacro)
	RegisterBaseFunction("macroexpand1", macroExpand1)
	RegisterBaseFunction("macroexpand", macroExpand)
	RegisterBaseFunction("macroexpand-all", macroExpandAll)

	RegisterSequencePredicate("macro?", isMacro)
	RegisterSequencePredicate("special-form?", isSpecialForm)
}
