package api

import "fmt"

// Macro is a Value that can be used to transform a Value
type Macro struct {
	Name Name
	Prep SequenceProcessor
	Data bool
}

// Apply makes Macro Applicable
func (m *Macro) Apply(c Context, args Sequence) Value {
	return m.Prep(c, args)
}

func (m *Macro) String() string {
	if m.Name != "" {
		return "(macro :name " + string(m.Name) + ")"
	}
	return fmt.Sprintf("(macro :addr %p)", &m)
}
