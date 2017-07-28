package api

import "bytes"

// ToList takes any sequence and converts it to a List
func ToList(s Sequence) List {
	if l, ok := s.(List); ok {
		return l
	}
	if c, ok := s.(Counted); ok {
		r := make([]Value, c.Count())
		for i, idx := s, 0; i.IsSequence(); i = i.Rest() {
			r[idx] = i.First()
			idx++
		}
		return NewList(r...)
	}
	return uncountedToList(s)
}

func uncountedToList(s Sequence) List {
	return NewList(uncountedToArray(s)...)
}

// ToVector takes any sequence and converts it to a Vector
func ToVector(s Sequence) Vector {
	if v, ok := s.(Vector); ok {
		return v
	}
	if c, ok := s.(Counted); ok {
		r := make([]Value, c.Count())
		for i, idx := s, 0; i.IsSequence(); i = i.Rest() {
			r[idx] = i.First()
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
	for i := s; i.IsSequence(); i = i.Rest() {
		r = append(r, i.First())
	}
	return r
}

// ToAssociative takes any sequence and converts it to an associative
func ToAssociative(s Sequence) Associative {
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
		for idx := 0; idx < ml; idx++ {
			k := i.First()
			i = i.Rest()

			v := i.First()
			i = i.Rest()

			r[idx] = NewVector(k, v)
		}
		return associative(r)
	}
	return uncountedToAssociative(s)
}

func uncountedToAssociative(s Sequence) Associative {
	r := []Vector{}
	for i := s; i.IsSequence(); i = i.Rest() {
		k := i.First()
		i = i.Rest()
		if i.IsSequence() {
			v := i.First()
			r = append(r, NewVector(k, v))
		} else {
			panic(ErrStr(ExpectedPair))
		}
	}
	return associative(r)
}

// ToStr takes any sequence and attempts to convert it to a Str
func ToStr(s Sequence) Str {
	if st, ok := s.(Str); ok {
		return st
	}
	var buf bytes.Buffer
	for i := s; i.IsSequence(); i = i.Rest() {
		v := i.First()
		if v == Nil {
			continue
		}
		buf.WriteString(string(MakeStr(v)))
	}
	return Str(buf.String())
}

// ToReaderStr takes a sequence and converts it to a readable Str
func ToReaderStr(s Sequence) Str {
	if st, ok := s.(Str); ok {
		return st
	}
	var buf bytes.Buffer
	if s.IsSequence() {
		buf.WriteString(string(s.First().Str()))
	}
	for i := s.Rest(); i.IsSequence(); i = i.Rest() {
		buf.WriteString(" ")
		buf.WriteString(string(i.First().Str()))
	}
	return Str(buf.String())
}
