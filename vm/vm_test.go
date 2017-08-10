package vm_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	"github.com/kode4food/sputter/vm"
)

var vmTestData = []a.Value{
	a.Str("The First Data"),
	a.Str("Hello There"),
}

func TestReturn(t *testing.T) {
	as := assert.New(t)
	m := &vm.Module{
		LocalsSize: 1,
		StackSize:  1,
		Data:       vmTestData,
		Instructions: []vm.Instruction{
			{OpCode: vm.Const, Args: [3]uint{1, 0, 0}},
			{OpCode: vm.Return},
		},
	}
	r := m.Apply(a.NewContext(), a.EmptyList)
	as.String("Hello There", r)
}
