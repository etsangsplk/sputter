package api

var (
	// MetaMacro identifies a Function as being a Macro
	MetaMacro = NewKeyword("macro")

	// MetaSplicing identifies a Macro as requiring its result to be spliced
	MetaSplicing = NewKeyword("splicing")
)

var macroMetadata = Metadata{
	MetaMacro:    True,
	MetaType:     Name("macro"),
	MetaSplicing: False,
}

// NewMacro instantiates a new Macro
func NewMacro(e SequenceProcessor) Function {
	return NewFunction(e).WithMetadata(macroMetadata).(Function)
}

// IsMacro tests an Applicable as being marked a Macro
func IsMacro(a Applicable) bool {
	if an, ok := a.(Annotated); ok {
		if v, ok := an.Metadata().Get(MetaMacro); ok {
			return v == True
		}
	}
	return false
}
