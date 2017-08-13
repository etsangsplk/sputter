package builtins

import a "github.com/kode4food/sputter/api"

// ExpectedBindings is raised if a binding vector isn't an even number
const ExpectedBindings = "expected bindings in the form: name value"

type (
	defFunction struct{ BaseBuiltIn }
	letFunction struct{ BaseBuiltIn }
)

func (*defFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	ns := a.GetContextNamespace(c)

	s := args.First()
	n := a.AssertUnqualified(s).Name()
	v := args.Rest().First()

	ns.Put(n, a.Eval(c, v))
	return s
}

func (*letFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := a.ChildContext(c)

	b := a.AssertVector(args.First())
	bc := b.Count()
	if bc%2 != 0 {
		panic(a.ErrStr(ExpectedBindings))
	}

	for i := 0; i < bc; i++ {
		s, _ := b.ElementAt(i)
		n := a.AssertUnqualified(s).Name()
		i++
		v, _ := b.ElementAt(i)
		l.Put(n, a.Eval(l, v))
	}

	return a.MakeBlock(args.Rest()).Eval(l)
}

func init() {
	var def *defFunction
	var let *letFunction

	RegisterBuiltIn("def", def)
	RegisterBuiltIn("let", let)
}
