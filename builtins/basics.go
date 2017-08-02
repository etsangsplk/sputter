package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

type (
	panicFunction   struct{ BaseBuiltIn }
	recoverFunction struct{ BaseBuiltIn }
	doFunction      struct{ BaseBuiltIn }
	readFunction    struct{ BaseBuiltIn }
	evalFunction    struct{ BaseBuiltIn }
)

func (f *panicFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	p := toProperties(a.SequenceToAssociative(args))
	panic(a.Err(p))
}

func (f *recoverFunction) Apply(c a.Context, args a.Sequence) (res a.Value) {
	a.AssertArity(args, 2)

	defer func() {
		post := a.AssertApplicable(a.Eval(c, args.Rest().First()))
		res = post.Apply(c, a.NewVector(res, recover().(a.Value)))
	}()

	return a.Eval(c, args.First())
}

func (f *doFunction) Apply(c a.Context, args a.Sequence) a.Value {
	var r a.Value = a.Nil
	for i := args; i.IsSequence(); i = i.Rest() {
		r = a.Eval(c, i.First())
	}
	return r
}

func (f *readFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	s := a.AssertSequence(v)
	return e.ReadStr(c, a.SequenceToStr(s))
}

func (f *evalFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	return a.Eval(c, v)
}

func init() {
	var _panic *panicFunction
	var _recover *recoverFunction
	var do *doFunction
	var read *readFunction
	var eval *evalFunction

	RegisterBuiltIn("panic", _panic)
	RegisterBuiltIn("make-recover", _recover)
	RegisterBuiltIn("do", do)
	RegisterBuiltIn("read", read)
	RegisterBuiltIn("eval", eval)
}
