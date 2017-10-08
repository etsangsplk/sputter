package builtins

import (
	a "github.com/kode4food/sputter/api"
	c "github.com/kode4food/sputter/compiler"
	"github.com/kode4food/sputter/vm"
)

var fib *vm.Module

func init() {
	inst := c.RemoveLabels([]vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args},           // load args
		{OpCode: vm.First},                        // take first arg
		{OpCode: vm.Dup},                          // duplicate it
		{OpCode: vm.Store, Op1: vm.Variables + 0}, // store dup into var 0
		{OpCode: vm.Zero},                         // push 0
		{OpCode: vm.Eq},                           // eq
		{OpCode: vm.CondJump, Op1: 0},             // goto label 0
		{OpCode: vm.Load, Op1: vm.Variables + 0},  // load var 0
		{OpCode: vm.One},                          // push 1
		{OpCode: vm.Eq},                           // eq
		{OpCode: vm.CondJump, Op1: 1},             // goto label 1
		{OpCode: vm.Load, Op1: vm.Variables + 0},  // load var 0
		{OpCode: vm.Const, Op1: 1},                // push 2
		{OpCode: vm.Eq},                           // eq
		{OpCode: vm.CondJump, Op1: 1},             // goto label 1
		{OpCode: vm.Const, Op1: 0},                // push fib
		{OpCode: vm.One},                          // push 1
		{OpCode: vm.Load, Op1: vm.Variables + 0},  // load var 0
		{OpCode: vm.Sub},                          // sub
		{OpCode: vm.Call, Op1: 1},                 // call fib
		{OpCode: vm.Const, Op1: 0},                // push fib
		{OpCode: vm.Const, Op1: 1},                // push 2
		{OpCode: vm.Load, Op1: vm.Variables + 0},  // load var 0
		{OpCode: vm.Sub},                          // sub
		{OpCode: vm.Call, Op1: 1},                 // call fib
		{OpCode: vm.Add},                          // add
		{OpCode: vm.Return},                       // return sum
		{OpCode: vm.Label, Op1: 0},                // label 0
		{OpCode: vm.Zero},                         // push 0
		{OpCode: vm.Return},                       // return
		{OpCode: vm.Label, Op1: 1},                // label 1
		{OpCode: vm.One},                          // push 1
		{OpCode: vm.Return},                       // return
	})

	fib = &vm.Module{
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

	RegisterFunction("fib-vm", fib)
}
