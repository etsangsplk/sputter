package builtins

import a "github.com/kode4food/sputter/api"

// RegisterPredicate registers a simple predicate
func RegisterPredicate(n a.Name, f a.SequenceProcessor) {
	not := a.Name("!" + n)
	RegisterBuiltIn(n, f)
	RegisterBuiltIn(not, func(c a.Context, args a.Sequence) a.Value {
		if f(c, args) == a.True {
			return a.False
		}
		return a.True
	})
}

// RegisterSequencePredicate registers a set-based predicate
func RegisterSequencePredicate(n a.Name, f a.ValueFilter) {
	not := a.Name("!" + n)
	RegisterBuiltIn(n, func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			if !f(i.First()) {
				return a.False
			}
		}
		return a.True
	})

	RegisterBuiltIn(not, func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			if f(i.First()) {
				return a.False
			}
		}
		return a.True
	})
}

func identical(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := args.First()
	for i := args.Rest(); i.IsSequence(); i = i.Rest() {
		if l != i.First() {
			return a.False
		}
	}
	return a.True
}

func isNil(v a.Value) bool {
	return v == a.Nil
}

func isKeyword(v a.Value) bool {
	if _, ok := v.(a.Keyword); ok {
		return true
	}
	return false
}

func init() {
	RegisterPredicate("eq", identical)
	RegisterSequencePredicate("nil?", isNil)
	RegisterSequencePredicate("keyword?", isKeyword)
}
