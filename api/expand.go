package api

// MacroExpand1 performs a single macro expansion
func MacroExpand1(c Context, v Value) (Value, bool) {
	if l, ok := v.(*List); ok {
		if f, r, ok := l.Split(); ok {
			if s, ok := f.(Symbol); ok {
				if sr, ok := s.Resolve(c); ok {
					if ap, ok := sr.(Applicable); ok {
						if IsMacro(ap) && !IsSpecialForm(ap) {
							return Apply(c, ap, r), true
						}
					}
				}
			}
		}
	}
	return v, false
}

// MacroExpand repeatedly performs a macro expansion until no more can occur
func MacroExpand(c Context, v Value) (Value, bool) {
	var ok bool
	r := v
	for i := 0; ; i++ {
		if r, ok = MacroExpand1(c, r); ok {
			continue
		}
		return r, i > 0
	}
}

// MacroExpandAll attempts to recursively expand the specified Value
func MacroExpandAll(c Context, v Value) Value {
	ex, _ := MacroExpand(c, v)
	if s, ok := ex.(Sequence); ok {
		return expandSequence(c, s)
	}
	return ex
}

func expandSequence(c Context, s Sequence) Value {
	if st, ok := s.(Str); ok {
		return st
	}
	if l, ok := s.(*List); ok {
		ex := NewList(expandElements(c, l)...)
		return ex
	}
	if v, ok := s.(Vector); ok {
		return NewVector(expandElements(c, v)...)
	}
	if as, ok := s.(Associative); ok {
		return expandAssociative(c, as)
	}
	return s
}

func expandAssociative(c Context, as Associative) Value {
	var res []Vector
	for f, r, ok := as.Split(); ok; f, r, ok = r.Split() {
		e := f.(Vector)
		k, _ := e.ElementAt(0)
		v, _ := e.ElementAt(1)
		res = append(res, Vector{
			MacroExpandAll(c, k),
			MacroExpandAll(c, v),
		})
	}
	return NewAssociative(res...)
}

func expandElements(c Context, s Sequence) Vector {
	res := Vector{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, MacroExpandAll(c, f))
	}
	return res
}
