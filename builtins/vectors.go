package builtins

import a "github.com/kode4food/sputter/api"

type vectorFunction struct{ BaseBuiltIn }

func (*vectorFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToVector(args)
}

func isVector(v a.Value) bool {
	if _, ok := v.(a.Vector); ok {
		return true
	}
	return false
}

func init() {
	var vector *vectorFunction

	RegisterBuiltIn("vector", vector)

	RegisterSequencePredicate("vector?", isVector)
}
