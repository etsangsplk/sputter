package builtins

import a "github.com/kode4food/sputter/api"

type (
	withMetaFunction struct{ BaseBuiltIn }
	getMetaFunction  struct{ BaseBuiltIn }
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
	r := []a.Vector{}
	for k, v := range m.Flatten() {
		r = append(r, a.Values{k, v})
	}
	return a.NewAssociative(r...)
}

func (*withMetaFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := a.AssertAnnotated(args.First())
	if m, ok := args.Rest().First().(a.MappedSequence); ok {
		return o.WithMetadata(toProperties(m))
	}
	panic(a.ErrStr(a.ExpectedMapped, o))
}

func (*getMetaFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	o := a.AssertAnnotated(args.First())
	return fromMetadata(o.Metadata())
}

func isAnnotated(v a.Value) bool {
	if _, ok := v.(a.Annotated); ok {
		return true
	}
	return false
}

func init() {
	var withMeta *withMetaFunction
	var getMeta *getMetaFunction

	RegisterBuiltIn("with-meta", withMeta)
	RegisterBuiltIn("meta", getMeta)

	RegisterSequencePredicate("meta?", isAnnotated)
}
