package vm

// OpCode represents a decoded VM instruction identifier
type OpCode uint

// These are the OpCodes recognized by the Virtual Machine
const (
	NoOp OpCode = iota
	Load
	Store
	Clear
	Dup
	Swap
	Nil
	EmptyList
	True
	False
	Zero
	One
	Const
	Def
	Let
	Eval
	Apply
	IsSeq
	First
	Rest
	Split
	Prepend
	CondJump
	Jump
	Return
	ReturnNil
	Panic
)
