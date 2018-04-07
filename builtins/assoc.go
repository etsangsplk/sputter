package builtins

import a "github.com/kode4food/sputter/api"

const (
	assocName    = "assoc"
	isAssocName  = "is-assoc"
	isMappedName = "is-mapped"
)

type (
	assocFunction    struct{ BaseBuiltIn }
	isAssocFunction  struct{ BaseBuiltIn }
	isMappedFunction struct{ BaseBuiltIn }
)

func (*assocFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return a.SequenceToAssociative(args)
}

func (*isAssocFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(a.Associative); ok {
		return a.True
	}
	return a.False
}

func (*isMappedFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(a.MappedSequence); ok {
		return a.True
	}
	return a.False
}

func init() {
	var assoc *assocFunction
	var isAssoc *isAssocFunction
	var isMapped *isMappedFunction

	RegisterBuiltIn(assocName, assoc)
	RegisterBuiltIn(isAssocName, isAssoc)
	RegisterBuiltIn(isMappedName, isMapped)
}
