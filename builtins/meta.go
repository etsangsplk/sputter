package builtins

import a "github.com/kode4food/sputter/api"

const (
	withMetaName = "with-meta"
	metaName     = "meta"
	isMetaName   = "is-meta"
)

type (
	withMetaFunction struct{ BaseBuiltIn }
	getMetaFunction  struct{ BaseBuiltIn }
	isMetaFunction   struct{ BaseBuiltIn }
)

func toProperties(args a.MappedSequence) a.Properties {
	res := make(a.Properties)
	for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
		p := f.(a.Sequence)
		k := p.First()
		v := p.Rest().First()
		res[k] = v
	}
	return res
}

func fromMetadata(m a.Object) a.Value {
	r := make([]a.Vector, 0)
	for k, v := range m.Flatten() {
		r = append(r, a.Values{k, v})
	}
	return a.NewAssociative(r...)
}

func (*withMetaFunction) Apply(_ a.Context, args a.Values) a.Value {
	a.AssertMinimumArity(args, 2)
	o := args[0].(a.AnnotatedValue)
	for _, f := range args[1:] {
		m := f.(a.MappedSequence)
		o = o.WithMetadata(toProperties(m))
	}
	return o
}

func (*getMetaFunction) Apply(_ a.Context, args a.Values) a.Value {
	a.AssertArity(args, 1)
	o := args[0].(a.Annotated)
	return fromMetadata(o.Metadata())
}

func (*isMetaFunction) Apply(_ a.Context, args a.Values) a.Value {
	if _, ok := args[0].(a.Annotated); ok {
		return a.True
	}
	return a.False
}

func init() {
	var withMeta *withMetaFunction
	var getMeta *getMetaFunction
	var isMeta *isMetaFunction

	RegisterBuiltIn(withMetaName, withMeta)
	RegisterBuiltIn(metaName, getMeta)
	RegisterBuiltIn(isMetaName, isMeta)
}
