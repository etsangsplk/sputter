package api

import "bytes"

type (
	// Applicable is the standard signature for any Value that can be applied
	// to a sequence of arguments
	Applicable interface {
		Apply(Context, Sequence) Value
	}

	// Evaluable identifies a Value as being directly evaluable
	Evaluable interface {
		Eval(Context) Value
	}

	// Block evaluates a Sequence as a Block, returning the last expression
	Block interface {
		Sequence
		Evaluable
		BlockType()
	}

	block struct {
		Sequence
	}
)

// Eval is a ValueProcessor that expands and evaluates a Value
func Eval(c Context, v Value) Value {
	ex, _ := MacroExpand(c, v)
	if e, ok := ex.(Evaluable); ok {
		return e.Eval(c)
	}
	return ex
}

// MakeBlock casts a Sequence into a Block for evaluation
func MakeBlock(s Sequence) Block {
	if b, ok := s.(Block); ok {
		return b
	}
	return &block{Sequence: s}
}

func (b *block) BlockType() {}

func (b *block) Eval(c Context) Value {
	var res Value = Nil
	for f, r, ok := b.Sequence.Split(); ok; f, r, ok = r.Split() {
		res = Eval(c, f)
	}
	return res
}

func (b *block) Str() Str {
	var buf bytes.Buffer
	for f, r, ok := b.Sequence.Split(); ok; f, r, ok = r.Split() {
		buf.WriteString(string(f.Str()))
	}
	return Str(buf.String())
}
