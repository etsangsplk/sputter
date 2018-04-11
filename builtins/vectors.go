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

func (*vectorFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return a.SequenceToVector(args)
}

func (*isVectorFunction) Apply(_ a.Context, args a.Vector) a.Value {
	_, ok := args[0].(a.Vector)
	return a.Bool(ok)
}

func init() {
	var vector *vectorFunction
	var isVector *isVectorFunction

	RegisterBuiltIn(vectorName, vector)
	RegisterBuiltIn(isVectorName, isVector)
}
