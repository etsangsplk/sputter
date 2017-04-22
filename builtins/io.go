package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

type outputFunc func(a.Value)

func raw(v a.Value) {
	fmt.Print(v.Str())
}

func pretty(v a.Value) {
	if s, ok := v.(a.Str); ok {
		fmt.Print(s)
		return
	}
	fmt.Print(v.Str())
}

func out(c a.Context, args a.Sequence, o outputFunc) a.Value {
	if args.IsSequence() {
		o(a.Eval(c, args.First()))
	}
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		fmt.Print(" ")
		o(a.Eval(c, i.First()))
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
	registerAnnotated(
		a.NewFunction(pr).WithMetadata(a.Metadata{
			a.MetaName: a.Name("pr"),
		}),
	)

	registerAnnotated(
		a.NewFunction(prn).WithMetadata(a.Metadata{
			a.MetaName: a.Name("prn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_print).WithMetadata(a.Metadata{
			a.MetaName: a.Name("print"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_println).WithMetadata(a.Metadata{
			a.MetaName: a.Name("println"),
		}),
	)
}
