package builtins

import a "github.com/kode4food/sputter/api"

// BuiltIns is a special Namespace for built-in identifiers
var BuiltIns = a.GetNamespace(a.BuiltInDomain)

func registerAnnotated(v a.Annotated) {
	n := v.Metadata()[a.MetaName].(a.Name)
	BuiltIns.Put(n, v)
}

func registerPredicate(f a.Function) {
	pn := a.Name("!" + f.Name())
	p := a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		r := f.Apply(c, args)
		if r == a.True {
			return a.False
		}
		return a.True
	}).WithMetadata(a.Metadata{
		a.MetaName: pn,
	})

	registerAnnotated(f)
	registerAnnotated(p)
}

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalSequence(c, args)
}

func quote(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return a.Quote(args.First())
}

func init() {
	registerAnnotated(
		a.NewFunction(do).WithMetadata(a.Metadata{
			a.MetaName: a.Name("do"),
		}),
	)

	registerAnnotated(
		a.NewMacro(quote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("quote"),
		}),
	)
}
