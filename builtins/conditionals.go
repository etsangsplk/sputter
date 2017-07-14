package builtins

import a "github.com/kode4food/sputter/api"

func _if(c a.Context, args a.Sequence) a.Value {
	i := a.AssertArityRange(args, 2, 3)
	cond := a.Eval(c, args.First())
	rest := args.Rest()
	if a.Truthy(cond) {
		return a.Eval(c, rest.First())
	}
	if i == 3 {
		return a.Eval(c, rest.Rest().First())
	}
	return a.Nil
}

func init() {
	RegisterBuiltIn("if", _if) // special
}
