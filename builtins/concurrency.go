package builtins

import a "github.com/kode4food/sputter/api"

var (
	// MetaEmitter is the key used to retrieve the Emitter from a Channel
	MetaEmitter = a.NewKeyword("emit")

	// MetaSequence is the key used to retrieve the Sequence from a Channel
	MetaSequence = a.NewKeyword("seq")
)

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
	a.AssertArity(args, 0)
	e, s := a.NewChannel()

	return a.ArrayMap{
		a.Vector{MetaEmitter, emitFunction(e)},
		a.Vector{MetaSequence, s},
	}
}

func _go(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	e, s := a.NewChannel()

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
