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

func (*strFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return a.SequenceToStr(args)
}

func (*readerStrFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if len(args) == 0 {
		return emptyString
	}

	var b bytes.Buffer
	b.WriteString(string(args[0].Str()))
	for _, f := range args[1:] {
		b.WriteString(" ")
		b.WriteString(string(f.Str()))
	}
	return a.Str(b.String())
}

func (*isStrFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(a.Str); ok {
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
