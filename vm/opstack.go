package vm

import a "github.com/kode4food/sputter/api"

type operandStack []a.Value

var emptyOperandStack = make(operandStack, 0)

func (o operandStack) IsSequence() bool {
	return len(o) > 0
}

func (o operandStack) First() a.Value {
	if len(o) > 0 {
		return o[0]
	}
	return a.Nil
}

func (o operandStack) Rest() a.Sequence {
	if len(o) > 1 {
		return o[1:]
	}
	return emptyOperandStack
}

func (o operandStack) Split() (a.Value, a.Sequence, bool) {
	lv := len(o)
	if lv > 1 {
		return o[0], o[1:], true
	} else if lv == 1 {
		return o[0], emptyOperandStack, true
	}
	return a.Nil, emptyOperandStack, false
}

func (o operandStack) Prepend(p a.Value) a.Sequence {
	return append(operandStack{p}, o...)
}

func (o operandStack) Str() a.Str {
	return a.MakeSequenceStr(o)
}
