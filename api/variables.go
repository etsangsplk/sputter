package api

// Name is a Variable name
type Name string

// Value is the generic interface for all 'Values'
type Value interface {
}

// Named is the generic interface for Values that are named
type Named interface {
	Name() Name
}

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

// Name makes Name Named
func (n Name) Name() Name {
	return n
}
