package builtins

import a "github.com/kode4food/sputter/api"

const (
	isIdenticalName = "eq"
	isNilName       = "nil?"
	isKeywordName   = "keyword?"
	isSymbolName    = "symbol?"
	isLocalName     = "local?"
)

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
func RegisterSequencePredicate(n a.Name, fn a.ValueFilter) {
	pos := NewPredicate(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
			if !fn(f) {
				return a.False
			}
		}
		return a.True
	})

	neg := NewPredicate(func(_ a.Context, args a.Sequence) a.Value {
		a.AssertMinimumArity(args, 1)
		for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
			if fn(f) {
				return a.False
			}
		}
		return a.True
	})

	RegisterFunction(n, pos)
	RegisterFunction(a.Name("!"+n), neg)
}

func isIdentical(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	l := args.First()
	for f, r, ok := args.Split(); ok; f, r, ok = r.Split() {
		if l != f {
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

func isSymbol(v a.Value) bool {
	if _, ok := v.(a.Symbol); ok {
		return true
	}
	return false
}

func isLocal(v a.Value) bool {
	if _, ok := v.(a.LocalSymbol); ok {
		return true
	}
	return false
}

func init() {
	RegisterPredicate(isIdenticalName, isIdentical)

	RegisterSequencePredicate(isNilName, isNil)
	RegisterSequencePredicate(isKeywordName, isKeyword)
	RegisterSequencePredicate(isSymbolName, isSymbol)
	RegisterSequencePredicate(isLocalName, isLocal)
}
