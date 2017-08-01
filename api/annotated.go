package api

import d "github.com/kode4food/sputter/docstring"

// ExpectedAnnotated is thrown if a Value is not Annotated
const ExpectedAnnotated = "value does not support annotation: %s"

var (
	// NameKey is the Metadata key for a Value's Name
	NameKey = NewKeyword("name")

	// TypeKey is the Metadata key for a Value's Type
	TypeKey = NewKeyword("type")

	// DocKey is the Metadata key for Documentation Strings
	DocKey = NewKeyword("doc")

	// DocAssetKey is the Metadata key for Asset Strings
	DocAssetKey = NewKeyword("doc-asset")

	// ArgsKey is the Metadata key for a Function's arguments
	ArgsKey = NewKeyword("args")

	// InstanceKey is the Metadata key for a Value's instance ID
	InstanceKey = NewKeyword("instance")

	// Undocumented is the default documentation for a symbol
	Undocumented = Str("this symbol is not documented")
)

type (
	// Annotated is implemented if a Value is Annotated with Metadata
	Annotated interface {
		Metadata() Object
		WithMetadata(md Object) AnnotatedValue
	}

	// AnnotatedValue is returned by the Child call
	AnnotatedValue interface {
		Annotated
		Value
	}
)

// IsTrue tests whether or not the specified key has a True value
func IsTrue(o Object, key Value) bool {
	if r, ok := o.Get(key); ok {
		return r == True
	}
	return false
}

// GetDocumentation retrieves the doc or doc-asset for an Annotated Value
func GetDocumentation(md Object) Str {
	if v, ok := md.Get(DocKey); ok {
		return MakeStr(v)
	}
	if v, ok := md.Get(DocAssetKey); ok {
		k := string(MakeStr(v))
		if d.Exists(k) {
			return Str(d.Get(k))
		}
	}
	return Undocumented
}

// AssertAnnotated will cast a Value to Annotated or die trying
func AssertAnnotated(v Value) Annotated {
	if a, ok := v.(Annotated); ok {
		return a
	}
	panic(ErrStr(ExpectedAnnotated, v))
}
