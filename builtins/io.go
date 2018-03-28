package builtins

import (
	"io"
	"os"

	a "github.com/kode4food/sputter/api"
)

const (
	newlineName = "*newline*"
	spaceName   = "*space*"
	stdinName   = "*in*"
	stdoutName  = "*out*"
	stderrName  = "*err*"

	newLine = a.Str("\n")
	space   = a.Str(" ")
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
		a.TypeKey: a.Name("writer"),
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
	return a.NewExecFunction(func(_ a.Context, args a.Values) a.Value {
		a.AssertMinimumArity(args, 1)
		for _, f := range args {
			w.Write(f)
		}
		return a.Nil
	})
}

func bindCloser(c a.Closer) a.Function {
	return a.NewExecFunction(func(_ a.Context, args a.Values) a.Value {
		a.AssertArity(args, 0)
		c.Close()
		return a.Nil
	})
}

func init() {
	Namespace.Put(newlineName, newLine)
	Namespace.Put(spaceName, space)
	Namespace.Put(stdinName, makeReader(os.Stdin, a.LineInput))
	Namespace.Put(stdoutName, makeWriter(os.Stdout, a.StrOutput))
	Namespace.Put(stderrName, makeWriter(os.Stderr, a.StrOutput))
}
