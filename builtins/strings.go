package builtins

import (
	"bytes"
	a "github.com/kode4food/sputter/api"
)

const (
	strName       = "str"
	readerStrName = "str!"
	isStrName     = "str?"
)

type (
	strFunction       struct{ BaseBuiltIn }
	readerStrFunction struct{ BaseBuiltIn }
)

var emptyString = a.Str("")

func (*strFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToStr(args)
}

func (*readerStrFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	f, r, ok := args.Split()
	if !ok {
		return emptyString
	}

	var b bytes.Buffer
	b.WriteString(string(f.Str()))
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		b.WriteString(" ")
		b.WriteString(string(f.Str()))
	}
	return a.Str(b.String())
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

	RegisterBuiltIn(strName, str)
	RegisterBuiltIn(readerStrName, readerStr)

	RegisterSequencePredicate(isStrName, isStr)
}
