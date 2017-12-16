package builtins

import a "github.com/kode4food/sputter/api"

const (
	isIdenticalName = "is-eq"
	isAtomName      = "is-atom"
	isNilName       = "is-nil"
	isKeywordName   = "is-keyword"
)

type (
	isIdenticalFunction struct{ BaseBuiltIn }
	isAtomFunction      struct{ BaseBuiltIn }
	isNilFunction       struct{ BaseBuiltIn }
	isKeywordFunction   struct{ BaseBuiltIn }
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

func (*isAtomFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Evaluable); !ok {
		return a.True
	}
	return a.False
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

func init() {
	var isIdentical *isIdenticalFunction
	var isAtom *isAtomFunction
	var isNil *isNilFunction
	var isKeyword *isKeywordFunction

	RegisterBuiltIn(isIdenticalName, isIdentical)
	RegisterBuiltIn(isAtomName, isAtom)
	RegisterBuiltIn(isNilName, isNil)
	RegisterBuiltIn(isKeywordName, isKeyword)
}
