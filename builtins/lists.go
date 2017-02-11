package builtins

import a "github.com/kode4food/sputter/api"

func isListCommand(c *a.Context, args a.Iterable) a.Value {
	AssertArity(args, 1)
	iter := args.Iterate()
	if val, ok := iter.Next(); ok {
		if _, ok := a.Evaluate(c, val).(*a.List); ok {
			return a.True
		}
	}
	return a.False
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "list?", Exec: isListCommand})
}
