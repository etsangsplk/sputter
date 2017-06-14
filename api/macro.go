package api

// MetaMacro identifies a Function as being a Macro
var MetaMacro = NewKeyword("macro")

var macroMetadata = Metadata{
	MetaMacro: True,
	MetaType:  Name("macro"),
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
