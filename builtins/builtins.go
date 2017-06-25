package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	e "github.com/kode4food/sputter/evaluator"
)

// Namespace is a special Namespace for built-in identifiers
var Namespace = a.GetNamespace(a.BuiltInDomain)

func registerAnnotated(v a.Annotated) {
	n := v.Metadata()[a.MetaName].(a.Name)
	Namespace.Put(n, v.(a.Value))
}

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalBlock(c, args)
}

func read(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	s := a.AssertStr(v)
	return e.ReadStr(c, s)
}

func eval(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	return a.Eval(c, v)
}

func init() {
	registerAnnotated(
		a.NewFunction(do).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("do"),
			a.MetaDoc:     d.Get("do"),
			a.MetaSpecial: a.True,
		}),
	)

	registerAnnotated(
		a.NewFunction(read).WithMetadata(a.Metadata{
			a.MetaName: a.Name("read"),
			a.MetaDoc:  d.Get("read"),
		}),
	)

	registerAnnotated(
		a.NewFunction(eval).WithMetadata(a.Metadata{
			a.MetaName: a.Name("eval"),
			a.MetaDoc:  d.Get("eval"),
		}),
	)
}
