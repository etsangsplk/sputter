package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func str(_ a.Context, args a.Sequence) a.Value {
	return a.ToStr(args)
}

func escapeString(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return args.First().Str()
}

func isStr(v a.Value) bool {
	if _, ok := v.(a.Str); ok {
		return true
	}
	return false
}

func init() {
	RegisterBuiltIn("str", str)
	RegisterBuiltIn("str!", escapeString)
	RegisterSequencePredicate("str?", isStr)
}
