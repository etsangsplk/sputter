package vm

const (
	NoOp byte = iota
	PushFrame
	PopFrame
	PushContext
	PopContext
	Nil
	Return
	Halt
)
