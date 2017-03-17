package builtins

import (
	a "github.com/kode4food/sputter/api"
)

func isSequence(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	if s, ok := a.Eval(c, v).(a.Sequence); ok {
		if s.IsSequence() {
			return a.True
		}
	}
	return a.False
}

func fetchSequence(c a.Context, args a.Sequence) a.Sequence {
	a.AssertArity(args, 1)
	return a.AssertSequence(a.Eval(c, args.First()))
}

func first(c a.Context, args a.Sequence) a.Value {
	return fetchSequence(c, args).First()
}

func rest(c a.Context, args a.Sequence) a.Value {
	return fetchSequence(c, args).Rest()
}

func cons(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := a.Eval(c, args.First())
	r := a.Eval(c, args.Rest().First())
	return a.AssertSequence(r).Prepend(f)
}

func _len(c a.Context, args a.Sequence) a.Value {
	s := fetchSequence(c, args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func init() {
	registerPredicate(&a.Function{Name: "seq?", Exec: isSequence})
	registerFunction(&a.Function{Name: "first", Exec: first})
	registerFunction(&a.Function{Name: "rest", Exec: rest})
	registerFunction(&a.Function{Name: "cons", Exec: cons})
	registerFunction(&a.Function{Name: "len", Exec: _len})
}
