package builtins

import a "github.com/kode4food/sputter/api"

func _if(c a.Context, args a.Sequence) a.Value {
	i := a.AssertArityRange(args, 2, 3)
	if a.Truthy(a.Eval(c, args.First())) {
		return a.Eval(c, args.Rest().First())
	}
	if i == 3 {
		return a.Eval(c, args.Rest().Rest().First())
	}
	return a.Nil
}

func init() {
	RegisterBuiltIn("if", _if)
}
