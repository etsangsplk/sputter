package api

// List contains a node to a singly-linked List
type List interface {
	Conjoiner
	MakeEvaluable
	Elementer
	Counted
	Applicable
	List() bool
}

type list struct {
	first Value
	rest  *list
	count int
}

type evaluableList struct {
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

// List is a disambiguating marker
func (l *list) List() bool {
	return true
}

// First returns the first element of a List
func (l *list) First() Value {
	return l.first
}

// Rest returns the rest of the List as a Sequence
func (l *list) Rest() Sequence {
	return l.rest
}

// IsSequence returns whether this instance is a consumable Sequence
func (l *list) IsSequence() bool {
	return l != EmptyList
}

// Prepend creates a new Sequence by prepending a Value
func (l *list) Prepend(v Value) Sequence {
	return &list{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

// Conjoin implements the Conjoiner interface
func (l *list) Conjoin(v Value) Sequence {
	return l.Prepend(v)
}

// Count returns the length of the List
func (l *list) Count() int {
	return l.count
}

// ElementAt returns the Value at the indexed position in the List
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

// Apply makes List applicable
func (l *list) Apply(c Context, args Sequence) Value {
	return IndexedApply(l, c, args)
}

// Evaluable turns List into an Evaluable Expression
func (l *list) Evaluable() Value {
	if l == EmptyList {
		return l
	}
	return &evaluableList{
		list: l,
	}
}

// Str converts this Value into a Str
func (l *list) Str() Str {
	return MakeSequenceStr(l)
}

// Eval makes evaluableList Evaluable
func (e *evaluableList) Eval(c Context) Value {
	t := e.first
	if a, ok := Eval(c, t).(Applicable); ok {
		return Apply(c, a, e.rest)
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
