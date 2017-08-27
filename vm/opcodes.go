package vm

// OpCode represents a decoded VM instruction identifier
type OpCode uint

// These are the OpCodes recognized by the Virtual Machine
const (
	NoOp OpCode = iota
	Pop
	Load
	Store
	StoreConst
	Clear
	Dup
	Swap
	Nil
	EmptyList
	True
	False
	Zero
	One
	NegOne
	Const
	Def
	Let
	Eval
	Apply
	Call
	Vector
	IsSeq
	First
	Rest
	Split
	Prepend
	Add
	Sub
	Mul
	Div
	Mod
	Eq
	Neq
	CondJump
	Jump
	Return
	Panic
)

// These are the default positions for certain locals
const (
	Context uint = iota
	Args
	Variables
)
