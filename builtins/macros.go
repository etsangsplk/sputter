package builtins

import a "github.com/kode4food/sputter/api"

func defmacro(c a.Context, args a.Sequence) a.Value {
	fd := getFunctionDefinition(c, args)
	m := &a.Macro{Function: defineFunction(fd)}
	a.GetContextNamespace(fd.closure).Put(fd.name, m)
	return m
}

func init() {
	registerMacro(&a.Function{Name: "defmacro", Exec: defmacro})
}
