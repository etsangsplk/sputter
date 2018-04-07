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
		vm.MakeInst(vm.GetArg, 0, vm.Vars+0),
		// (cond (= i 0) 0)
		vm.MakeInst(vm.Zero, vm.Vars+1),
		vm.MakeInst(vm.Eq, vm.Vars+0, vm.Vars+1, vm.Vars+1),
		vm.MakeInst(vm.CondJumpLabel, 0, vm.Vars+1), // goto return 0
		// (cond (= i 1) 1)
		vm.MakeInst(vm.One, vm.Vars+1),
		vm.MakeInst(vm.Eq, vm.Vars+0, vm.Vars+1, vm.Vars+1),
		vm.MakeInst(vm.CondJumpLabel, 1, vm.Vars+1), // goto return 1
		// (cond (= i 2) 1)
		vm.MakeInst(vm.Const, 1, vm.Vars+2), // 2
		vm.MakeInst(vm.Eq, vm.Vars+0, vm.Vars+2, vm.Vars+1),
		vm.MakeInst(vm.CondJumpLabel, 1, vm.Vars+1), // goto return 1
		// (cond :else)
		vm.MakeInst(vm.Const, 0, vm.Vars+3), // fib
		// (fib (- i 1))
		vm.MakeInst(vm.One, vm.Vars+1),
		vm.MakeInst(vm.Sub, vm.Vars+0, vm.Vars+1, vm.Vars+1),
		vm.MakeInst(vm.Vector, vm.Vars+1, vm.Vars+2, vm.Vars+1),
		vm.MakeInst(vm.Apply, vm.Vars+1, vm.Vars+3, vm.Vars+1),
		// (fib (- i 2))
		vm.MakeInst(vm.Sub, vm.Vars+0, vm.Vars+2, vm.Vars+2),
		vm.MakeInst(vm.Vector, vm.Vars+2, vm.Vars+3, vm.Vars+2),
		vm.MakeInst(vm.Apply, vm.Vars+2, vm.Vars+3, vm.Vars+2),
		// (+ (fib (- i 1)) (fib (- i 2)))
		vm.MakeInst(vm.Add, vm.Vars+1, vm.Vars+2, vm.Vars+2),
		vm.MakeInst(vm.Return, vm.Vars+2),
		// return 0
		vm.MakeInst(vm.Label, 0),
		vm.MakeInst(vm.Zero, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
		// return 1
		vm.MakeInst(vm.Label, 1),
		vm.MakeInst(vm.One, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	})

	fib = &vm.Module{
		BaseFunction: a.DefaultBaseFunction,
		Data:         nil,
		Instructions: inst,
	}

	fib.Data = a.Vector{
		fib,
		a.NewFloat(2),
	}

	RegisterFunction("fib-vm", fib)
}
