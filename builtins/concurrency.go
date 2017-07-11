package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
)

var (
	// MetaChannel is the key used to identify a Channel
	MetaChannel = a.NewKeyword("channel")

	// MetaEmitter is the key used to emit to a Channel
	MetaEmitter = a.NewKeyword("emit")

	// MetaSequence is the key used to retrieve the Sequence from a Channel
	MetaSequence = a.NewKeyword("seq")
)

var channelPrototype = a.Properties{
	MetaChannel: a.True,
}

func channel(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	e, s := a.NewChannel()

	return channelPrototype.Child(a.Properties{
		MetaEmitter:  bindWriter(e),
		MetaClose:    bindCloser(e),
		MetaSequence: s,
	})
}

func promise(_ a.Context, args a.Sequence) a.Value {
	if a.AssertArityRange(args, 0, 1) == 1 {
		p := a.NewPromise()
		p.Deliver(args.First())
		return p
	}
	return a.NewPromise()
}

func isPromise(v a.Value) bool {
	if _, ok := v.(a.Promise); ok {
		return true
	}
	return false
}

func doAsync(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	go a.EvalBlock(a.ChildContext(c), args)
	return a.Nil
}

func init() {
	registerAnnotated(
		a.NewFunction(channel).WithMetadata(a.Properties{
			a.MetaName: a.Name("channel"),
			a.MetaDoc:  d.Get("channel"),
		}),
	)

	registerAnnotated(
		a.NewFunction(promise).WithMetadata(a.Properties{
			a.MetaName: a.Name("promise"),
			a.MetaDoc:  d.Get("promise"),
		}),
	)

	registerAnnotated(
		a.NewFunction(doAsync).WithMetadata(a.Properties{
			a.MetaName:    a.Name("do-async"),
			a.MetaSpecial: a.True,
		}),
	)

	registerSequencePredicate(isPromise, a.Properties{
		a.MetaName: a.Name("promise?"),
		a.MetaDoc:  d.Get("is-promise"),
	})
}
