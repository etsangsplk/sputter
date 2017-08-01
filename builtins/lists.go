package builtins

import a "github.com/kode4food/sputter/api"

type listFunction struct{ a.ReflectedFunction }

func (f *listFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToList(args)
}

func isList(v a.Value) bool {
	if _, ok := v.(a.List); ok {
		return true
	}
	return false
}

func init() {
	var list *listFunction

	RegisterBaseFunction("list", list)
	RegisterSequencePredicate("list?", isList)
}
