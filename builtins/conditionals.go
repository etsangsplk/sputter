package builtins

import a "github.com/kode4food/sputter/api"

func _if(c a.Context, args a.Sequence) a.Value {
	AssertArityRange(args, 2, 3)
	i := args.Iterate()
	condVal, _ := i.Next()
	cond := a.Eval(c, condVal)
	if !a.Truthy(cond) {
		i.Next()
	}
	result, _ := i.Next()
	return a.Eval(c, result)
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "if", Exec: _if})
}
