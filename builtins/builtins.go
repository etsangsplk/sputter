package builtins

import a "github.com/kode4food/sputter/api"

// BuiltIns is a special Namespace for built-in identifiers
var BuiltIns = a.GetNamespace(a.BuiltInDomain)

// PutFunction puts a Function into a Namespace by its Name
func putFunction(ns a.Namespace, f *a.Function) {
	ns.Put(f.Name, f)
}

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
	putFunction(BuiltIns, &a.Function{Name: "do", Apply: do})

	putFunction(BuiltIns, &a.Function{
		Name:  "quote",
		Apply: quote,
		Data:  true,
	})
}
