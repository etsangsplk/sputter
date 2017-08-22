package api

// Values turns an array of Values into a basic Sequence
type Values []Value

var emptyValues = Values{}

// IsSequence returns whether or not the array is empty
func (a Values) IsSequence() bool {
	return len(a) > 0
}

// First returns the first element of the array
func (a Values) First() Value {
	if len(a) > 0 {
		return a[0]
	}
	return Nil
}

// Rest returns the remaining elements of the array
func (a Values) Rest() Sequence {
	if len(a) > 1 {
		return a[1:]
	}
	return emptyValues
}

// Split returns the components of the array (first, rest, isSeq?)
func (a Values) Split() (Value, Sequence, bool) {
	lv := len(a)
	if lv > 1 {
		return a[0], a[1:], true
	} else if lv == 1 {
		return a[0], emptyValues, true
	}
	return Nil, emptyValues, false
}

// Prepend inserts an element at the beginning of the Array
func (a Values) Prepend(p Value) Sequence {
	return append(Values{p}, a...)
}

// Str converts the Array into a Str
func (a Values) Str() Str {
	return MakeSequenceStr(a)
}
