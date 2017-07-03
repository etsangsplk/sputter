package builtins

import (
	"bytes"

	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func str(_ a.Context, args a.Sequence) a.Value {
	var b bytes.Buffer
	for i := args; i.IsSequence(); i = i.Rest() {
		v := i.First()
		if v == a.Nil {
			continue
		} else if s, ok := v.(a.Str); ok {
			b.WriteString(string(s))
		} else {
			b.WriteString(string(v.Str()))
		}
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
	registerAnnotated(
		a.NewFunction(str).WithMetadata(a.Metadata{
			a.MetaName: a.Name("str"),
			a.MetaDoc:  d.Get("str"),
		}),
	)

	registerSequencePredicate(isStr, a.Metadata{
		a.MetaName: a.Name("str?"),
		a.MetaDoc:  d.Get("is-str"),
	})
}
