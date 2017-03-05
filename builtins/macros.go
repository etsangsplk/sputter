package builtins

import a "github.com/kode4food/sputter/api"

func defmacro(_ a.Context, form a.Sequence) a.Value {
	a.AssertArity(form, 4)
	g := a.GetNamespace(a.UserDomain)

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
	a.PutFunction(g, m)
	return m
}

func init() {
	a.PutFunction(Context, &a.Function{
		Name:    "defmacro",
		Prepare: defmacro,
		Data:    true,
	})
}
