package builtins

import a "github.com/kode4food/sputter/api"

type forProc func(a.Context)

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

func first(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).First()
}

func rest(_ a.Context, args a.Sequence) a.Value {
	return fetchSequence(args).Rest()
}

func cons(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f := args.First()
	r := args.Rest().First()
	return a.AssertSequence(r).Prepend(f)
}

func conj(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	s := a.AssertConjoiner(args.First())
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		v := i.First()
		s = s.Conjoin(v).(a.Conjoiner)
	}
	return s
}

func _len(_ a.Context, args a.Sequence) a.Value {
	s := fetchSequence(args)
	l := a.Count(s)
	return a.NewFloat(float64(l))
}

func nth(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	s := a.AssertIndexed(args.First())
	return a.IndexedApply(s, args.Rest())
}

func get(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	s := a.AssertMapped(args.First())
	return a.MappedApply(s, args.Rest())
}

func init() {
	RegisterSequencePredicate("seq?", isSequence)
	RegisterSequencePredicate("len?", isCounted)
	RegisterSequencePredicate("indexed?", isIndexed)

	RegisterBuiltIn("first", first)
	RegisterBuiltIn("rest", rest)
	RegisterBuiltIn("cons", cons)
	RegisterBuiltIn("conj", conj)
	RegisterBuiltIn("len", _len)
	RegisterBuiltIn("nth", nth)
	RegisterBuiltIn("get", get)
}
