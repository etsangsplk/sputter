package builtins

import a "github.com/kode4food/sputter/api"

func defineMacro(closure a.Context, d *functionDefinition) a.Function {
	an := makeArgProcessor(closure, d.args)
	db := d.body

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		l := an(c, args)
		return a.EvalSequence(l, db)
	}).WithMetadata(d.meta).(a.Function)
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	m := defineMacro(c, fd)
	a.GetContextNamespace(c).Put(m.Name(), m)
	return m
}

func init() {
	registerAnnotated(
		a.NewFunction(defmacro).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defmacro"),
		}),
	)
}
