package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

func printCommand(c *a.Context, args a.Iterable) a.Value {
	i := args.Iterate()
	for v, ok := i.Next(); ok; v, ok = i.Next() {
		r := a.Evaluate(c, v)
		fmt.Print(r)
	}
	return a.Nil
}

func printlnCommand(c *a.Context, args a.Iterable) a.Value {
	printCommand(c, args)
	fmt.Println("")
	return a.Nil
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "print", Exec: printCommand})
	BuiltIns.PutFunction(&a.Function{Name: "println", Exec: printlnCommand})
}
