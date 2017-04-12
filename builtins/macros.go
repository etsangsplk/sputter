package builtins

import a "github.com/kode4food/sputter/api"

func defineMacro(closure a.Context, d *functionDefinition) a.Macro {
	an := argNames(d.args)
	ac := len(an)
	db := d.body

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, ac)
		l := a.ChildContext(closure)
		i := args
		for _, n := range an {
			l.Put(n, a.Eval(c, i.First()))
			i = i.Rest()
		}
		return a.EvalSequence(l, db)
	}).WithMetadata(d.meta).(a.Macro)
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	m := defineMacro(c, fd)
	a.GetContextNamespace(c).Put(m.Name(), m)
	return m
}

func init() {
	registerAnnotated(
		a.NewMacro(defmacro).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defmacro"),
		}),
	)
}
