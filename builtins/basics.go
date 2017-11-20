package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

const (
	errorName   = "make-error"
	raiseName   = "raise"
	recoverName = "recover"
	doName      = "do"
	readName    = "read"
	evalName    = "eval"
)

type (
	errorFunction   struct{ BaseBuiltIn }
	raiseFunction   struct{ BaseBuiltIn }
	recoverFunction struct{ BaseBuiltIn }
	doFunction      struct{ BaseBuiltIn }
	readFunction    struct{ BaseBuiltIn }
	evalFunction    struct{ BaseBuiltIn }
)

func (*errorFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	return a.Err(toProperties(args.First().(a.Associative)))
}

func (*raiseFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	panic(args.First())
}

func (*recoverFunction) Apply(c a.Context, args a.Sequence) (res a.Value) {
	a.AssertArity(args, 2)

	defer func() {
		if rec := recover(); rec != nil {
			post := a.Eval(c, args.Rest().First()).(a.Applicable)
			res = post.Apply(c, a.Values{rec.(a.Value)})
		}
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
	s := v.(a.Sequence)
	return e.ReadStr(c, a.SequenceToStr(s))
}

func (*evalFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	return a.Eval(c, v)
}

func init() {
	var _error *errorFunction
	var raise *raiseFunction
	var _recover *recoverFunction
	var do *doFunction
	var read *readFunction
	var eval *evalFunction

	RegisterBuiltIn(errorName, _error)
	RegisterBuiltIn(raiseName, raise)
	RegisterBuiltIn(recoverName, _recover)
	RegisterBuiltIn(doName, do)
	RegisterBuiltIn(readName, read)
	RegisterBuiltIn(evalName, eval)
}
