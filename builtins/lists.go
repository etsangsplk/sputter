package builtins

import (
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	u "github.com/kode4food/sputter/util"
)

func list(c a.Context, args a.Sequence) a.Value {
	s := u.NewStack()
	for i := args; i.IsSequence(); i = i.Rest() {
		s.Push(a.Eval(c, i.First()))
	}

	l := a.Sequence(a.EmptyList)
	for v, ok := s.Pop(); ok; v, ok = s.Pop() {
		l = l.Prepend(v)
	}
	return l
}

func toList(c a.Context, args a.Sequence) a.Value {
	return list(c, concat(c, args).(a.Sequence))
}

func isList(v a.Value) bool {
	if _, ok := v.(*a.List); ok {
		return true
	}
	return false
}

func init() {
	registerAnnotated(
		a.NewFunction(list).WithMetadata(a.Metadata{
			a.MetaName: a.Name("list"),
			a.MetaDoc:  d.Get("list"),
		}),
	)

	registerAnnotated(
		a.NewFunction(toList).WithMetadata(a.Metadata{
			a.MetaName: a.Name("to-list"),
			a.MetaDoc:  d.Get("to-list"),
		}),
	)

	registerSequencePredicate(isList, a.Metadata{
		a.MetaName: a.Name("list?"),
		a.MetaDoc:  d.Get("is-list"),
	})
}
