package builtins

import (
	"bytes"

	a "github.com/kode4food/sputter/api"
)

const (
	strName       = "str"
	readerStrName = "str!"
	isStrName     = "is-str"
)

type (
	strFunction       struct{ BaseBuiltIn }
	readerStrFunction struct{ BaseBuiltIn }
	isStrFunction     struct{ BaseBuiltIn }
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

func (*isStrFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Str); ok {
		return a.True
	}
	return a.False
}

func init() {
	var str *strFunction
	var readerStr *readerStrFunction
	var isStr *isStrFunction

	RegisterBuiltIn(strName, str)
	RegisterBuiltIn(readerStrName, readerStr)
	RegisterBuiltIn(isStrName, isStr)
}
