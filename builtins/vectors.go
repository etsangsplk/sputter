package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func vector(_ a.Context, args a.Sequence) a.Value {
	return a.ToVector(args)
}

func isVector(v a.Value) bool {
	if _, ok := v.(a.Vector); ok {
		return true
	}
	return false
}

func init() {
	RegisterBuiltIn("vector", vector)
	RegisterSequencePredicate("vector?", isVector)
}
