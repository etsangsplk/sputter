package api

import "sync"

// AlreadyBound is thrown when an attempt is made to rebind a Name
const AlreadyBound = "symbol is already bound in this scope: %s"

type (
	// Variables represents a mapping from Name to Value
	Variables map[Name]Value

	// WriteOnceVariables implements Variables with write-once semantics
	WriteOnceVariables struct {
		Variables
		sync.RWMutex
	}
)

// Get retrieves a variable by name
func (v Variables) Get(n Name) (Value, bool) {
	if r, ok := v[n]; ok {
		return r, ok
	}
	return Nil, false
}

// Has checks for the existence of a variable and returns its context
func (v Variables) Has(n Name) (Context, bool) {
	if _, ok := v[n]; ok {
		return v, ok
	}
	return v, false
}

// Put sets a variable by name
func (v Variables) Put(n Name, val Value) {
	v[n] = val
}

// Delete removes a variable by name
func (v Variables) Delete(n Name) {
	delete(v, n)
}

// Str converts this Value into a Str
func (v Variables) Str() Str {
	return MakeDumpStr(v)
}

// Get retrieves a variable by name
func (w *WriteOnceVariables) Get(n Name) (Value, bool) {
	w.RLock()
	defer w.RUnlock()
	return w.Variables.Get(n)
}

// Has checks for the existence of a variable and returns its context
func (w *WriteOnceVariables) Has(n Name) (Context, bool) {
	w.RLock()
	defer w.RUnlock()
	if _, ok := w.Variables.Has(n); ok {
		return w, true
	}
	return w, false
}

// Delete removes a variable by name
func (w *WriteOnceVariables) Delete(n Name) {
	w.Lock()
	defer w.Unlock()
	w.Variables.Delete(n)
}

// Put sets a variable by name if it hasn't already been set
func (w *WriteOnceVariables) Put(n Name, v Value) {
	w.Lock()
	defer w.Unlock()
	if _, ok := w.Variables.Has(n); ok {
		panic(ErrStr(AlreadyBound, n))
	}
	w.Variables.Put(n, v)
}
