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

// IsMacro tests an Applicable as being marked a Macro and is a special form
func IsMacro(a Applicable) (bool, bool) {
	if an, ok := a.(Annotated); ok {
		md := an.Metadata()
		return md.IsTrue(MetaMacro), md.IsTrue(MetaSpecial)
	}
	return false, false
}

// IsSpecialForm tests an Applicable as being marked a special form
func IsSpecialForm(a Applicable) bool {
	if an, ok := a.(Annotated); ok {
		return an.Metadata().IsTrue(MetaSpecial)
	}
	return false
}
