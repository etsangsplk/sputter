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

func (*isIdenticalFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 2)
	l := args[0]
	for _, f := range args {
		if l != f {
			return a.False
		}
	}
	return a.True
}

func (*isAtomFunction) Apply(_ a.Context, args a.Vector) a.Value {
	_, ok := args[0].(a.Evaluable)
	return a.Bool(!ok)
}

func (*isNilFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return a.Bool(args[0] == a.Nil)
}

func (*isKeywordFunction) Apply(_ a.Context, args a.Vector) a.Value {
	_, ok := args[0].(a.Keyword)
	return a.Bool(ok)
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
