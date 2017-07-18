package builtins

import (
	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
)

var (
	// Namespace is a special Namespace for built-in identifiers
	Namespace = a.GetNamespace(a.BuiltInDomain)

	builtInFuncs = map[a.Name]a.SequenceProcessor{}
)

// RegisterBuiltIn adds a built-in SequenceProcessor by Name
func RegisterBuiltIn(n a.Name, proc a.SequenceProcessor) {
	builtInFuncs[n] = proc
}

// GetBuiltIn returns a registered built-in function
func GetBuiltIn(n a.Name) (a.SequenceProcessor, bool) {
	if v, ok := builtInFuncs[n]; ok {
		return v, true
	}
	return nil, false
}

func defBuiltIn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	if f, ok := GetBuiltIn(n); ok {
		var md a.Object = toProperties(a.ToAssociative(args.Rest()))
		md = loadDocumentation(md)

		r := a.NewFunction(f).WithMetadata(md).(a.Value)
		a.GetContextNamespace(c).Put(n, r)
		return r
	}
	panic(a.Err(a.KeyNotFound, n))
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

func read(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	s := a.AssertSequence(v)
	return e.ReadStr(c, a.ToStr(s))
}

func eval(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	v := args.First()
	return a.Eval(c, v)
}

func init() {
	Namespace.Put("def-builtin",
		a.NewFunction(defBuiltIn).WithMetadata(a.Properties{
			a.MetaSpecial: a.True,
		}).(a.Function),
	)

	RegisterBuiltIn("do", a.EvalBlock)
	RegisterBuiltIn("read", read)
	RegisterBuiltIn("eval", eval)
}
