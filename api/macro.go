package api

import "bytes"

// Macro is a Function that can be used for Reader transformation
type Macro interface {
	Function
	IsMacro() bool
}

type basicMacro struct {
	Function
}

// NewMacro instantiates f new Macro
func NewMacro(e SequenceProcessor) Macro {
	return &basicMacro{
		Function: NewFunction(e),
	}
}

func (m *basicMacro) IsMacro() bool {
	return true
}

// WithMetadata copies the Function with new Metadata
func (m *basicMacro) WithMetadata(md Variables) Annotated {
	return &basicMacro{
		Function: m.Function.WithMetadata(md).(Function),
	}
}

func (m *basicMacro) String() string {
	var b bytes.Buffer
	b.WriteString("(macro ")
	b.WriteString(" :name " + String(m.Name()))
	doc := m.Documentation()
	if doc != "" {
		b.WriteString(" :doc " + String(doc))
	}
	b.WriteString(")")
	return b.String()
}
