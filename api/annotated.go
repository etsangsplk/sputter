package api

import "bytes"

// ExpectedAnnotated is thrown if a Value is not Annotated
const ExpectedAnnotated = "value does not support annotation: %s"

var (
	// MetaName is the Metadata key for a Value's Name
	MetaName = NewKeyword("name")

	// MetaType is the Metadata key for a Value's Type
	MetaType = NewKeyword("type")

	// MetaMeta is the Metadata key for a Value's Metadata
	MetaMeta = NewKeyword("meta")

	// MetaDoc is the Metadata key for Documentation Strings
	MetaDoc = NewKeyword("doc")

	// MetaArgs is the Metadata key for a Function's arguments
	MetaArgs = NewKeyword("args")

	// MetaInstance is the Metadata key for a Value's instance ID
	MetaInstance = NewKeyword("instance")

	// MetaClosure marks a Value as requiring a Closure
	MetaClosure = NewKeyword("closure")
)

// Metadata stores metadata about an Annotated Value
type Metadata map[Value]Value

// Annotated is implemented if a Value is Annotated with Metadata
type Annotated interface {
	Metadata() Metadata
	WithMetadata(md Metadata) Annotated
}

// Merge merges two Metadata sets into a new one
func (m Metadata) Merge(nv Metadata) Metadata {
	if len(m) == 0 {
		return nv
	}
	r := make(Metadata)
	for k, v := range m {
		r[k] = v
	}
	for k, v := range nv {
		r[k] = v
	}
	return r
}

// Get returns the Value corresponding to the key in the Metadata
func (m Metadata) Get(key Value) (Value, bool) {
	if r, ok := m[key]; ok {
		return r, true
	}
	return Nil, false
}

// Eval is self-evaluating
func (m Metadata) Eval(_ Context) Value {
	return m
}

// Str converts this Value into a Str
func (m Metadata) Str() Str {
	var b bytes.Buffer
	c := false
	b.WriteString("{")
	for k, v := range m {
		if c {
			b.WriteString(", ")
		} else {
			c = true
		}
		b.WriteString(string(k.Str()))
		b.WriteString(" ")
		b.WriteString(string(v.Str()))
	}
	b.WriteString("}")
	return Str(b.String())
}

// AssertAnnotated will cast a Value to Annotated or die trying
func AssertAnnotated(v Value) Annotated {
	if a, ok := v.(Annotated); ok {
		return a
	}
	panic(Err(ExpectedAnnotated, v))
}
