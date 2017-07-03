package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

var (
	// MetaEmitter is the key used to emit to a Channel
	MetaEmitter = a.NewKeyword("emit")

	// MetaClose is the key used to close a Channel
	MetaClose = a.NewKeyword("close")

	// MetaSequence is the key used to retrieve the Sequence from a Channel
	MetaSequence = a.NewKeyword("seq")

	// MetaPromise is the key used to identify a Promise
	MetaPromise = a.NewKeyword("promise")
)

var (
	emitterMetadata = a.Metadata{MetaEmitter: a.True}
	promiseMetadata = a.Metadata{MetaPromise: a.True}
)

func closeFunction(e a.Emitter) a.Function {
	return a.NewFunction(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, 0)
		e.Close()
		return a.Nil
	})
}

func emitFunction(e a.Emitter) a.Function {
	return a.NewFunction(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			e.Emit(i.First())
		}
		return a.Nil
	}).WithMetadata(emitterMetadata).(a.Function)
}

func channel(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	e, s := a.NewChannel()

	return a.NewAssociative(
		a.NewVector(MetaEmitter, emitFunction(e)),
		a.NewVector(MetaClose, closeFunction(e)),
		a.NewVector(MetaSequence, s),
	)
}

func promise(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	p := a.NewPromise()

	return a.NewFunction(
		func(_ a.Context, args a.Sequence) a.Value {
			if a.AssertArityRange(args, 0, 1) == 1 {
				return p.Deliver(args.First())
			}
			return p.Resolve()
		},
	).WithMetadata(promiseMetadata).(a.Function)
}

func doAsync(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	go a.EvalBlock(a.ChildContext(c), args)
	return a.Nil
}

func init() {
	registerAnnotated(
		a.NewFunction(channel).WithMetadata(a.Metadata{
			a.MetaName: a.Name("channel"),
			a.MetaDoc:  d.Get("channel"),
		}),
	)

	registerAnnotated(
		a.NewFunction(promise).WithMetadata(a.Metadata{
			a.MetaName: a.Name("promise"),
			a.MetaDoc:  d.Get("promise"),
		}),
	)

	registerAnnotated(
		a.NewFunction(doAsync).WithMetadata(a.Metadata{
			a.MetaName:    a.Name("do-async"),
			a.MetaSpecial: a.True,
		}),
	)
}
