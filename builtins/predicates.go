package builtins

import a "github.com/kode4food/sputter/api"

// PredicateKey identifies a Function as being a predicate
var PredicateKey = a.NewKeyword("predicate")

// NewPredicate creates a new Predicate instance
func NewPredicate(f a.SequenceProcessor) a.Function {
	return a.NewExecFunction(f).WithMetadata(a.Properties{
		PredicateKey: a.True,
	}).(a.Function)
}

// RegisterPredicate registers a simple predicate
func RegisterPredicate(n a.Name, f a.SequenceProcessor) {
	pos := NewPredicate(f)
	neg := NewPredicate(func(c a.Context, args a.Sequence) a.Value {
		if f(c, args) == a.True {
			return a.False
		}
		return a.True
	})

	RegisterFunction(n, pos)
	RegisterFunction(a.Name("!"+n), neg)
}

// RegisterSequencePredicate registers a set-based predicate
func RegisterSequencePredicate(n a.Name, f a.ValueFilter) {
	pos := NewPredicate(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			if !f(i.First()) {
				return a.False
			}
		}
		return a.True
	})

	neg := NewPredicate(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for i := args; i.IsSequence(); i = i.Rest() {
			if f(i.First()) {
				return a.False
			}
		}
		return a.True
	})

	RegisterFunction(n, pos)
	RegisterFunction(a.Name("!"+n), neg)
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
