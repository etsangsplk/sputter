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

func RegisterBuiltIn(n a.Name, proc a.SequenceProcessor) {
	builtInFuncs[n] = proc
}

// GetBuiltIn Returns a registered built-in function
func GetBuiltIn(n a.Name) (a.SequenceProcessor, bool) {
	if v, ok := builtInFuncs[n]; ok {
		return v, true
	}
	return nil, false
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

func defBuiltIn(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	if f, ok := builtInFuncs[n]; ok {
		var md a.Object = toProperties(a.ToAssociative(args.Rest()))
		if s, ok := md.Get(a.MetaName); ok {
			n = a.AssertUnqualified(s).Name()
			md = md.Child(a.Properties{
				a.MetaName: n,
			})
		}
		md = loadDocumentation(md)

		var r a.Function
		if a.IsTrue(md, a.MetaMacro) {
			r = a.NewMacro(f)
		} else {
			r = a.NewFunction(f)
		}

		r = r.WithMetadata(md).(a.Function)
		a.GetContextNamespace(c).Put(n, r)
		return r
	}
	panic(a.Err(a.KeyNotFound, n))
}

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalBlock(c, args)
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

	RegisterBuiltIn("do", do)
	RegisterBuiltIn("read", read)
	RegisterBuiltIn("eval", eval)
}
