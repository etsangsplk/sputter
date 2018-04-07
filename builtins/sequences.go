package builtins

import a "github.com/kode4food/sputter/api"

const (
	seqName       = "seq"
	firstName     = "first"
	restName      = "rest"
	lastName      = "last"
	consName      = "cons"
	conjName      = "conj"
	lenName       = "len"
	nthName       = "nth"
	getName       = "get"
	isSeqName     = "is-seq"
	isLenName     = "is-len"
	isIndexedName = "is-indexed"
)

type (
	seqFunction       struct{ BaseBuiltIn }
	firstFunction     struct{ BaseBuiltIn }
	restFunction      struct{ BaseBuiltIn }
	lastFunction      struct{ BaseBuiltIn }
	consFunction      struct{ BaseBuiltIn }
	conjFunction      struct{ BaseBuiltIn }
	lenFunction       struct{ BaseBuiltIn }
	nthFunction       struct{ BaseBuiltIn }
	getFunction       struct{ BaseBuiltIn }
	isSeqFunction     struct{ BaseBuiltIn }
	isLenFunction     struct{ BaseBuiltIn }
	isIndexedFunction struct{ BaseBuiltIn }

	forProc func(a.Context)
)

func (*seqFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 1)
	if s, ok := args[0].(a.Sequence); ok {
		return s
	}
	return a.Nil
}

func fetchSequence(args a.Vector) a.Sequence {
	a.AssertArity(args, 1)
	return args[0].(a.Sequence)
}

func (*firstFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return fetchSequence(args).First()
}

func (*restFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return fetchSequence(args).Rest()
}

func (*lastFunction) Apply(_ a.Context, args a.Vector) a.Value {
	var l a.Value = a.Nil
	for f, r, ok := fetchSequence(args).Split(); ok; f, r, ok = r.Split() {
		l = f
	}
	return l
}

func (*consFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 2)
	h := args[0]
	r := args[1]
	return r.(a.Sequence).Prepend(h)
}

func (*conjFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 2)
	s := args[0].(a.Conjoiner)
	for _, f := range args[1:] {
		s = s.Conjoin(f).(a.Conjoiner)
	}
	return s
}

func (*lenFunction) Apply(_ a.Context, args a.Vector) a.Value {
	s := fetchSequence(args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func (*nthFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 1)
	s := args[0].(a.Indexed)
	return a.IndexedApply(s, args[1:])
}

func (*getFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 1)
	s := args[0].(a.Mapped)
	return a.MappedApply(s, args[1:])
}

func (*isSeqFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if s, ok := args[0].(a.Sequence); ok && s.IsSequence() {
		return a.True
	}
	return a.False
}

func (*isLenFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(a.CountedSequence); ok {
		return a.True
	}
	return a.False
}

func (*isIndexedFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(a.IndexedSequence); ok {
		return a.True
	}
	return a.False
}

func init() {
	var seq *seqFunction
	var first *firstFunction
	var rest *restFunction
	var last *lastFunction
	var cons *consFunction
	var conj *conjFunction
	var _len *lenFunction
	var nth *nthFunction
	var get *getFunction
	var isSeq *isSeqFunction
	var isLen *isLenFunction
	var isIndexed *isIndexedFunction

	RegisterBuiltIn(seqName, seq)
	RegisterBuiltIn(firstName, first)
	RegisterBuiltIn(restName, rest)
	RegisterBuiltIn(lastName, last)
	RegisterBuiltIn(consName, cons)
	RegisterBuiltIn(conjName, conj)
	RegisterBuiltIn(lenName, _len)
	RegisterBuiltIn(nthName, nth)
	RegisterBuiltIn(getName, get)
	RegisterBuiltIn(isSeqName, isSeq)
	RegisterBuiltIn(isLenName, isLen)
	RegisterBuiltIn(isIndexedName, isIndexed)
}
