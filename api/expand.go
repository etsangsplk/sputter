package api

type (
	specialForm struct {
		List
		fn   Applicable
		args Vector
	}

	evaluatingForm struct {
		List
		fn   Applicable
		args Vector
	}
)

// MacroExpand1 performs a single macro expansion
func MacroExpand1(c Context, v Value) (Value, bool) {
	if l, ok := v.(List); ok {
		if f, r, ok := l.Split(); ok {
			if s, ok := f.(Symbol); ok {
				if sr, ok := s.Resolve(c); ok {
					if a, ok := sr.(Applicable); ok {
						if IsMacro(a) && !IsSpecialForm(a) {
							return Apply(c, a, r), true
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
	if l, ok := s.(List); ok {
		ex := NewList(expandElements(c, l)...)
		if f, ok := MakeForm(ex); ok {
			return f
		}
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
	for f, r, ok := as.(Sequence).Split(); ok; f, r, ok = r.Split() {
		e := f.(Vector)
		k, _ := e.ElementAt(0)
		v, _ := e.ElementAt(1)
		res = append(res, NewVector(
			MacroExpandAll(c, k),
			MacroExpandAll(c, v),
		))
	}
	return NewAssociative(res...)
}

func expandElements(c Context, s Sequence) Values {
	res := Values{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, MacroExpandAll(c, f))
	}
	return res
}

// MakeForm attempts to convert a List into a directly applying form
func MakeForm(l List) (List, bool) {
	if f, _, ok := l.Split(); ok {
		if s, ok := f.(Symbol); ok {
			if d := s.Domain(); d != LocalDomain {
				ns := GetNamespace(d)
				if g, ok := ns.Get(s.Name()); ok {
					if ap, ok := g.(Applicable); ok {
						return makeFormObject(l, ap), true
					}
				}
			}
		}
	}
	return nil, false
}

func makeFormObject(l List, a Applicable) List {
	if IsSpecialForm(a) {
		return &specialForm{
			List: l,
			fn:   a,
			args: SequenceToVector(l.Rest()),
		}
	}
	return &evaluatingForm{
		List: l,
		fn:   a,
		args: SequenceToVector(l.Rest()),
	}
}

func (s *specialForm) Eval(c Context) Value {
	return Apply(c, s.fn, s.args)
}

func (f *evaluatingForm) Eval(c Context) Value {
	return Apply(c, f.fn, f.args.Eval(c).(Sequence))
}
