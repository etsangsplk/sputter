package builtins

import (
	"bytes"
	a "github.com/kode4food/sputter/api"
)

type (
	strFunction       struct{ a.ReflectedFunction }
	readerStrFunction struct{ a.ReflectedFunction }
)

func (f *strFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToStr(args)
}

func (f *readerStrFunction) Apply(_ a.Context, args a.Sequence) a.Value {
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
	var str *strFunction
	var readerStr *readerStrFunction

	RegisterBaseFunction("str", str)
	RegisterBaseFunction("str!", readerStr)

	RegisterSequencePredicate("str?", isStr)
}
