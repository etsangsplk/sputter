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

func (*listFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToList(args)
}

func (*isListFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.List); ok {
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
