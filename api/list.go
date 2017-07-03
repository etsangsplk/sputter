package api

// List contains a node to a singly-linked List
type List interface {
	Conjoiner
	Indexed
	Counted
	Applicable
	Evaluable
	IsList() bool
}

type list struct {
	first Value
	rest  *list
	count int
}

// EmptyList represents an empty List
var EmptyList List
var emptyList *list

// NewList creates a new List instance
func NewList(v ...Value) List {
	r := emptyList
	for i := len(v) - 1; i >= 0; i-- {
		r = &list{
			first: v[i],
			rest:  r,
			count: r.count + 1,
		}
	}
	return r
}

func (l *list) IsList() bool {
	return true
}

func (l *list) First() Value {
	return l.first
}

func (l *list) Rest() Sequence {
	return l.rest
}

func (l *list) IsSequence() bool {
	return l != EmptyList
}

func (l *list) Prepend(v Value) Sequence {
	return &list{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

func (l *list) Conjoin(v Value) Sequence {
	return l.Prepend(v)
}

func (l *list) Count() int {
	return l.count
}

func (l *list) ElementAt(index int) (Value, bool) {
	if index > l.count-1 || index < 0 {
		return Nil, false
	}

	e := l
	for i := 0; i < index; i++ {
		e = e.rest
	}
	return e.first, true
}

func (l *list) Apply(_ Context, args Sequence) Value {
	return IndexedApply(l, args)
}

func (l *list) Eval(c Context) Value {
	if l == emptyList {
		return emptyList
	}

	t := l.first
	if a, ok := Eval(c, t).(Applicable); ok {
		if IsSpecialForm(a) {
			return a.Apply(c, l.rest)
		}
		return a.Apply(c, l.evalArgs(c, l.rest))
	}
	panic(Err(ExpectedApplicable, t))
}

func (l *list) evalArgs(c Context, args *list) Vector {
	ac := args.count
	r := make(vector, ac)
	for idx, i := 0, args; idx < ac; idx++ {
		r[idx] = Eval(c, i.first)
		i = i.rest
	}
	return r
}

func (l *list) Str() Str {
	return MakeSequenceStr(l)
}

func init() {
	emptyList = &list{
		first: Nil,
		count: 0,
	}
	emptyList.rest = emptyList
	EmptyList = emptyList
}
