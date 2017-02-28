package builtins

import a "github.com/kode4food/sputter/api"

// Context is a special Context of built-in identifiers
var Context = a.GetNamespace(a.BuiltInDomain)

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalSequence(c, args)
}

func quote(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	return v
}

func init() {
	a.PutFunction(Context, &a.Function{Name: "do", Exec: do})

	a.PutFunction(Context, &a.Function{
		Name: "quote",
		Exec: quote,
		Data: true,
	})
}
