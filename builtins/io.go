package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

func _print(c a.Context, args a.Sequence) a.Value {
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r := a.Eval(c, v)
		fmt.Print(r)
	}
	return a.Nil
}

func _println(c a.Context, args a.Sequence) a.Value {
	r := _print(c, args)
	fmt.Println("")
	return r
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "print", Apply: _print})
	a.PutFunction(Context, &a.Function{Name: "println", Apply: _println})
}
