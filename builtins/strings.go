package builtins

import (
	"bytes"
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

func str(c a.Context, args a.Sequence) a.Value {
	var b bytes.Buffer
	for i := args; i.IsSequence(); i = i.Rest() {
		v := a.Eval(c, i.First())
		if s, ok := v.(string); ok {
			b.WriteString(s)
		} else {
			b.WriteString(a.String(v))
		}
	}
	return b.String()
}

func init() {
	registerAnnotated(
		a.NewFunction(str).WithMetadata(a.Metadata{
			a.MetaName: a.Name("str"),
			a.MetaDoc:  d.Get("str"),
		}),
	)
}
