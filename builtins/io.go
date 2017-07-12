package builtins

import (
	"io"
	"os"

	a "github.com/kode4food/sputter/api"
)

const (
	stdinName  = a.Name("*stdin*")
	stdoutName = a.Name("*stdout*")
	stderrName = a.Name("*stderr*")
	space      = a.Str(" ")
	newLine    = a.Str("\n")
)

type outputFunc func(a.Writer, a.Value)

var (
	// MetaWriter is the key used to wrap a Writer
	MetaWriter = a.NewKeyword("writer")

	// MetaWrite is key used to write to a Writer
	MetaWrite = a.NewKeyword("write")

	// MetaClose is the key used to close a file
	MetaClose = a.NewKeyword("close")

	writerPrototype = a.Properties{
		a.MetaType: a.Name("writer"),
	}
)

func makeReader(r io.Reader, i a.InputFunc) a.Reader {
	return a.NewReader(r, i)
}

func makeWriter(w io.Writer, o a.OutputFunc) a.Object {
	wrapped := a.NewWriter(w, o)

	wrapper := a.Properties{
		MetaWriter: wrapped,
		MetaWrite:  bindWriter(wrapped),
	}

	if c, ok := w.(a.Closer); ok {
		wrapper[MetaClose] = bindCloser(c)
	}

	return writerPrototype.Child(wrapper)
}

func bindWriter(w a.Writer) a.Function {
	return a.NewFunction(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			w.Write(i.First())
		}
		return a.Nil
	})
}

func bindCloser(c a.Closer) a.Function {
	return a.NewFunction(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, 0)
		c.Close()
		return a.Nil
	})
}

func getWriter(c a.Context, n a.Name) a.Writer {
	if v, ok := c.Get(n); ok {
		if o, ok := v.(a.Object); ok {
			if p, ok := o.Get(MetaWriter); ok {
				if w, ok := p.(a.Writer); ok {
					return w
				}
			}
		}
		panic(a.Err(a.ExpectedWriter, v.Str()))
	}
	panic(a.Err(a.UnknownSymbol, n))
}

func raw(e a.Writer, v a.Value) {
	e.Write(v)
}

func pretty(e a.Writer, v a.Value) {
	if s, ok := v.(a.Str); ok {
		e.Write(s.Str())
	} else {
		e.Write(v)
	}
}

func out(c a.Context, args a.Sequence, o outputFunc) a.Writer {
	e := getWriter(c, stdoutName)
	if args.IsSequence() {
		o(e, args.First())
	}
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		e.Write(space)
		o(e, i.First())
	}
	return e
}

func outn(c a.Context, args a.Sequence, o outputFunc) a.Writer {
	e := out(c, args, o)
	e.Write(newLine)
	return e
}

func pr(c a.Context, args a.Sequence) a.Value {
	out(c, args, pretty)
	return a.Nil
}

func prn(c a.Context, args a.Sequence) a.Value {
	outn(c, args, pretty)
	return a.Nil
}

func _print(c a.Context, args a.Sequence) a.Value {
	out(c, args, raw)
	return a.Nil
}

func _println(c a.Context, args a.Sequence) a.Value {
	outn(c, args, raw)
	return a.Nil
}

func init() {
	Namespace.Put(stdinName, makeReader(os.Stdin, a.LineInput))
	Namespace.Put(stdoutName, makeWriter(os.Stdout, a.StrOutput))
	Namespace.Put(stderrName, makeWriter(os.Stderr, a.StrOutput))

	registerAnnotated(
		a.NewFunction(pr).WithMetadata(a.Properties{
			a.MetaName: a.Name("pr"),
		}),
	)

	registerAnnotated(
		a.NewFunction(prn).WithMetadata(a.Properties{
			a.MetaName: a.Name("prn"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_print).WithMetadata(a.Properties{
			a.MetaName: a.Name("print"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_println).WithMetadata(a.Properties{
			a.MetaName: a.Name("println"),
		}),
	)
}
