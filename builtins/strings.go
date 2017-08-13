package builtins

import (
	"bytes"
	a "github.com/kode4food/sputter/api"
)

type (
	strFunction       struct{ BaseBuiltIn }
	readerStrFunction struct{ BaseBuiltIn }
)

func (f *strFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToStr(args)
}

func (f *readerStrFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	var buf bytes.Buffer
	if args.IsSequence() {
		buf.WriteString(string(args.First().Str()))
	}
	var t a.Value
	for i := args.Rest(); i.IsSequence(); {
		t, i = i.Split()
		buf.WriteString(" ")
		buf.WriteString(string(t.Str()))
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

	RegisterBuiltIn("str", str)
	RegisterBuiltIn("str!", readerStr)

	RegisterSequencePredicate("str?", isStr)
}
