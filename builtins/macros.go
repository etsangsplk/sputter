package builtins

import a "github.com/kode4food/sputter/api"

func defmacro(c a.Context, form a.Sequence) a.Value {
	a.AssertArity(form, 4)
	ns := a.GetContextNamespace(c)

	i := a.Iterate(form)
	i.Next() // skip the form name

	mv, _ := i.Next()
	mn := a.AssertUnqualified(mv).Name

	av, _ := i.Next()
	an := a.AssertSequence(av)
	fc := a.Count(an)
	b, _ := i.Next()

	m := &a.Macro{
		Name: mn,
		Data: true,
		Prep: func(c a.Context, form a.Sequence) a.Value {
			i := a.Iterate(form)
			i.Next() // skip the macro name
			a.AssertArity(i.Rest(), fc)
			return b
		},
	}

	ns.Put(mn, m)
	return m
}

func init() {
	registerMacro(&a.Macro{
		Name: "defmacro",
		Prep: defmacro,
		Data: true,
	})
}
