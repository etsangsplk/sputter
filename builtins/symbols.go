package builtins

import a "github.com/kode4food/sputter/api"

const genSymName = "gensym"

type genSymFunction struct{ BaseBuiltIn }

var anonName = a.Name("anon")

func (*genSymFunction) Apply(c a.Context, args a.Sequence) a.Value {
	if a.AssertArityRange(args, 0, 1) == 1 {
		return a.NewGeneratedSymbol(a.Name(args.First().(a.Str)))
	}
	return a.NewGeneratedSymbol(anonName)
}

func init() {
	var genSym *genSymFunction

	RegisterBuiltIn(genSymName, genSym)
}
