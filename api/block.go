package api

import "bytes"

type (
	// Block evaluates a Sequence as a Block, returning the last expression
	Block interface {
		Sequence
		Evaluable
		BlockType()
	}

	block struct {
		Vector
	}

	emptyBlock struct {
		block
	}

	singleBlock struct {
		block
		value Value
	}

	multiBlock struct {
		block
		first Vector
		last  Value
	}
)

// EvalVectorAsBlock evaluates Vector as if they were a Block
func EvalVectorAsBlock(c Context, v Vector) Value {
	l := len(v)
	switch l {
	case 0:
		return Nil
	case 1:
		return Eval(c, v[0])
	default:
		fl := l - 1
		for _, e := range v[0:fl] {
			Eval(c, e)
		}
		return Eval(c, v[fl])
	}
}

// MakeBlock casts a Sequence into a Block for evaluation
func MakeBlock(s Sequence) Block {
	if b, ok := s.(Block); ok {
		return b
	}
	v := SequenceToVector(s)
	switch len(v) {
	case 0:
		return &emptyBlock{
			block: block{Vector: v},
		}
	case 1:
		return &singleBlock{
			block: block{Vector: v},
			value: v[0],
		}
	default:
		fl := len(v) - 1
		return &multiBlock{
			block: block{Vector: v},
			first: v[0:fl],
			last:  v[fl],
		}
	}
}

func (*emptyBlock) BlockType()  {}
func (*singleBlock) BlockType() {}
func (*multiBlock) BlockType()  {}

func (*emptyBlock) Eval(_ Context) Value {
	return Nil
}

func (b *singleBlock) Eval(c Context) Value {
	return Eval(c, b.value)
}

func (b *multiBlock) Eval(c Context) Value {
	for _, f := range b.first {
		Eval(c, f)
	}
	return Eval(c, b.last)
}

func (b *block) Str() Str {
	var buf bytes.Buffer
	for _, f := range b.Vector {
		buf.WriteString(string(f.Str()))
	}
	return Str(buf.String())
}
