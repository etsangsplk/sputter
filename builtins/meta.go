package builtins

import a "github.com/kode4food/sputter/api"

func toProperties(args a.MappedSequence) a.Properties {
	r := make(a.Properties)
	for i := args.(a.Sequence); i.IsSequence(); i = i.Rest() {
		p := i.First().(a.Sequence)
		k := p.First()
		v := p.Rest().First()
		r[k] = v
	}
	return r
}

func fromMetadata(m a.Object) a.Value {
	r := []a.Vector{}
	for k, v := range m.Flatten() {
		r = append(r, a.NewVector(k, v))
	}
	return a.NewAssociative(r...)
}

func withMeta(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	o := a.AssertAnnotated(args.First())
	if m, ok := args.Rest().First().(a.MappedSequence); ok {
		return o.WithMetadata(toProperties(m)).(a.Value)
	}
	panic(a.ErrStr(a.ExpectedSequence, o))
}

func meta(_ a.Context, args a.Sequence) a.Value {
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
	RegisterBuiltIn("with-meta", withMeta)
	RegisterBuiltIn("meta", meta)
	RegisterSequencePredicate("meta?", isAnnotated)
}
