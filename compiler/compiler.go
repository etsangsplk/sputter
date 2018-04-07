package compiler

import (
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/vm"
)

type (
	// InstructionsProcessor transforms a set of Instructions
	InstructionsProcessor func(vm.Instructions) vm.Instructions

	compilerFunc struct{}
)

var compiler = &compilerFunc{}

// Compile expands and (possibly) compiles a Sequence of Function calls
func Compile(c a.Context, s a.Sequence) a.Sequence {
	return a.Map(c, s, compiler)
}

func (*compilerFunc) Apply(ctx a.Context, args a.Vector) a.Value {
	// Eventually this will do something
	return args[0]
}
