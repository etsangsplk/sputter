package builtins

import (
	a "github.com/kode4food/sputter/api"
	c "github.com/kode4food/sputter/compiler"
	"github.com/kode4food/sputter/vm"
)

var fib *vm.Module

func init() {
	inst := c.RemoveLabels([]vm.Instruction{
		// (defn fib [i] ...)
		vm.MakeInst(vm.First, vm.Args, vm.Locals+0),
		// (cond (= i 0) 0)
		vm.MakeInst(vm.Zero, vm.Locals+1),
		vm.MakeInst(vm.Eq, vm.Locals+0, vm.Locals+1, vm.Locals+1),
		vm.MakeInst(vm.CondJumpLabel, 0, vm.Locals+1), // goto return 0
		// (cond (= i 1) 1)
		vm.MakeInst(vm.One, vm.Locals+1),
		vm.MakeInst(vm.Eq, vm.Locals+0, vm.Locals+1, vm.Locals+1),
		vm.MakeInst(vm.CondJumpLabel, 1, vm.Locals+1), // goto return 1
		// (cond (= i 2) 1)
		vm.MakeInst(vm.Const, 1, vm.Locals+2), // 2
		vm.MakeInst(vm.Eq, vm.Locals+0, vm.Locals+2, vm.Locals+1),
		vm.MakeInst(vm.CondJumpLabel, 1, vm.Locals+1), // goto return 1
		// (cond :else)
		vm.MakeInst(vm.Const, 0, vm.Locals+3), // fib
		// (fib (- i 1))
		vm.MakeInst(vm.One, vm.Locals+1),
		vm.MakeInst(vm.Sub, vm.Locals+0, vm.Locals+1, vm.Locals+1),
		vm.MakeInst(vm.Vector, vm.Locals+1, vm.Locals+2, vm.Locals+1),
		vm.MakeInst(vm.Apply, vm.Locals+1, vm.Locals+3, vm.Locals+1),
		// (fib (- i 2))
		vm.MakeInst(vm.Sub, vm.Locals+0, vm.Locals+2, vm.Locals+2),
		vm.MakeInst(vm.Vector, vm.Locals+2, vm.Locals+3, vm.Locals+2),
		vm.MakeInst(vm.Apply, vm.Locals+2, vm.Locals+3, vm.Locals+2),
		// (+ (fib (- i 1)) (fib (- i 2)))
		vm.MakeInst(vm.Add, vm.Locals+1, vm.Locals+2, vm.Locals+2),
		vm.MakeInst(vm.Return, vm.Locals+2),
		// return 0
		vm.MakeInst(vm.Label, 0),
		vm.MakeInst(vm.Zero, vm.Locals+0),
		vm.MakeInst(vm.Return, vm.Locals+0),
		// return 1
		vm.MakeInst(vm.Label, 1),
		vm.MakeInst(vm.One, vm.Locals+0),
		vm.MakeInst(vm.Return, vm.Locals+0),
	})

	fib = &vm.Module{
		BaseFunction: a.DefaultBaseFunction,
		LocalsSize:   8,
		Data:         nil,
		Instructions: inst,
	}

	fib.Data = a.Values{
		fib,
		a.NewFloat(2),
	}

	RegisterFunction("fib-vm", fib)
}
