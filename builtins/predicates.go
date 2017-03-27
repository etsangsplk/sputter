package builtins

import a "github.com/kode4food/sputter/api"

func identical(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	l := args.First()
	r := args.Rest().First()
	if a.Eval(c, l) == a.Eval(c, r) {
		return a.True
	}
	return a.False
}

func isNil(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	for i := args; i.IsSequence(); i = i.Rest() {
		if a.Eval(c, i.First()) != a.Nil {
			return a.False
		}
	}
	return a.True
}

func init() {
	registerPredicate(
		a.NewFunction(identical).WithMetadata(a.Variables{
			a.MetaName: a.Name("eq"),
		}).(a.Function),
	)

	registerPredicate(
		a.NewFunction(isNil).WithMetadata(a.Variables{
			a.MetaName: a.Name("nil?"),
		}).(a.Function),
	)
}
