package api

import "bytes"

type (
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
		Values
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
	return &block{
		Values: SequenceToValues(s),
	}
}

func (*block) BlockType() {}

func (b *block) Eval(c Context) Value {
	var res Value = Nil
	for _, f := range b.Values {
		res = Eval(c, f)
	}
	return res
}

func (b *block) Str() Str {
	var buf bytes.Buffer
	for _, f := range b.Values {
		buf.WriteString(string(f.Str()))
	}
	return Str(buf.String())
}
