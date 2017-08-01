package builtins

import a "github.com/kode4food/sputter/api"

var (
	// Namespace is a special Namespace for built-in identifiers
	Namespace = a.GetNamespace(a.BuiltInDomain)

	builtInFuncs = map[a.Name]a.Function{}
)

// RegisterFunction registers a built-in Function by Name
func RegisterFunction(n a.Name, f a.Function) {
	builtInFuncs[n] = f
}

// RegisterBaseFunction registers a Base-derived Function by Name
func RegisterBaseFunction(n a.Name, f a.HasReflectedFunction) {
	builtInFuncs[n] = a.NewReflectedFunction(f)
}

// GetBuiltIn returns a registered built-in function
func GetBuiltIn(n a.Name) (a.Function, bool) {
	if f, ok := builtInFuncs[n]; ok {
		return f, true
	}
	return nil, false
}

func defBuiltIn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	if f, ok := GetBuiltIn(n); ok {
		var md a.Object = toProperties(a.SequenceToAssociative(args.Rest()))
		r := f.WithMetadata(md)
		a.GetContextNamespace(c).Put(n, r)
		return r
	}
	panic(a.ErrStr(a.KeyNotFound, n))
}

func isBuiltInDomain(s a.Symbol) bool {
	return s.Domain() == a.BuiltInDomain
}

func isBuiltInCall(n a.Name, v a.Value) (a.List, bool) {
	if l, ok := v.(a.List); ok && l.Count() > 0 {
		if s, ok := l.First().(a.Symbol); ok {
			return l, isBuiltInDomain(s) && s.Name() == n
		}
	}
	return nil, false
}

func init() {
	Namespace.Put("def-builtin",
		a.NewExecFunction(defBuiltIn).WithMetadata(a.Properties{
			a.SpecialKey: a.True,
		}),
	)
}
