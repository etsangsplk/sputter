package vm

// OpCode represents a decoded VM instruction identifier
type OpCode uint

// These are the OpCodes recognized by the Virtual Machine
const (
	NoOp OpCode = iota
	Load
	Store
	Dup
	Nil
	EmptyList
	True
	False
	Zero
	One
	Const
	PushContext
	PopContext
	Def
	Let
	Eval
	Apply
	First
	Rest
	Prepend
	Truthy
	CondJump
	Jump
	Return
	ReturnNil
	Panic
)
