package builtins

import a "github.com/kode4food/sputter/api"

func defmacro(c a.Context, form a.Sequence) a.Value {
	a.AssertArity(form, 4)
	ns := a.GetContextNamespace(c)

	i := form.Iterate()
	i.Next() // skip the form name

	mv, _ := i.Next()
	mn := a.AssertUnqualified(mv).Name

	av, _ := i.Next()
	an := a.AssertSequence(av)
	fc := a.Count(an)
	b, _ := i.Next()

	m := &a.Function{
		Name: mn,
		Prepare: func(c a.Context, form a.Sequence) a.Value {
			i := form.Iterate()
			i.Next() // skip the macro name
			a.AssertArity(i.Rest(), fc)
			return b
		},
	}
	
	ns.Put(mn, m)
	return m
}

func init() {
	registerFunction(&a.Function{
		Name:    "defmacro",
		Prepare: defmacro,
		Data:    true,
	})
}
