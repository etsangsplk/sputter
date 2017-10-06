package vm

import (
	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
)

var fib *Module

func init() {
	fib = &Module{
		BaseFunction: a.DefaultBaseFunction,
		StackSize:    16,
		LocalsSize:   16,
		Data:         nil,
		Instructions: []Instruction{
			{OpCode: Load, Op1: Args},           //  1: load args
			{OpCode: First},                     //  2: take first arg
			{OpCode: Dup},                       //  3: duplicate it
			{OpCode: Store, Op1: Variables + 0}, //  4: store dup into var 0
			{OpCode: Zero},                      //  5: push 0
			{OpCode: Eq},                        //  6: eq
			{OpCode: CondJump, Op1: 27},         //  7: goto 'return 0'
			{OpCode: Load, Op1: Variables + 0},  //  9: load var 0
			{OpCode: One},                       // 10: push 1
			{OpCode: Eq},                        // 11: eq
			{OpCode: CondJump, Op1: 29},         // 12: goto 'return 1'
			{OpCode: Load, Op1: Variables + 0},  // 13: load var 0
			{OpCode: Const, Op1: 1},             // 14: push 2
			{OpCode: Eq},                        // 15: eq
			{OpCode: CondJump, Op1: 29},         // 16: goto 'return 1'
			{OpCode: Const, Op1: 0},             // 17: push fib
			{OpCode: One},                       // 18: push 1
			{OpCode: Load, Op1: Variables + 0},  // 19: load var 0
			{OpCode: Sub},                       // 20: sub
			{OpCode: Call, Op1: 1},              // 21: call fib
			{OpCode: Const, Op1: 0},             // 22: push fib
			{OpCode: Const, Op1: 1},             // 23: push 2
			{OpCode: Load, Op1: Variables + 0},  // 24: load var 0
			{OpCode: Sub},                       // 25: sub
			{OpCode: Call, Op1: 1},              // 26: call fib
			{OpCode: Add},                       // 27: add
			{OpCode: Return},                    // 28: return sum
			{OpCode: Zero},                      // 29: push 0
			{OpCode: Return},                    // 30: return
			{OpCode: One},                       // 31: push 0
			{OpCode: Return},                    // 32: return
		},
	}

	fib.Data = a.Values{
		fib,
		a.NewFloat(2),
	}

	b.RegisterFunction("fib-vm", fib)
}
