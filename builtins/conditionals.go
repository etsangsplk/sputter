package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func cond(c a.Context, args a.Sequence) a.Value {
	for i := args; i.IsSequence(); i = i.Rest() {
		p := a.Eval(c, i.First())
		i = i.Rest()
		if !i.IsSequence() {
			return p
		}
		if a.Truthy(p) {
			return a.Eval(c, i.First())
		}
	}
	return a.Nil
}

// this will be replaced by a macro -> cond
func _if(c a.Context, args a.Sequence) a.Value {
	a.AssertArityRange(args, 2, 3)
	i := a.Iterate(args)
	condVal, _ := i.Next()
	cond := a.Eval(c, condVal)
	if !a.Truthy(cond) {
		i.Next()
	}
	result, _ := i.Next()
	return a.Eval(c, result)
}

func init() {
	registerFunction(&a.Function{Name: "cond", Exec: cond})
	registerFunction(&a.Function{Name: "if", Exec: _if})
}
