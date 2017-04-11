// Package builtins defines the functions and macros that serve as the
// standard runtime for the Sputter interpreter.
package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

// Namespace is a special Namespace for built-in identifiers
var Namespace = a.GetNamespace(a.BuiltInDomain)

func registerAnnotated(v a.Annotated) {
	n := v.Metadata()[a.MetaName].(a.Name)
	Namespace.Put(n, v)
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
			a.MetaDoc:  d.Get("do"),
		}),
	)

	registerAnnotated(
		a.NewMacro(quote).WithMetadata(a.Metadata{
			a.MetaName: a.Name("quote"),
			a.MetaDoc:  d.Get("quote"),
		}),
	)
}
