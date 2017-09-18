package builtins

import a "github.com/kode4food/sputter/api"

const (
	listName   = "list"
	isListName = "list?"
)

type listFunction struct{ BaseBuiltIn }

func (*listFunction) Apply(_ a.Context, args a.Sequence) a.Value {
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

	RegisterBuiltIn(listName, list)
	RegisterSequencePredicate(isListName, isList)
}
