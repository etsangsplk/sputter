package api

import "bytes"

// SequenceToList takes any sequence and converts it to a List
func SequenceToList(s Sequence) List {
	if l, ok := s.(List); ok {
		return l
	}
	if c, ok := s.(Counted); ok {
		res := make(Values, c.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return NewList(res...)
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
		res := make(vector, c.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return res
	}
	return uncountedToVector(s)
}

func uncountedToVector(s Sequence) Vector {
	return vector(uncountedToArray(s))
}

func uncountedToArray(s Sequence) Values {
	res := Values{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, f)
	}
	return res
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
			k, i, _ = i.Split()
			v, i, _ = i.Split()
			r[idx] = NewVector(k, v)
		}
		return associative(r)
	}
	return uncountedToAssociative(s)
}

func uncountedToAssociative(s Sequence) Associative {
	res := []Vector{}
	var v Value
	for k, r, ok := s.Split(); ok; k, r, ok = r.Split() {
		if v, r, ok = r.Split(); ok {
			res = append(res, NewVector(k, v))
		} else {
			panic(ErrStr(ExpectedPair))
		}
	}
	return associative(res)
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
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if f == Nil {
			continue
		}
		buf.WriteString(string(MakeStr(f)))
	}
	return Str(buf.String())
}
