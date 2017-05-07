package builtins

import (
	"fmt"
	"io"
	"os"

	a "github.com/kode4food/sputter/api"
	n "github.com/kode4food/sputter/native"
)

type outputFunc func(io.Writer, a.Value)

const stdinReader = a.Name("*stdin*")
const stdoutWriter = a.Name("*stdout*")

func getStdOut(c a.Context) io.Writer {
	if v, ok := c.Get(stdoutWriter); ok {
		if n, ok := v.(n.Wrapped); ok {
			if w, ok := n.Wrapped().(io.Writer); ok {
				return w
			}
		}
	}
	return os.Stdout
}

func raw(w io.Writer, v a.Value) {
	fmt.Fprint(w, v.Str())
}

func pretty(w io.Writer, v a.Value) {
	if s, ok := v.(a.Str); ok {
		fmt.Fprint(w, s)
		return
	}
	fmt.Fprint(w, v.Str())
}

func out(c a.Context, args a.Sequence, o outputFunc) a.Value {
	w := getStdOut(c)
	if args.IsSequence() {
		o(w, a.Eval(c, args.First()))
	}
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		fmt.Fprint(w, " ")
		o(w, a.Eval(c, i.First()))
	}
	return a.Nil
}

func outn(c a.Context, args a.Sequence, o outputFunc) a.Value {
	r := out(c, args, o)
	w := getStdOut(c)
	fmt.Fprintln(w, "")
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
