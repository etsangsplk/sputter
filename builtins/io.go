package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

type outputFunc func(a.Value)

func raw(v a.Value) {
	fmt.Print(a.String(v))
}

func pretty(v a.Value) {
	fmt.Print(v)
}

func out(c a.Context, args a.Sequence, o outputFunc) a.Value {
	i := a.Iterate(args)
	for v, ok := i.Next(); ok; {
		o(a.Eval(c, v))
		v, ok = i.Next()
		if ok {
			fmt.Print(" ")
		}
	}
	return a.Nil
}

func outn(c a.Context, args a.Sequence, o outputFunc) a.Value {
	r := out(c, args, o)
	fmt.Println("")
	return r
}

func pr(c a.Context, args a.Sequence) a.Value {
	return out(c, args, raw)
}

func prn(c a.Context, args a.Sequence) a.Value {
	return outn(c, args, raw)
}

func _print(c a.Context, args a.Sequence) a.Value {
	return out(c, args, pretty)
}

func _println(c a.Context, args a.Sequence) a.Value {
	return outn(c, args, pretty)
}

func init() {
	registerFunction(&a.Function{Name: "pr", Exec: pr})
	registerFunction(&a.Function{Name: "prn", Exec: prn})
	registerFunction(&a.Function{Name: "print", Exec: _print})
	registerFunction(&a.Function{Name: "println", Exec: _println})
}
