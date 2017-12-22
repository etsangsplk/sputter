package api

// AlreadyBound is thrown when an attempt is made to rebind a Name
const AlreadyBound = "symbol is already bound in this scope: %s"

// Variables represents a mapping from Name to Value
type Variables map[Name]Value

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
