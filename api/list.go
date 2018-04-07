package api

// List contains a node to a singly-linked List
type List struct {
	first Value
	rest  *List
	count int
}

// EmptyList represents an empty List
var EmptyList *List

// NewList creates a new List instance
func NewList(v ...Value) *List {
	r := EmptyList
	for i := len(v) - 1; i >= 0; i-- {
		r = &List{
			first: v[i],
			rest:  r,
			count: r.count + 1,
		}
	}
	return r
}

// First returns the first element of the List
func (l *List) First() Value {
	return l.first
}

// Rest returns the elements of the List that follow the first
func (l *List) Rest() Sequence {
	return l.rest
}

// IsSequence returns whether or not this List has any elements
func (l *List) IsSequence() bool {
	return l != EmptyList
}

// Split breaks the List into its components (first, rest, isSequence)
func (l *List) Split() (Value, Sequence, bool) {
	return l.first, l.rest, l != EmptyList
}

// Prepend inserts an element at the beginning of the List
func (l *List) Prepend(v Value) Sequence {
	return &List{
		first: v,
		rest:  l,
		count: l.count + 1,
	}
}

// Conjoin appends an element to the beginning of the List
func (l *List) Conjoin(v Value) Sequence {
	return l.Prepend(v)
}

// Count returns the number of elements in the List
func (l *List) Count() int {
	return l.count
}

// ElementAt returns a specific element of the List
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

// Apply makes List Applicable
func (l *List) Apply(_ Context, args Vector) Value {
	return IndexedApply(l, args)
}

// Eval evaluates its elements, returning a new List
func (l *List) Eval(c Context) Value {
	if l == EmptyList {
		return EmptyList
	}

	t := l.first
	a := Eval(c, t).(Applicable)
	if IsSpecialForm(a) {
		return Apply(c, a, l.rest)
	}
	return Apply(c, a, l.evalArgs(c, l.rest))
}

// Str converts this List to a Str
func (l *List) Str() Str {
	return MakeSequenceStr(l)
}

func (l *List) evalArgs(c Context, args *List) Vector {
	ac := args.count
	r := make(Vector, ac)
	for idx, i := 0, args; idx < ac; idx++ {
		r[idx] = Eval(c, i.first)
		i = i.rest
	}
	return r
}

func init() {
	EmptyList = &List{
		first: Nil,
		count: 0,
	}
	EmptyList.rest = EmptyList
}
