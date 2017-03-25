package api

import "bytes"

// Macro is a Function that can be used for Reader transformation
type Macro struct {
	*Function
	Prep SequenceProcessor
}

func (m *Macro) String() string {
	var b bytes.Buffer
	b.WriteString("(macro ")
	b.WriteString(" :name " + String(m.Name))
	if m.Doc != "" {
		b.WriteString(" :doc " + String(m.Doc))
	}
	b.WriteString(")")
	return b.String()
}
