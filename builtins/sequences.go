package builtins

import a "github.com/kode4food/sputter/api"

type (
	firstFunction struct{ BaseBuiltIn }
	restFunction  struct{ BaseBuiltIn }
	consFunction  struct{ BaseBuiltIn }
	conjFunction  struct{ BaseBuiltIn }
	lenFunction   struct{ BaseBuiltIn }
	nthFunction   struct{ BaseBuiltIn }
	getFunction   struct{ BaseBuiltIn }

	forProc func(a.Context)
)

func isSequence(v a.Value) bool {
	if s, ok := v.(a.Sequence); ok {
		return s.IsSequence()
	}
	return false
}

func isCounted(v a.Value) bool {
	if _, ok := v.(a.CountedSequence); ok {
		return true
	}
	return false
}

func isIndexed(v a.Value) bool {
	if _, ok := v.(a.IndexedSequence); ok {
		return true
	}
	return false
}

func fetchSequence(args a.Sequence) a.Sequence {
	a.AssertArity(args, 1)
	return a.AssertSequence(args.First())
}

func (*firstFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).First()
}

func (*restFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).Rest()
}

func (*consFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	h := args.First()
	r := args.Rest().First()
	return a.AssertSequence(r).Prepend(h)
}

func (*conjFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, ok := args.Split()
	s := a.AssertConjoiner(f)
	for f, r, ok = r.Split(); ok; f, r, ok = r.Split() {
		s = s.Conjoin(f).(a.Conjoiner)
	}
	return s
}

func (*lenFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	s := fetchSequence(args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func (*nthFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	s := a.AssertIndexed(args.First())
	return a.IndexedApply(s, args.Rest())
}

func (*getFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	s := a.AssertMapped(args.First())
	return a.MappedApply(s, args.Rest())
}

func init() {
	var first *firstFunction
	var rest *restFunction
	var cons *consFunction
	var conj *conjFunction
	var _len *lenFunction
	var nth *nthFunction
	var get *getFunction

	RegisterBuiltIn("first", first)
	RegisterBuiltIn("rest", rest)
	RegisterBuiltIn("cons", cons)
	RegisterBuiltIn("conj", conj)
	RegisterBuiltIn("len", _len)
	RegisterBuiltIn("nth", nth)
	RegisterBuiltIn("get", get)

	RegisterSequencePredicate("seq?", isSequence)
	RegisterSequencePredicate("len?", isCounted)
	RegisterSequencePredicate("indexed?", isIndexed)
}
