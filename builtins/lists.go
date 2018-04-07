package builtins

import a "github.com/kode4food/sputter/api"

const (
	listName   = "list"
	isListName = "is-list"
)

type (
	listFunction   struct{ BaseBuiltIn }
	isListFunction struct{ BaseBuiltIn }
)

func (*listFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return a.SequenceToList(args)
}

func (*isListFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(*a.List); ok {
		return a.True
	}
	return a.False
}

func init() {
	var list *listFunction
	var isList *isListFunction

	RegisterBuiltIn(listName, list)
	RegisterBuiltIn(isListName, isList)
}
