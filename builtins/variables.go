package builtins

import a "github.com/kode4food/sputter/api"

const (
	// ExpectedBindings is raised if a binding vector isn't an even number
	ExpectedBindings = "expected bindings in the form: name value"

	defName = "def"
	letName = "let"
)

type (
	defFunction struct{ BaseBuiltIn }
	letFunction struct{ BaseBuiltIn }
)

func (*defFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	ns := a.GetContextNamespace(c)

	f, r, _ := args.Split()
	n := f.(a.LocalSymbol).Name()
	v := r.First()

	ns.Put(n, a.Eval(c, v))
	return f
}

func (*letFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := a.ChildContext(c)

	f, r, _ := args.Split()
	b := f.(a.Vector)
	bc := b.Count()
	if bc%2 != 0 {
		panic(a.ErrStr(ExpectedBindings))
	}

	for i := 0; i < bc; i++ {
		s, _ := b.ElementAt(i)
		n := s.(a.LocalSymbol).Name()
		i++
		v, _ := b.ElementAt(i)
		l.Put(n, a.Eval(l, v))
	}

	return a.MakeBlock(r).Eval(l)
}

func init() {
	var def *defFunction
	var let *letFunction

	RegisterBuiltIn(defName, def)
	RegisterBuiltIn(letName, let)
}
