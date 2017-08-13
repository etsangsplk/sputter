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
	if l, ok := v.(List); ok && l.IsSequence() {
		f := l.First()
		if s, ok := f.(Symbol); ok {
			if r, ok := s.Resolve(c); ok {
				if a, ok := r.(Applicable); ok {
					if ok, sp := IsMacro(a); ok && !sp {
						return a.Apply(c, l.Rest()), true
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
	r := []Vector{}
	var t Value
	for i := as.(Sequence); i.IsSequence(); {
		t, i = i.Split()
		e := t.(Vector)
		k, _ := e.ElementAt(0)
		v, _ := e.ElementAt(1)
		r = append(r, NewVector(
			MacroExpandAll(c, k),
			MacroExpandAll(c, v),
		))
	}
	return NewAssociative(r...)
}

func expandElements(c Context, s Sequence) []Value {
	r := []Value{}
	var t Value
	for i := s; i.IsSequence(); {
		t, i = i.Split()
		r = append(r, MacroExpandAll(c, t))
	}
	return r
}

// MakeForm attempts to convert a List into a directly applying form
func MakeForm(l List) (List, bool) {
	if !l.IsSequence() {
		return l, false
	}
	if s, ok := l.First().(Symbol); ok {
		if d := s.Domain(); d != LocalDomain {
			ns := GetNamespace(d)
			if g, ok := ns.Get(s.Name()); ok {
				if ap, ok := g.(Applicable); ok {
					return makeFormObject(l, ap), true
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
	return s.fn.Apply(c, s.args)
}

func (f *evaluatingForm) Eval(c Context) Value {
	return f.fn.Apply(c, f.args.Eval(c).(Vector))
}
