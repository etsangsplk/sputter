package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

func print(c *a.Context, args a.Iterable) a.Value {
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r := a.Evaluate(c, v)
		fmt.Print(r)
	}
	return a.Nil
}

func println(c *a.Context, args a.Iterable) a.Value {
	print(c, args)
	fmt.Println("")
	return a.Nil
}

func init() {
	Context.PutFunction(&a.Function{Name: "print", Exec: print})
	Context.PutFunction(&a.Function{Name: "println", Exec: println})
}
