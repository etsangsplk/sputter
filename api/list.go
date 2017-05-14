package api

// List contains a node to a singly-linked List
type List struct {
	first Value
	rest  *List
	count int
}

// EvaluableList represents an Evaluable List
type EvaluableList struct {
	*List
}

// EmptyList represents an empty List
var EmptyList *List

// NewList creates a new List instance
func NewList(v Value) *List {
	return &List{
		first: v,
		rest:  EmptyList,
		count: 1,
	}
}

// First returns the first element of a List
func (l *List) First() Value {
	return l.first
}

// Rest returns the rest of the List as a Sequence
func (l *List) Rest() Sequence {
	return l.rest
}

// IsSequence returns whether this instance is a consumable Sequence
func (l *List) IsSequence() bool {
	return l != EmptyList
}

// Prepend creates a new Sequence by prepending a Value
func (l *List) Prepend(v Value) Sequence {
	return &List{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

// Conjoin implements the Conjoiner interface
func (l *List) Conjoin(v Value) Sequence {
	return l.Prepend(v)
}

// Count returns the length of the List
func (l *List) Count() int {
	return l.count
}

// ElementAt returns the Value at the indexed position in the List
func (l *List) ElementAt(index int) (Value, bool) {
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
func (l *List) Apply(c Context, args Sequence) Value {
	return IndexedApply(l, c, args)
}

// Evaluable turns List into an Evaluable Expression
func (l *List) Evaluable() Value {
	if l == EmptyList {
		return l
	}
	return &EvaluableList{
		List: l,
	}
}

// Str converts this Value into a Str
func (l *List) Str() Str {
	return MakeSequenceStr(l)
}

// Eval makes EvaluableList Evaluable
func (e *EvaluableList) Eval(c Context) Value {
	l := e.List
	t := l.first
	if a, ok := Eval(c, t).(Applicable); ok {
		return Apply(c, a, l.rest)
	}
	panic(Err(ExpectedApplicable, t))
}

func init() {
	EmptyList = &List{
		first: Nil,
		count: 0,
	}
	EmptyList.rest = EmptyList
}
