package vm

import (
	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
)

var fib *Module

func init() {
	inst := RemoveLabels([]Instruction{
		{OpCode: Load, Op1: Args},           // load args
		{OpCode: First},                     // take first arg
		{OpCode: Dup},                       // icate it
		{OpCode: Store, Op1: Variables + 0}, // store dup into var 0
		{OpCode: Zero},                      // push 0
		{OpCode: Eq},                        // eq
		{OpCode: CondJump, Op1: 0},          // goto label 0
		{OpCode: Load, Op1: Variables + 0},  // load var 0
		{OpCode: One},                       // push 1
		{OpCode: Eq},                        // eq
		{OpCode: CondJump, Op1: 1},          // goto label 1
		{OpCode: Load, Op1: Variables + 0},  // load var 0
		{OpCode: Const, Op1: 1},             // push 2
		{OpCode: Eq},                        // eq
		{OpCode: CondJump, Op1: 1},          // goto label 1
		{OpCode: Const, Op1: 0},             // push fib
		{OpCode: One},                       // push 1
		{OpCode: Load, Op1: Variables + 0},  // load var 0
		{OpCode: Sub},                       // sub
		{OpCode: Call, Op1: 1},              // call fib
		{OpCode: Const, Op1: 0},             // push fib
		{OpCode: Const, Op1: 1},             // push 2
		{OpCode: Load, Op1: Variables + 0},  // load var 0
		{OpCode: Sub},                       // sub
		{OpCode: Call, Op1: 1},              // call fib
		{OpCode: Add},                       // add
		{OpCode: Return},                    // return sum
		{OpCode: Label, Op1: 0},             // label 0
		{OpCode: Zero},                      // push 0
		{OpCode: Return},                    // return
		{OpCode: Label, Op1: 1},             // label 1
		{OpCode: One},                       // push 1
		{OpCode: Return},                    // return
	})

	fib = &Module{
		BaseFunction: a.DefaultBaseFunction,
		StackSize:    16,
		LocalsSize:   16,
		Data:         nil,
		Instructions: inst,
	}

	fib.Data = a.Values{
		fib,
		a.NewFloat(2),
	}

	b.RegisterFunction("fib-vm", fib)
}
