package builtins

import a "github.com/kode4food/sputter/api"

func ifCommand(c *a.Context, args a.Iterable) a.Value {
	AssertArityRange(args, 2, 3)
	iter := args.Iterate()
	condValue, _ := iter.Next()
	cond := a.Evaluate(c, condValue)
	if !a.Truthy(cond) {
		iter.Next()
	}
	result, _ := iter.Next()
	return a.Evaluate(c, result)
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "if", Exec: ifCommand})
}
