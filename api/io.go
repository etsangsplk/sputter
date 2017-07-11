package api

import (
	"bufio"
	"io"
)

const (
	// ExpectedWriter is thrown when a Value is not a Writer
	ExpectedWriter = "value is not a writer: %s"

	// ExpectedReader is thrown when a Value is not a Reader
	ExpectedReader = "value is not a reader: %s"
)

type (
	// Reader is used to retrieve Values from a File
	Reader interface {
		Sequence
	}

	// Writer is used to emit Values to a File
	Writer interface {
		Value
		Write(Value)
	}

	// Closer is used to close a File
	Closer interface {
		Close()
	}

	// OutputFunc is a callback used to marshal Values to a Writer
	OutputFunc func(*bufio.Writer, Value)

	// InputFunc is a callback used to unmarshal Values from a Reader
	InputFunc func(*bufio.Reader) (Value, bool)

	wrappedWriter struct {
		writer *bufio.Writer
		output OutputFunc
	}

	wrappedClosingWriter struct {
		*wrappedWriter
		closer io.Closer
	}
)

// NewReader wraps a Go Reader, coupling it with an input function
func NewReader(r io.Reader, i InputFunc) Reader {
	var resolver Resolver
	br := bufio.NewReader(r)

	resolver = func() (Value, bool, Resolver) {
		if v, ok := i(br); ok {
			return v, ok, resolver
		}
		return nil, false, nil
	}

	return NewLazySequence(resolver)
}

// NewWriter wraps a Go Writer, coupling it with an output function
func NewWriter(w io.Writer, o OutputFunc) Writer {
	wrapped := &wrappedWriter{
		writer: bufio.NewWriter(w),
		output: o,
	}
	if c, ok := w.(io.Closer); ok {
		return &wrappedClosingWriter{
			wrappedWriter: wrapped,
			closer:        c,
		}
	}
	return wrapped
}

func (w *wrappedWriter) Write(v Value) {
	w.output(w.writer, v)
	w.writer.Flush()
}

func (w *wrappedClosingWriter) Close() {
	w.writer.Flush()
	w.closer.Close()
}

func (w *wrappedWriter) Str() Str {
	return MakeDumpStr(w)
}

func (w *wrappedWriter) Type() Name {
	return "writer"
}

func strToBytes(s Str) []byte {
	return []byte(string(s))
}

// StrOutput is the standard string-based output function
func StrOutput(w *bufio.Writer, v Value) {
	if s, ok := v.(Str); ok {
		w.Write(strToBytes(s))
		return
	}
	w.Write(strToBytes(v.Str()))
}

// LineInput is the standard single line input function
func LineInput(r *bufio.Reader) (Value, bool) {
	if l, err := r.ReadBytes('\n'); err == nil {
		return Str(l[0 : len(l)-1]), true
	}
	return nil, false
}

// RuneInput is the standard single rune input function
func RuneInput(r *bufio.Reader) (Value, bool) {
	if c, _, err := r.ReadRune(); err == nil {
		return Str(string(c)), true
	}
	return nil, false
}
