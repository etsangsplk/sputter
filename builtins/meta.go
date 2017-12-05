package builtins

import a "github.com/kode4food/sputter/api"

const (
	withMetaName = "with-meta"
	metaName     = "meta"
	isMetaName   = "meta?"
)

type (
	withMetaFunction struct{ BaseBuiltIn }
	getMetaFunction  struct{ BaseBuiltIn }

	isAnnotatedFunction struct{ a.BaseFunction }
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

func (*withMetaFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := args.First().(a.Annotated)
	m := args.Rest().First().(a.MappedSequence)
	return o.WithMetadata(toProperties(m))
}

func (*getMetaFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	o := args.First().(a.Annotated)
	return fromMetadata(o.Metadata())
}

func (*isAnnotatedFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Annotated); ok {
		return a.True
	}
	return a.False
}

func init() {
	var withMeta *withMetaFunction
	var getMeta *getMetaFunction
	var isAnnotated *isAnnotatedFunction

	RegisterBuiltIn(withMetaName, withMeta)
	RegisterBuiltIn(metaName, getMeta)
	RegisterSequencePredicate(isMetaName, isAnnotated)
}
