package vm_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	"github.com/kode4food/sputter/vm"
)

var (
	vmTestData = []a.Value{
		a.Str("The first bit of data"),
		a.Str("Hello there!"),
	}

	vmTestArgs = a.NewList(a.Str("first"), a.Str("second"), a.Str("third"))
)

func moduleFromInstructions(inst []vm.Instruction) *vm.Module {
	return &vm.Module{
		LocalsSize:   8,
		StackSize:    8,
		Data:         vmTestData,
		Instructions: inst,
	}
}

func TestReturn(t *testing.T) {
	as := assert.New(t)
	m := moduleFromInstructions([]vm.Instruction{
		{OpCode: vm.Const, Args: [3]uint{1, 0, 0}},
		{OpCode: vm.Return},
	})
	r := m.Apply(a.NewContext(), a.EmptyList)
	as.String("Hello there!", r)
}

func TestFirst(t *testing.T) {
	as := assert.New(t)
	m := moduleFromInstructions([]vm.Instruction{
		{OpCode: vm.Load, Args: [3]uint{0, 0, 0}},
		{OpCode: vm.First},
		{OpCode: vm.Return},
	})
	r := m.Apply(a.NewContext(), vmTestArgs)
	as.String("first", r)
}

func TestRest(t *testing.T) {
	as := assert.New(t)
	m := moduleFromInstructions([]vm.Instruction{
		{OpCode: vm.Load, Args: [3]uint{0, 0, 0}},
		{OpCode: vm.Rest},
		{OpCode: vm.Return},
	})
	r := m.Apply(a.NewContext(), vmTestArgs)
	as.String(`("second" "third")`, r)
}

func TestPrepend(t *testing.T) {
	as := assert.New(t)
	m := moduleFromInstructions([]vm.Instruction{
		{OpCode: vm.Const, Args: [3]uint{1, 0, 0}},
		{OpCode: vm.Load, Args: [3]uint{0, 0, 0}},
		{OpCode: vm.Prepend},
		{OpCode: vm.Return},
	})
	r := m.Apply(a.NewContext(), vmTestArgs)
	as.String(`("Hello there!" "first" "second" "third")`, r)
}
