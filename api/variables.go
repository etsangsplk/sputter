package api

import "bytes"

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

// Merge merges two Variables sets into a new one
func (v Variables) Merge(nv Variables) Variables {
	r := make(Variables)
	for k, v := range v {
		r[k] = v
	}
	for k, v := range nv {
		r[k] = v
	}
	return r
}

func (v Variables) String() string {
	var b bytes.Buffer
	c := false
	b.WriteString("{")
	for k, v := range v {
		if c {
			b.WriteString(", ")
		} else {
			c = true
		}
		b.WriteString(":")
		b.WriteString(String(k))
		b.WriteString(" ")
		b.WriteString(String(v))
	}
	b.WriteString("}")
	return b.String()
}
