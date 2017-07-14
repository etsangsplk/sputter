package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func str(_ a.Context, args a.Sequence) a.Value {
	return a.ToStr(args)
}

func isStr(v a.Value) bool {
	if _, ok := v.(a.Str); ok {
		return true
	}
	return false
}

func init() {
	RegisterBuiltIn("str", str)
	RegisterSequencePredicate("str?", isStr)
}
