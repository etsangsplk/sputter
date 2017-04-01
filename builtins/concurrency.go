package builtins

import a "github.com/kode4food/sputter/api"

var (
	// MetaEmitter is the key used to emit to a Channel
	MetaEmitter = a.NewKeyword("emit")

	// MetaClose is the key used to close a Channel
	MetaClose = a.NewKeyword("close")

	// MetaSequence is the key used to retrieve the Sequence from a Channel
	MetaSequence = a.NewKeyword("seq")
)

func closeFunction(e a.Emitter) a.Function {
	return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, 0)
		e.Close()
		return a.Nil
	}).(a.Function)
}

func emitFunction(e a.Emitter) a.Function {
	return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			e.Emit(a.Eval(c, i.First()))
		}
		return a.Nil
	}).WithMetadata(a.Metadata{
		MetaEmitter: a.True,
	}).(a.Function)
}

func channel(_ a.Context, args a.Sequence) a.Value {
	buf := 0
	if a.AssertArityRange(args, 0, 1) == 1 {
		v := a.AssertNumber(args.First())
		f, _ := v.Float64()
		buf = int(f)
	}
	e, s := a.NewChannel(buf)

	return a.ArrayMap{
		a.Vector{MetaEmitter, emitFunction(e)},
		a.Vector{MetaClose, closeFunction(e)},
		a.Vector{MetaSequence, s},
	}
}

func _go(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	e, s := a.NewChannel(0)

	l := a.ChildContext(c)
	l.Put("emit", emitFunction(e))

	go func() {
		a.EvalSequence(l, args)
		e.Close()
	}()

	return s
}

func init() {
	registerAnnotated(
		a.NewFunction(channel).WithMetadata(a.Metadata{
			a.MetaName: a.Name("channel"),
		}),
	)

	registerAnnotated(
		a.NewFunction(_go).WithMetadata(a.Metadata{
			a.MetaName: a.Name("go"),
		}),
	)
}
