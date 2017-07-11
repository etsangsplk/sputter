package api

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
)

type (
	// Annotated is implemented if a Value is Annotated with Metadata
	Annotated interface {
		Metadata() Object
		WithMetadata(md Object) Annotated
	}
)

// IsTrue tests whether or not the specified key has a True value
func IsTrue(o Object, key Value) bool {
	if r, ok := o.Get(key); ok {
		return r == True
	}
	return false
}

// AssertAnnotated will cast a Value to Annotated or die trying
func AssertAnnotated(v Value) Annotated {
	if a, ok := v.(Annotated); ok {
		return a
	}
	panic(Err(ExpectedAnnotated, v))
}
