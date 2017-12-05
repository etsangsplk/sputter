package builtins

import a "github.com/kode4food/sputter/api"

const (
	isIdenticalName = "eq"
	isNilName       = "nil?"
	isKeywordName   = "keyword?"
	isSymbolName    = "symbol?"
	isLocalName     = "local?"
)

type (
	isIdenticalFunction struct{ a.BaseFunction }
	isNilFunction       struct{ a.BaseFunction }
	isKeywordFunction   struct{ a.BaseFunction }
	isSymbolFunction    struct{ a.BaseFunction }
	isLocalFunction     struct{ a.BaseFunction }
)

// RegisterPredicate registers a simple predicate
func RegisterPredicate(n a.Name, fn a.Applicable) {
	pos := a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		if !a.Truthy(fn.Apply(c, args)) {
			return a.False
		}
		return a.True
	})

	neg := a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		if a.Truthy(pos.Apply(c, args)) {
			return a.False
		}
		return a.True
	})

	RegisterFunction(n, pos)
	RegisterFunction(a.Name("!"+n), neg)
}

// RegisterSequencePredicate registers a set-based predicate
func RegisterSequencePredicate(n a.Name, fn a.Applicable) {
	pos := a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
			if !a.Truthy(fn.Apply(c, a.Values{f})) {
				return a.False
			}
		}
		return a.True
	})

	neg := a.NewExecFunction(func(c a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
			if a.Truthy(fn.Apply(c, a.Values{f})) {
				return a.False
			}
		}
		return a.True
	})

	RegisterFunction(n, pos)
	RegisterFunction(a.Name("!"+n), neg)
}

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

	RegisterPredicate(isIdenticalName, isIdentical)
	RegisterSequencePredicate(isNilName, isNil)
	RegisterSequencePredicate(isKeywordName, isKeyword)
	RegisterSequencePredicate(isSymbolName, isSymbol)
	RegisterSequencePredicate(isLocalName, isLocal)
}
