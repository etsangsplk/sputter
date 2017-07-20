package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

func _panic(_ a.Context, args a.Sequence) a.Value {
	panic(a.ToAssociative(args))
}

func makeRecover(c a.Context, args a.Sequence) (res a.Value) {
	a.AssertArity(args, 2)

	defer func() {
		post := a.AssertApplicable(a.Eval(c, args.Rest().First()))
		res = post.Apply(c, a.NewVector(res, recover().(a.Value)))
	}()

	return a.Eval(c, args.First())
}

func read(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	s := a.AssertSequence(v)
	return e.ReadStr(c, a.ToStr(s))
}

func eval(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	return a.Eval(c, v)
}

func init() {
	RegisterBuiltIn("panic", _panic)
	RegisterBuiltIn("make-recover", makeRecover)
	RegisterBuiltIn("do", a.EvalBlock)
	RegisterBuiltIn("read", read)
	RegisterBuiltIn("eval", eval)
}
