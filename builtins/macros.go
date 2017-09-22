package builtins

import a "github.com/kode4food/sputter/api"

const (
	defMacroName       = "defmacro"
	macroExpand1Name   = "macroexpand1"
	macroExpandName    = "macroexpand"
	macroExpandAllName = "macroexpand-all"
	isMacroName        = "macro?"
)

type (
	defMacroFunction  struct{ BaseBuiltIn }
	expand1Function   struct{ BaseBuiltIn }
	expandFunction    struct{ BaseBuiltIn }
	expandAllFunction struct{ BaseBuiltIn }
)

var macroMetadata = a.Properties{
	a.MacroKey: a.True,
}

func (*defMacroFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fd := parseNamedFunction(args)
	n := a.NewLocalSymbol(fd.name)
	fn := makeFunction(c, fd)
	r := fn.WithMetadata(macroMetadata)
	return new(defFunction).Apply(c, a.Values{n, r})
}

func (*expand1Function) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	r, _ := a.MacroExpand1(c, args.First())
	return r
}

func (*expandFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	r, _ := a.MacroExpand(c, args.First())
	return r
}

func (*expandAllFunction) Apply(c a.Context, args a.Sequence) a.Value {
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

	RegisterBuiltIn(defMacroName, defMacro)
	RegisterBuiltIn(macroExpand1Name, macroExpand1)
	RegisterBuiltIn(macroExpandName, macroExpand)
	RegisterBuiltIn(macroExpandAllName, macroExpandAll)

	RegisterSequencePredicate(isMacroName, isMacro)
}
