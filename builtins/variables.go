package builtins

import a "github.com/kode4food/sputter/api"

const (
	// ExpectedBindings is raised if a binding vector isn't an even number
	ExpectedBindings = "expected bindings in the form: name value"

	letName = "let*"
)

type letFunction struct{ BaseBuiltIn }

func (l *letFunction) Apply(c a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 2)

	b := args[0].(a.Vector)
	bc := b.Count()
	if bc%2 != 0 {
		panic(a.ErrStr(ExpectedBindings))
	}

	v := make(a.Variables, bc/2)
	lc := a.ChildContext(c, v)
	for i := 0; i < bc; i++ {
		s := b[i]
		n := s.(a.LocalSymbol).Name()
		i++
		val := b[i]
		v[n] = a.Eval(lc, val)
	}
	return a.EvalVectorAsBlock(lc, args[1:])
}

func init() {
	var let *letFunction

	RegisterBuiltIn(letName, let)
}
