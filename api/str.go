package api

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// Str is the Sequence-compatible representation of string values
type Str string

var (
	escape = regexp.MustCompile(`[\\"]`)

	emptyStr Str
)

// First returns the first character of the Str
func (s Str) First() Value {
	if r, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return Str(r)
	}
	return Nil
}

// Rest returns a Str of all characters after the first
func (s Str) Rest() Sequence {
	if _, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return Str(s[w:])
	}
	return emptyStr
}

// Split returns the split form (First and Rest) of the Sequence
func (s Str) Split() (Value, Sequence, bool) {
	if r, w := utf8.DecodeRuneInString(string(s)); w > 0 {
		return Str(r), Str(s[w:]), true
	}
	return Nil, emptyStr, false
}

// IsSequence returns true if the Str is not empty
func (s Str) IsSequence() bool {
	return len(s) != 0
}

// Prepend prepends a Value to the beginning of the Str. If the Value
// is a single character, the resulting Str will be retained in native
// form, otherwise a List is returned.
func (s Str) Prepend(v Value) Sequence {
	if e, ok := v.(Str); ok && len(e) == 1 {
		return Str(e + s)
	}
	return s.list().Prepend(v)
}

func (s Str) list() List {
	c := []rune(string(s))
	r := EmptyList
	for i := len(c) - 1; i >= 0; i-- {
		r = r.Prepend(Str(c[i])).(List)
	}
	return r
}

// Conjoin appends a Value to the end of the Str. If the Value is a
// single character, the resulting Str will be retained in native
// form, otherwise a Vector is returned.
func (s Str) Conjoin(v Value) Sequence {
	if e, ok := v.(Str); ok && len(e) == 1 {
		return Str(s + e)
	}
	return s.vector().Conjoin(v)
}

func (s Str) vector() Vector {
	c := []rune(string(s))
	r := make(Values, len(c))
	for i := 0; i < len(c); i++ {
		r[i] = Str(c[i])
	}
	return NewVector(r...)
}

// Count returns the length of the Str
func (s Str) Count() int {
	return utf8.RuneCountInString(string(s))
}

// ElementAt returns the Character at the indexed position in the Str
func (s Str) ElementAt(index int) (Value, bool) {
	if index < 0 {
		return Nil, false
	}
	ns := string(s)
	p := 0
	for i := 0; i < index; i++ {
		if _, w := utf8.DecodeRuneInString(ns[p:]); w > 0 {
			p += w
		} else {
			return Nil, false
		}
	}
	if r, w := utf8.DecodeRuneInString(ns[p:]); w > 0 {
		return Str(r), true
	}
	return Nil, false
}

// Apply makes Str applicable
func (s Str) Apply(_ Context, args Sequence) Value {
	return IndexedApply(s, args)
}

// Str converts this Value into a Str
func (s Str) Str() Str {
	r := escape.ReplaceAllStringFunc(string(s), func(e string) string {
		return "\\" + e
	})
	return Str(`"` + r + `"`)
}

// MakeDumpStr takes a Value and attempts to spit out a bunch of info
func MakeDumpStr(v Value) Str {
	p := Str(fmt.Sprintf("%p", v))
	m := Object(Properties{InstanceKey: p})
	if t, ok := v.(Typed); ok {
		m = m.Child(Properties{TypeKey: t.Type()})
	}
	return m.Str()
}
