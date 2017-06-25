package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

// ExpectedBindings is raised if a binding vector isn't an even number
const ExpectedBindings = "expected bindings in the form: name value"

func def(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	ns := a.GetContextNamespace(c)

	s := args.First()
	n := a.AssertUnqualified(s).Name()
	v := args.Rest().First()

	ns.Put(n, a.Eval(c, v))
	return s
}

func let(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := a.ChildContext(c)

	b := a.AssertVector(args.First())
	bc := b.Count()
	if bc%2 != 0 {
		panic(ExpectedBindings)
	}

	for i := 0; i < bc; i++ {
		s, _ := b.ElementAt(i)
		n := a.AssertUnqualified(s).Name()
		i++
		v, _ := b.ElementAt(i)
		l.Put(n, a.Eval(l, v))
	}

	return a.EvalBlock(l, args.Rest())
}

func init() {
	registerAnnotated(
		a.NewFunction(def).WithMetadata(a.Metadata{
			a.MetaName: a.Name("def"),
			a.MetaDoc:  d.Get("def"),
		}),
	)

	registerAnnotated(
		a.NewFunction(let).WithMetadata(a.Metadata{
			a.MetaName: a.Name("let"),
			a.MetaDoc:  d.Get("let"),
		}),
	)
}
