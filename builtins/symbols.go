package builtins

import a "github.com/kode4food/sputter/api"

const (
	symName      = "sym"
	genSymName   = "gensym"
	isSymbolName = "is-symbol"
	isLocalName  = "is-local"
)

type (
	symFunction      struct{ BaseBuiltIn }
	genSymFunction   struct{ BaseBuiltIn }
	isSymbolFunction struct{ BaseBuiltIn }
	isLocalFunction  struct{ BaseBuiltIn }
)

var anonName = a.Name("anon")

func (*symFunction) Apply(c a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 1)
	return a.ParseSymbol(a.Name(args[0].(a.Str)))
}

func (*genSymFunction) Apply(c a.Context, args a.Vector) a.Value {
	if a.AssertArityRange(args, 0, 1) == 1 {
		return a.NewGeneratedSymbol(a.Name(args[0].(a.Str)))
	}
	return a.NewGeneratedSymbol(anonName)
}

func (*isSymbolFunction) Apply(_ a.Context, args a.Vector) a.Value {
	_, ok := args[0].(a.Symbol)
	return a.Bool(ok)
}

func (*isLocalFunction) Apply(_ a.Context, args a.Vector) a.Value {
	_, ok := args[0].(a.LocalSymbol)
	return a.Bool(ok)
}

func init() {
	var sym *symFunction
	var genSym *genSymFunction
	var isSymbol *isSymbolFunction
	var isLocal *isLocalFunction

	RegisterBuiltIn(symName, sym)
	RegisterBuiltIn(genSymName, genSym)
	RegisterBuiltIn(isSymbolName, isSymbol)
	RegisterBuiltIn(isLocalName, isLocal)
}
