package api

import (
	"fmt"
	"regexp"
)

// Str is the Sequence-compatible representation of string values
type Str string

var escape = regexp.MustCompile(`\\|\"`)

// First returns the first character of the Str
func (s Str) First() Value {
	c := []rune(string(s))
	return Str(c[0])
}

// Rest returns a Str of all characters after the first
func (s Str) Rest() Sequence {
	c := []rune(string(s))
	return Str(c[1:])
}

// IsSequence returns true if the Str is not empty
func (s Str) IsSequence() bool {
	return len(s) != 0
}

// Prepend prepends a Value to the beginning of the Str. If the Value
// is a single character, the resulting Str will be retained in native
// form, otherwise a List is returned.
func (s Str) Prepend(v Value) Sequence {
	if e, ok := v.(Str); ok && e.Count() == 1 {
		return Str(e + s)
	}
	return s.list().Prepend(v)
}

func (s Str) list() *List {
	c := []rune(string(s))
	r := EmptyList
	for i := len(c) - 1; i >= 0; i-- {
		r = r.Prepend(Str(c[i])).(*List)
	}
	return r
}

// Conjoin appends a Value to the end of the Str. If the Value is a
// single character, the resulting Str will be retained in native
// form, otherwise a Vector is returned.
func (s Str) Conjoin(v Value) Sequence {
	if e, ok := v.(Str); ok && e.Count() == 1 {
		return Str(s + e)
	}
	return s.vector().Conjoin(v)
}

func (s Str) vector() Vector {
	c := []rune(string(s))
	r := make(Vector, len(c))
	for i := 0; i < len(c); i++ {
		r[i] = Str(c[i])
	}
	return r
}

// Count returns the length of the Str
func (s Str) Count() int {
	c := []rune(string(s))
	return len(c)
}

// Get returns the Character at the indexed position in the Str
func (s Str) Get(index int) (Value, bool) {
	if index >= 0 && index < s.Count() {
		c := []rune(string(s))
		return Str(c[index : index+1]), true
	}
	return Nil, false
}

// Apply makes Str applicable
func (s Str) Apply(c Context, args Sequence) Value {
	return IndexedApply(s, c, args)
}

// Str converts this Value into a Str
func (s Str) Str() Str {
	r := escape.ReplaceAllStringFunc(string(s), func(e string) string {
		return "\\" + e
	})
	return Str(`"` + r + `"`)
}

// MakeDumpStr takes a Value and attempts to spit out a bunch of metadata
func MakeDumpStr(v Value) Str {
	m := Metadata{}
	if n, ok := v.(Named); ok {
		m = m.Merge(Metadata{MetaName: n.Name()})
	}
	if t, ok := v.(Typed); ok {
		m = m.Merge(Metadata{MetaType: t.Type()})
	}
	p := Str(fmt.Sprintf("%p", &v))
	m = m.Merge(Metadata{MetaInstance: p})
	if a, ok := v.(Annotated); ok {
		m = m.Merge(Metadata{MetaMeta: a.Metadata()})
	}
	return m.Str()
}
