package builtins

import a "github.com/kode4food/sputter/api"

func cond(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	i := args.Iterate()
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		if b, ok := e.(*a.Cons); ok {
			a.AssertMinimumArity(b, 2)
			bi := b.Iterate()
			cond, _ := bi.Next()
			if a.Truthy(a.Eval(c, cond)) {
				return a.EvalSequence(c, bi.Rest())
			}
		}
	}
	return a.Nil
}

func _if(c a.Context, args a.Sequence) a.Value {
	a.AssertArityRange(args, 2, 3)
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
	a.PutFunction(Context, &a.Function{Name: "cond", Exec: cond})
	a.PutFunction(Context, &a.Function{Name: "if", Exec: _if})
}
