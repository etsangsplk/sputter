package builtins

import a "github.com/kode4food/sputter/api"

const (
	vectorName   = "vector"
	isVectorName = "is-vector"
)

type (
	vectorFunction   struct{ BaseBuiltIn }
	isVectorFunction struct{ BaseBuiltIn }
)

func (*vectorFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToVector(args)
}

func (*isVectorFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Vector); ok {
		return a.True
	}
	return a.False
}

func init() {
	var vector *vectorFunction
	var isVector *isVectorFunction

	RegisterBuiltIn(vectorName, vector)
	RegisterBuiltIn(isVectorName, isVector)
}
