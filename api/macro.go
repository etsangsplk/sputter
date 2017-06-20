package api

var (
	// MetaMacro identifies a Function as being a Macro
	MetaMacro = NewKeyword("macro")
	// MetaSpecial identifies a Macro as being a special form
	MetaSpecial = NewKeyword("special-form")
)

var macroMetadata = Metadata{
	MetaMacro:   True,
	MetaType:    Name("macro"),
	MetaSpecial: False,
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

// IsSpecialForm tests a Macro as being marked as a special form
func IsSpecialForm(a Applicable) bool {
	if !IsMacro(a) {
		return false
	}
	an, _ := a.(Annotated)
	if v, ok := an.Metadata().Get(MetaSpecial); ok {
		return v == True
	}
	return false
}
