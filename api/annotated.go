package api

import (
	"bytes"
	"fmt"
)

// ExpectedAnnotated is thrown if a Value is not Annotated
const ExpectedAnnotated = "value does not support annotation: %s"

var (
	// MetaName is the Metadata key for a Value's Name
	MetaName = NewKeyword("name")

	// MetaType is the Metadata key for a Value's Type
	MetaType = NewKeyword("type")

	// MetaDoc is the Metadata key for Documentation Strings
	MetaDoc = NewKeyword("doc")
)

// Metadata stores metadata about an Annotated Value
type Metadata map[Value]Value

// Annotated is implemented if a Value is Annotated with Metadata
type Annotated interface {
	Metadata() Metadata
	WithMetadata(md Metadata) Annotated
}

// Merge merges two Metadata sets into a new one
func (v Metadata) Merge(nv Metadata) Metadata {
	r := make(Metadata)
	for k, v := range v {
		r[k] = v
	}
	for k, v := range nv {
		r[k] = v
	}
	return r
}

func (v Metadata) String() string {
	var b bytes.Buffer
	c := false
	b.WriteString("{")
	for k, v := range v {
		if c {
			b.WriteString(", ")
		} else {
			c = true
		}
		b.WriteString(String(k))
		b.WriteString(" ")
		b.WriteString(String(v))
	}
	b.WriteString("}")
	return b.String()
}

// AssertAnnotated will cast a Value to Annotated or die trying
func AssertAnnotated(v Value) Annotated {
	if a, ok := v.(Annotated); ok {
		return a
	}
	panic(fmt.Sprintf(ExpectedAnnotated, String(v)))
}
