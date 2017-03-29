package builtins

import a "github.com/kode4food/sputter/api"

func defineMacro(closure a.Context, d *functionDefinition) a.Macro {
	an := argNames(d.argNames)
	ac := len(an)
	db := d.body

	return a.NewMacro(func(c a.Context, args a.Sequence) a.Value {
		a.AssertArity(args, ac)
		l := a.ChildContext(closure)
		i := a.Iterate(args)
		for _, n := range an {
			v, _ := i.Next()
			l.Put(n, a.Eval(c, v))
		}
		return a.EvalSequence(l, db)
	}).WithMetadata(a.Metadata{
		a.MetaName: d.name,
		a.MetaDoc:  d.doc,
	}).(a.Macro)
}

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(args)
	m := defineMacro(c, fd)
	a.GetContextNamespace(c).Put(fd.name, m)
	return m
}

func init() {
	registerAnnotated(
		a.NewMacro(defmacro).WithMetadata(a.Metadata{
			a.MetaName: a.Name("defmacro"),
		}),
	)
}
