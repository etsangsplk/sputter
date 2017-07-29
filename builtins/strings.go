package builtins

import (
	"bytes"
	a "github.com/kode4food/sputter/api"
)

func str(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToStr(args)
}

func readerString(_ a.Context, args a.Sequence) a.Value {
	var buf bytes.Buffer
	if args.IsSequence() {
		buf.WriteString(string(args.First().Str()))
	}
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		buf.WriteString(" ")
		buf.WriteString(string(i.First().Str()))
	}
	return a.Str(buf.String())
}

func isStr(v a.Value) bool {
	if _, ok := v.(a.Str); ok {
		return true
	}
	return false
}

func init() {
	RegisterBuiltIn("str", str)
	RegisterBuiltIn("str!", readerString)
	RegisterSequencePredicate("str?", isStr)
}
