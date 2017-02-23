package builtins

import a "github.com/kode4food/sputter/api"

const defaultVarsSize = 128

// Context is a special Context of built-in identifiers
var Context = &builtInsContext{make(a.Variables, defaultVarsSize)}

type builtInsContext struct {
	vars a.Variables
}

func (b *builtInsContext) Globals() a.Context {
	return b
}

func (b *builtInsContext) Get(n a.Name) (a.Value, bool) {
	if v, ok := b.vars[n]; ok {
		return v, true
	}
	return a.Nil, false
}

func (b *builtInsContext) Put(n a.Name, v a.Value) a.Context {
	b.vars[n] = v
	return b
}

func quote(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	i := args.Iterate()
	v, _ := i.Next()
	return v
}

func do(c a.Context, args a.Sequence) a.Value {
	return a.EvalSequence(c, args)
}

func init() {
	Context.Put("nil", a.Nil)
	Context.Put("true", a.True)
	Context.Put("false", a.False)

	a.PutFunction(Context, &a.Function{Name: "quote", Exec: quote})
	a.PutFunction(Context, &a.Function{Name: "do", Exec: do})
}
