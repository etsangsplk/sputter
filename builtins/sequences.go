package builtins

import a "github.com/kode4food/sputter/api"

type (
	firstFunction struct{ a.ReflectedFunction }
	restFunction  struct{ a.ReflectedFunction }
	consFunction  struct{ a.ReflectedFunction }
	conjFunction  struct{ a.ReflectedFunction }
	lenFunction   struct{ a.ReflectedFunction }
	nthFunction   struct{ a.ReflectedFunction }
	getFunction   struct{ a.ReflectedFunction }

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

func (f *firstFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).First()
}

func (f *restFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).Rest()
}

func (f *consFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	h := args.First()
	r := args.Rest().First()
	return a.AssertSequence(r).Prepend(h)
}

func (f *conjFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	s := a.AssertConjoiner(args.First())
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		v := i.First()
		s = s.Conjoin(v).(a.Conjoiner)
	}
	return s
}

func (f *lenFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	s := fetchSequence(args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func (f *nthFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	s := a.AssertIndexed(args.First())
	return a.IndexedApply(s, args.Rest())
}

func (f *getFunction) Apply(_ a.Context, args a.Sequence) a.Value {
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

	RegisterBaseFunction("first", first)
	RegisterBaseFunction("rest", rest)
	RegisterBaseFunction("cons", cons)
	RegisterBaseFunction("conj", conj)
	RegisterBaseFunction("len", _len)
	RegisterBaseFunction("nth", nth)
	RegisterBaseFunction("get", get)

	RegisterSequencePredicate("seq?", isSequence)
	RegisterSequencePredicate("len?", isCounted)
	RegisterSequencePredicate("indexed?", isIndexed)
}
