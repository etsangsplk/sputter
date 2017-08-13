package api

import "bytes"

// SequenceToList takes any sequence and converts it to a List
func SequenceToList(s Sequence) List {
	if l, ok := s.(List); ok {
		return l
	}
	if c, ok := s.(Counted); ok {
		r := make([]Value, c.Count())
		var t Value
		for i, idx := s, 0; i.IsSequence(); {
			t, i = i.Split()
			r[idx] = t
			idx++
		}
		return NewList(r...)
	}
	return uncountedToList(s)
}

func uncountedToList(s Sequence) List {
	return NewList(uncountedToArray(s)...)
}

// SequenceToVector takes any sequence and converts it to a Vector
func SequenceToVector(s Sequence) Vector {
	if v, ok := s.(Vector); ok {
		return v
	}
	if c, ok := s.(Counted); ok {
		r := make([]Value, c.Count())
		var t Value
		for i, idx := s, 0; i.IsSequence(); {
			t, i = i.Split()
			r[idx] = t
			idx++
		}
		return vector(r)
	}
	return uncountedToVector(s)
}

func uncountedToVector(s Sequence) Vector {
	return vector(uncountedToArray(s))
}

func uncountedToArray(s Sequence) []Value {
	r := []Value{}
	var t Value
	for i := s; i.IsSequence(); {
		t, i = i.Split()
		r = append(r, t)
	}
	return r
}

// SequenceToAssociative takes any sequence and converts it to an associative
func SequenceToAssociative(s Sequence) Associative {
	if a, ok := s.(Associative); ok {
		return a
	}
	if c, ok := s.(Counted); ok {
		l := c.Count()
		if l%2 != 0 {
			panic(ErrStr(ExpectedPair))
		}
		ml := l / 2
		r := make([]Vector, ml)
		i := s
		var k, v Value
		for idx := 0; idx < ml; idx++ {
			k, i = i.Split()
			v, i = i.Split()
			r[idx] = NewVector(k, v)
		}
		return associative(r)
	}
	return uncountedToAssociative(s)
}

func uncountedToAssociative(s Sequence) Associative {
	r := []Vector{}
	var k, v Value
	for i := s; i.IsSequence(); {
		k, i = i.Split()
		if i.IsSequence() {
			v, i = i.Split()
			r = append(r, NewVector(k, v))
		} else {
			panic(ErrStr(ExpectedPair))
		}
	}
	return associative(r)
}

// MakeStr converts a Value to a Str if it's not already one
func MakeStr(v Value) Str {
	if s, ok := v.(Str); ok {
		return s
	}
	return v.Str()
}

// SequenceToStr takes any sequence and attempts to convert it to a Str
func SequenceToStr(s Sequence) Str {
	if st, ok := s.(Str); ok {
		return st
	}
	var buf bytes.Buffer
	var t Value
	for i := s; i.IsSequence(); {
		t, i = i.Split()
		if t == Nil {
			continue
		}
		buf.WriteString(string(MakeStr(t)))
	}
	return Str(buf.String())
}
