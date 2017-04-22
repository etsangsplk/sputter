package api

// MetaMacro is the Metadata key that identifies a Function as being a Macro
var MetaMacro = NewKeyword("macro")

// Macro is a Function that can be used for Reader transformation
type Macro interface {
	Function
	DataMode() bool
}

type macro struct {
	Function
	dataMode bool
}

var macroMetadata = Metadata{MetaMacro: True}

// NewMacro instantiates a new Macro
func NewMacro(e SequenceProcessor) Macro {
	f := NewFunction(e).WithMetadata(macroMetadata).(Function)
	return &macro{
		Function: f,
		dataMode: true,
	}
}

func (m *macro) DataMode() bool {
	return m.dataMode
}

func (m *macro) Type() Name {
	return "macro"
}

// WithMetadata copies the Function with new Metadata
func (m *macro) WithMetadata(md Metadata) Annotated {
	return &macro{
		Function: m.Function.WithMetadata(md).(Function),
		dataMode: m.dataMode,
	}
}
