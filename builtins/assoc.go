package builtins

import a "github.com/kode4food/sputter/api"

const (
	assocName    = "assoc"
	isAssocName  = "assoc?"
	isMappedName = "mapped?"
)

type (
	assocFunction struct{ BaseBuiltIn }

	isAssociativeFunction struct{ a.BaseFunction }
	isMappedFunction      struct{ a.BaseFunction }
)

func (*assocFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToAssociative(args)
}

func (*isAssociativeFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.Associative); ok {
		return a.True
	}
	return a.False
}

func (*isMappedFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if _, ok := args.First().(a.MappedSequence); ok {
		return a.True
	}
	return a.False
}

func init() {
	var assoc *assocFunction
	var isAssociative *isAssociativeFunction
	var isMapped *isMappedFunction

	RegisterBuiltIn(assocName, assoc)
	RegisterSequencePredicate(isAssocName, isAssociative)
	RegisterSequencePredicate(isMappedName, isMapped)
}
