package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

func printCommand(c *a.Context, args a.Iterable) a.Value {
	iter := args.Iterate()
	for value, ok := iter.Next(); ok; value, ok = iter.Next() {
		result := a.Evaluate(c, value)
		fmt.Print(result)
	}
	return a.EmptyList
}

func printlnCommand(c *a.Context, args a.Iterable) a.Value {
	printCommand(c, args)
	fmt.Println("")
	return a.EmptyList
}

func init() {
	BuiltIns.PutFunction(&a.Function{Name: "print", Exec: printCommand})
	BuiltIns.PutFunction(&a.Function{Name: "println", Exec: printlnCommand})	
}
