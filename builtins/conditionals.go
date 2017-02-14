package builtins

import a "github.com/kode4food/sputter/api"

func _if(c *a.Context, args a.Iterable) a.Value {
	AssertArityRange(args, 2, 3)
	i := args.Iterate()
	condVal, _ := i.Next()
	cond := a.Evaluate(c, condVal)
	if !a.Truthy(cond) {
		i.Next()
	}
	result, _ := i.Next()
	return a.Evaluate(c, result)
}

func init() {
	Context.PutFunction(&a.Function{Name: "if", Exec: _if})
}
