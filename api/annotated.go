package api

const (
	// ExpectedAnnotated is thrown if a Value is not Annotated
	ExpectedAnnotated = "value is not annotated with metadata"

	// MetaName is the Metadata key for a Value's Name
	MetaName = Name("name")

	// MetaType is the Metadata key for a Value's Type
	MetaType = Name("type")

	// MetaDoc is the Metadata key for Documentation Strings
	MetaDoc = Name("doc")
)

// Annotated is implemented if a Value is Annotated with Metadata
type Annotated interface {
	Metadata() Variables
	WithMetadata(md Variables) Annotated
}

func AssertAnnotated(v Value) Annotated {
	if a, ok := v.(Annotated); ok {
		return a
	}
	panic(ExpectedAnnotated)
}
