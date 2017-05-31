package evaluator

import a "github.com/kode4food/sputter/api"

type Closure interface {
	a.Value
	Names() []a.Name
}

type closure struct {
	names []a.Name
	value a.Value
}

var emptyNames = []a.Name{}

// NewClosure may produce a new value whose evaluation context is decoupled
func NewClosure(v a.Value) a.Value {
	n := visitValue(v)
	if len(n) > 0 {
		return &closure{
			names: visitValue(v),
			value: v,
		}
	}
	return v
}

func (cl *closure) Names() []a.Name {
	return cl.names
}

func (cl *closure) Eval(c a.Context) a.Value {
	names := cl.names
	vars := make(a.Variables, len(names))
	for _, n := range names {
		if v, ok := c.Get(n); ok {
			vars[n] = v
		}
	}

	ns := a.GetContextNamespace(c)
	l := a.ChildContextVars(ns, vars)
	return cl.value.Eval(l)
}

func (cl *closure) Str() a.Str {
	return a.MakeDumpStr(cl)
}

func visitValue(v a.Value) []a.Name {
	if s, ok := v.(a.Sequence); ok {
		return visitSequence(s)
	}
	if _, ok := v.(a.Expression); !ok {
		return emptyNames
	}
	if s, ok := v.(a.Symbol); ok && s.Domain() == a.LocalDomain {
		return []a.Name{s.Name()}
	}
	return emptyNames
}

func visitSequence(s a.Sequence) []a.Name {
	if _, ok := s.(a.Str); ok {
		return emptyNames
	}
	r := []a.Name{}
	for i := s; i.IsSequence(); i = i.Rest() {
		n := visitValue(i.First())
		r = append(r, n...)
	}
	return r
}
