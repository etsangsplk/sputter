package api

// ExpectedApplicable is thrown when a Value is not Applicable
const ExpectedApplicable = "value does not support application: %s"

// Applicable is the standard signature for any Value that can be applied
// to a sequence of arguments
type Applicable interface {
	Apply(Context, Sequence) Value
}

// Evaluable identifies a Value as being directly evaluable
type Evaluable interface {
	Eval(Context) Value
}

// Block evaluates a Sequence as a Block, returning the last expression
type Block interface {
	Sequence
	IsBlock() bool
}

type block struct {
	Sequence
}

// MacroExpand1 performs a single macro expansion
func MacroExpand1(c Context, v Value) (Value, bool) {
	if l, ok := v.(List); ok && l.IsSequence() {
		f := l.First()
		if s, ok := f.(Symbol); ok {
			if r, ok := s.Resolve(c); ok {
				if f, ok := r.(Function); ok {
					if ok, sp := IsMacro(f); ok && !sp {
						return f.Apply(c, l.Rest()), true
					}
				}
			}
		}
	}
	return v, false
}

// MacroExpand repeatedly performs a macro expansion until no more can occur
func MacroExpand(c Context, v Value) (Value, bool) {
	r := v
	ok := false
	for i := 0; ; i++ {
		if r, ok = MacroExpand1(c, r); ok {
			continue
		}
		return r, i > 0
	}
}

// Eval is a ValueProcessor that expands and evaluates a Value
func Eval(c Context, v Value) Value {
	ex, _ := MacroExpand(c, v)
	if e, ok := ex.(Evaluable); ok {
		return e.Eval(c)
	}
	return ex
}

// EvalBlock expands and evaluates a Sequence as a Block
func EvalBlock(c Context, s Sequence) Value {
	var r Value
	for i := s; i.IsSequence(); i = i.Rest() {
		r = Eval(c, i.First())
	}
	return r
}

// NewBlock casts a Sequence into a Block for evaluation
func NewBlock(s Sequence) Block {
	return &block{Sequence: s}
}

func (b *block) IsBlock() bool {
	return true
}

func (b *block) Eval(c Context) Value {
	return EvalBlock(c, b)
}

// AssertApplicable will cast a Value into an Applicable or explode violently
func AssertApplicable(v Value) Applicable {
	if r, ok := v.(Applicable); ok {
		return r
	}
	panic(Err(ExpectedApplicable, v))
}
