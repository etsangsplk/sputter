package builtins

import a "github.com/kode4food/sputter/api"

const (
	assocName    = "assoc"
	isAssocName  = "assoc?"
	isMappedName = "mapped?"
)

type assocFunction struct{ BaseBuiltIn }

func (*assocFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	return a.SequenceToAssociative(args)
}

func isAssociative(v a.Value) bool {
	if _, ok := v.(a.Associative); ok {
		return true
	}
	return false
}

func isMapped(v a.Value) bool {
	if _, ok := v.(a.MappedSequence); ok {
		return true
	}
	return false
}

func init() {
	var assoc *assocFunction

	RegisterBuiltIn(assocName, assoc)
	RegisterSequencePredicate(isAssocName, isAssociative)
	RegisterSequencePredicate(isMappedName, isMapped)
}
