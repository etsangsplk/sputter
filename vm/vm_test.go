package vm_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	"github.com/kode4food/sputter/vm"
)

var (
	vmTestData = []a.Value{
		s("The first bit of data"),
		s("Hello there!"),
		a.ErrStr("i blew up!"),
		a.True,
		a.False,
	}

	vmTestArgs = a.NewList(s("first"), s("second"), s("third"))
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}

func testInstructions(t *testing.T, inst []vm.Instruction, expect a.Value) {
	as := assert.New(t)

	m1 := &vm.Module{
		BaseFunction: a.DefaultBaseFunction,
		LocalsSize:   16,
		StackSize:    32,
		Data:         vmTestData,
		Instructions: inst,
	}

	r1 := m1.Apply(a.NewContext(), vmTestArgs)
	as.Equal(expect, r1)

	m2 := m1.WithMetadata(a.Properties{
		a.NameKey: a.Name("newMetadata"),
	}).(*vm.Module)
	r2 := m2.Apply(a.NewContext(), vmTestArgs)
	as.Equal(expect, r2)
}

func testMapped(t *testing.T, o vm.OpCode, expect a.Value) {
	testInstructions(t, []vm.Instruction{
		{OpCode: o},
		{OpCode: vm.Return},
	}, expect)
}

func TestMappedOpCodes(t *testing.T) {
	testMapped(t, vm.Nil, a.Nil)
	testMapped(t, vm.EmptyList, a.EmptyList)
	testMapped(t, vm.True, a.True)
	testMapped(t, vm.False, a.False)
	testMapped(t, vm.Zero, f(0))
	testMapped(t, vm.One, f(1))
}

func TestSwap(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Arg0: 0},
		{OpCode: vm.Const, Arg0: 1},
		{OpCode: vm.Swap},
		{OpCode: vm.Return},
	}, s("The first bit of data"))
}

func TestReturnNil(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.ReturnNil},
	}, a.Nil)
}

func TestReturn(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Arg0: 1},
		{OpCode: vm.Return},
	}, s("Hello there!"))
}

func TestFirst(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Arg0: 0},
		{OpCode: vm.First},
		{OpCode: vm.Return},
	}, s("first"))
}

func TestRest(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Arg0: 0},
		{OpCode: vm.Rest},
		{OpCode: vm.Return},
	}, s(`("second" "third")`))
}

func TestPrepend(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Arg0: 0},
		{OpCode: vm.Const, Arg0: 1},
		{OpCode: vm.Prepend},
		{OpCode: vm.Return},
	}, s(`("Hello there!" "first" "second" "third")`))
}

func TestDup(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Arg0: 0},
		{OpCode: vm.Dup},
		{OpCode: vm.First},
		{OpCode: vm.Dup},
		{OpCode: vm.Store, Arg0: 1},
		{OpCode: vm.Prepend},
		{OpCode: vm.Load, Arg0: 1},
		{OpCode: vm.Prepend},
		{OpCode: vm.Return},
	}, s(`("first" "first" "first" "second" "third")`))
}

func TestClear(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Clear, Arg0: 0},
		{OpCode: vm.Load, Arg0: 0},
		{OpCode: vm.Return},
	}, a.Nil)
}

func TestJump(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Arg0: 1},
		{OpCode: vm.Jump, Arg0: 3},
		{OpCode: vm.Const, Arg0: 0},
		{OpCode: vm.Return},
	}, s("Hello there!"))
}

func TestCondJump(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Arg0: 3},
		{OpCode: vm.NoOp},
		{OpCode: vm.CondJump, Arg0: 5},
		{OpCode: vm.Const, Arg0: 0},
		{OpCode: vm.Return},
		{OpCode: vm.Const, Arg0: 1},
		{OpCode: vm.Return},
	}, s("Hello there!"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Arg0: 4},
		{OpCode: vm.CondJump, Arg0: 4},
		{OpCode: vm.Const, Arg0: 0},
		{OpCode: vm.Return},
		{OpCode: vm.Const, Arg0: 1},
		{OpCode: vm.Return},
	}, s("The first bit of data"))
}

func TestPanic(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectError(a.ErrStr("i blew up!"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Arg0: 2},
		{OpCode: vm.Panic},
	}, a.Nil)
}
