package api

// List contains a node to a singly-linked List
type List interface {
	Conjoiner
	MakeExpression
	Elementer
	Counted
	Applicable
	IsList() bool
}

type list struct {
	first Value
	rest  *list
	count int
}

type listExpression struct {
	*list
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

func (l *list) Apply(c Context, args Sequence) Value {
	return IndexedApply(l, c, args)
}

func (l *list) Eval(_ Context) Value {
	return l
}

func (l *list) Expression() Value {
	if l == EmptyList {
		return l
	}
	return &listExpression{
		list: l,
	}
}

func (l *list) Str() Str {
	return MakeSequenceStr(l)
}

func (l *listExpression) IsExpression() bool {
	return true
}

func (l *listExpression) Eval(c Context) Value {
	t := l.first
	if a, ok := t.Eval(c).(Applicable); ok {
		return a.Apply(c, l.rest)
	}
	panic(Err(ExpectedApplicable, t))
}

func init() {
	emptyList = &list{
		first: Nil,
		count: 0,
	}
	emptyList.rest = emptyList
	EmptyList = emptyList
}
