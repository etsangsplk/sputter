package compiler

import "github.com/kode4food/sputter/vm"

// InstructionsProcessor transforms a set of Instructions
type InstructionsProcessor func(vm.Instructions) vm.Instructions
