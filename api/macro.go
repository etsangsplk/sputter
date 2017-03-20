package api

import (
	"bytes"
	"fmt"
)

// Macro is a Value that can be used to transform a Value
type Macro struct {
	Name Name
	Doc  string
	Prep SequenceProcessor
	Data bool
}

// Docstring makes Macro Documented
func (m *Macro) Docstring() string {
	return m.Doc
}

// Apply makes Macro Applicable
func (m *Macro) Apply(c Context, args Sequence) Value {
	return m.Prep(c, args)
}

func (m *Macro) String() string {
	var b bytes.Buffer
	b.WriteString("(macro")
	if m.Name != "" {
		b.WriteString(" :name " + String(m.Name))
	}
	b.WriteString(fmt.Sprintf(" :instance %p", &m))
	if m.Doc != "" {
		b.WriteString(" :doc " + String(m.Doc))
	}
	b.WriteString(")")
	return b.String()
}
