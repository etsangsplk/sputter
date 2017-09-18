package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

const (
	panicName   = "panic"
	recoverName = "make-recover"
	doName      = "do"
	readName    = "read"
	evalName    = "eval"
)

type (
	panicFunction   struct{ BaseBuiltIn }
	recoverFunction struct{ BaseBuiltIn }
	doFunction      struct{ BaseBuiltIn }
	readFunction    struct{ BaseBuiltIn }
	evalFunction    struct{ BaseBuiltIn }
)

func (*panicFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	p := toProperties(a.SequenceToAssociative(args))
	panic(a.Err(p))
}

func (*recoverFunction) Apply(c a.Context, args a.Sequence) (res a.Value) {
	a.AssertArity(args, 2)

	defer func() {
		post := a.AssertApplicable(a.Eval(c, args.Rest().First()))
		res = post.Apply(c, a.Values{res, recover().(a.Value)})
	}()

	return a.Eval(c, args.First())
}

func (*doFunction) Apply(c a.Context, args a.Sequence) a.Value {
	var res a.Value = a.Nil
	for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
		res = a.Eval(c, f)
	}
	return res
}

func (*readFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	s := a.AssertSequence(v)
	return e.ReadStr(c, a.SequenceToStr(s))
}

func (*evalFunction) Apply(c a.Context, args a.Sequence) a.Value {
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

	RegisterBuiltIn(panicName, _panic)
	RegisterBuiltIn(recoverName, _recover)
	RegisterBuiltIn(doName, do)
	RegisterBuiltIn(readName, read)
	RegisterBuiltIn(evalName, eval)
}
