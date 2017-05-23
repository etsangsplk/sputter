package evaluator

import a "github.com/kode4food/sputter/api"

// ExpectedSequence is raised when unquote-splicing expands a non sequence
const ExpectedSequence = "can not splice a non-sequence: %s"

type expander struct {
	context a.Context
}

type splice []a.Value

// Expand will do just that, expand macros within a sequence
func Expand(c a.Context, v a.Value) a.Value {
	e := &expander{
		context: c,
	}

	return e.expandValue(v)
}

// ExpandSequence will expand a sequence into a new sequence
func ExpandSequence(c a.Context, s a.Sequence) a.Sequence {
	e := &expander{
		context: c,
	}

	return a.Map(s, func(v a.Value) a.Value {
		return e.expandValue(v)
	})
}

func (e *expander) expandValue(v a.Value) a.Value {
	if s, ok := v.(a.Sequence); ok {
		return e.expandSequence(s)
	}
	return e.makeExpression(v)
}

func (e *expander) makeExpression(v a.Value) a.Value {
	if m, ok := v.(a.MakeExpression); ok {
		return m.Expression()
	}
	return v
}

func (e *expander) expandSequence(s a.Sequence) a.Value {
	if l, ok := s.(a.List); ok {
		return e.expandList(l)
	}
	if v, ok := s.(a.Vector); ok {
		return e.expandVector(v)
	}
	if as, ok := s.(a.Associative); ok {
		return e.expandAssociative(as)
	}
	return s
}

func (e *expander) expandList(l a.List) a.Value {
	f := l.First()
	if m, ok := e.macro(f); ok {
		res := m.Apply(e.context, l.Rest())
		if isSplicing(m) {
			if s, ok := res.(a.Sequence); ok {
				return makeSplice(s)
			}
			panic(a.Err(ExpectedSequence, res))
		}
		return res
	}
	return a.NewList(e.expandElements(l)...).Expression()
}

func (e *expander) expandVector(v a.Vector) a.Value {
	return a.NewVector(e.expandElements(v)...)
}

func (e *expander) expandAssociative(as a.Associative) a.Value {
	r := e.expandElements(as)
	v := make([]a.Vector, len(r))
	for i, e := range r {
		v[i] = e.(a.Vector)
	}
	return a.NewAssociative(v...)
}

func makeSplice(s a.Sequence) splice {
	sp := splice{}
	for i := s; i.IsSequence(); i = i.Rest() {
		sp = append(sp, i.First())
	}
	return sp
}

func (e *expander) expandElements(s a.Sequence) []a.Value {
	r := []a.Value{}
	for i := s; i.IsSequence(); i = i.Rest() {
		e := e.expandValue(i.First())
		if sp, ok := e.(splice); ok {
			r = append(r, sp...)
		} else {
			r = append(r, e)
		}
	}
	return r
}

func (e *expander) macro(v a.Value) (a.Function, bool) {
	if s, ok := v.(a.Symbol); ok {
		if r, ok := s.Resolve(e.context); ok {
			if f, ok := r.(a.Function); ok {
				if a.IsMacro(f) {
					return f, true
				}
			}
		}
	}
	return nil, false
}

func isSplicing(m a.Function) bool {
	if v, ok := m.Metadata().Get(a.MetaSplicing); ok {
		return v == a.True
	}
	return false
}

func (s splice) Eval(_ a.Context) a.Value {
	return s
}

func (s splice) Str() a.Str {
	return a.MakeSequenceStr(a.NewVector(s...))
}
