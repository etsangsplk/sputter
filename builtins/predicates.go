package builtins

import a "github.com/kode4food/sputter/api"

const (
	isIdenticalName = "is-eq"
	isNilName       = "is-nil"
	isKeywordName   = "is-keyword"
	isSymbolName    = "is-symbol"
	isLocalName     = "is-local"
)

type (
	isIdenticalFunction struct{ BaseBuiltIn }
	isNilFunction       struct{ BaseBuiltIn }
	isKeywordFunction   struct{ BaseBuiltIn }
	isSymbolFunction    struct{ BaseBuiltIn }
	isLocalFunction     struct{ BaseBuiltIn }
)

func (*isIdenticalFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := args.First()
	for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
		if l != f {
			return a.False
		}
	}
	return a.True
}

func (*isNilFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if args.First() == a.Nil {
		return a.True
	}
	return a.False
}

func (*isKeywordFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Keyword); ok {
		return a.True
	}
	return a.False
}

func (*isSymbolFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Symbol); ok {
		return a.True
	}
	return a.False
}

func (*isLocalFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.LocalSymbol); ok {
		return a.True
	}
	return a.False
}

func init() {
	var isIdentical *isIdenticalFunction
	var isNil *isNilFunction
	var isKeyword *isKeywordFunction
	var isSymbol *isSymbolFunction
	var isLocal *isLocalFunction

	RegisterBuiltIn(isIdenticalName, isIdentical)
	RegisterBuiltIn(isNilName, isNil)
	RegisterBuiltIn(isKeywordName, isKeyword)
	RegisterBuiltIn(isSymbolName, isSymbol)
	RegisterBuiltIn(isLocalName, isLocal)
}
