package vm_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	b "github.com/kode4food/sputter/builtins"
	"github.com/kode4food/sputter/vm"
)

var (
	str, _ = b.GetFunction("str")

	vmTestData = a.Values{
		s("The first bit of data"),
		s("Hello there!"),
		a.ErrStr("i blew up!"),
		a.True,
		a.False,
		a.EmptyList,
		str,
		f(20),
		f(3),
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

func testBinary(t *testing.T, o vm.OpCode, expect a.Value) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 8},
		{OpCode: vm.Const, Op1: 7},
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
	testMapped(t, vm.NegOne, f(-1))
}

func TestBinaryMath(t *testing.T) {
	testBinary(t, vm.Add, f(23))
	testBinary(t, vm.Sub, f(17))
	testBinary(t, vm.Mul, f(60))
	testBinary(t, vm.Div, f(20).Div(f(3)))
	testBinary(t, vm.Mod, f(2))

	testBinary(t, vm.Eq, a.False)
	testBinary(t, vm.Neq, a.True)
	testBinary(t, vm.Gt, a.True)
	testBinary(t, vm.Gte, a.True)
	testBinary(t, vm.Lt, a.False)
	testBinary(t, vm.Lte, a.False)
}

func TestConst(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 3},
		{OpCode: vm.Return},
	}, a.True)

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.StoreConst, Op1: 1, Op2: vm.Variables + 0},
		{OpCode: vm.Load, Op1: vm.Variables + 0},
		{OpCode: vm.Return},
	}, s("Hello there!"))
}

func TestSwap(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 0},
		{OpCode: vm.Const, Op1: 1},
		{OpCode: vm.Swap},
		{OpCode: vm.Return},
	}, s("The first bit of data"))
}

func TestReturn(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 1},
		{OpCode: vm.Return},
	}, s("Hello there!"))
}

func TestSequences(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args},
		{OpCode: vm.IsSeq},
		{OpCode: vm.Return},
	}, a.True)

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 3},
		{OpCode: vm.IsSeq},
		{OpCode: vm.Return},
	}, a.False)

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args},
		{OpCode: vm.First},
		{OpCode: vm.Return},
	}, s("first"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args},
		{OpCode: vm.Rest},
		{OpCode: vm.Return},
	}, s(`("second" "third")`))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args},
		{OpCode: vm.Const, Op1: 1},
		{OpCode: vm.Prepend},
		{OpCode: vm.Return},
	}, s(`("Hello there!" "first" "second" "third")`))
}

func TestSequenceSplit(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 3}, // boolean
		{OpCode: vm.Split},
		{OpCode: vm.Return},
	}, a.False)

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 5}, // empty list
		{OpCode: vm.Split},
		{OpCode: vm.Return},
	}, a.False)

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args}, // list
		{OpCode: vm.Split},
		{OpCode: vm.Return},
	}, a.True)

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args}, // list
		{OpCode: vm.Split},
		{OpCode: vm.CondJump, Op1: 5},
		{OpCode: vm.Const, Op1: 2}, // error
		{OpCode: vm.Panic},
		{OpCode: vm.Return},
	}, s("first"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args}, // list
		{OpCode: vm.Split},
		{OpCode: vm.CondJump, Op1: 5},
		{OpCode: vm.Const, Op1: 2}, // error
		{OpCode: vm.Panic},
		{OpCode: vm.Pop},
		{OpCode: vm.Return},
	}, s(`("second" "third")`))
}

func TestDup(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Load, Op1: vm.Args},
		{OpCode: vm.Dup},
		{OpCode: vm.First},
		{OpCode: vm.Dup},
		{OpCode: vm.Store, Op1: vm.Variables + 0},
		{OpCode: vm.Prepend},
		{OpCode: vm.Load, Op1: vm.Variables + 0},
		{OpCode: vm.Prepend},
		{OpCode: vm.Return},
	}, s(`("first" "first" "first" "second" "third")`))
}

func TestIncDec(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.One},
		{OpCode: vm.Inc},
		{OpCode: vm.Return},
	}, f(2))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.One},
		{OpCode: vm.Dec},
		{OpCode: vm.Return},
	}, a.Zero)
}

func TestClear(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Clear, Op1: vm.Args},
		{OpCode: vm.Load, Op1: vm.Args},
		{OpCode: vm.Return},
	}, a.Nil)
}

func TestJump(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 1},
		{OpCode: vm.Jump, Op1: 3},
		{OpCode: vm.Const, Op1: 0},
		{OpCode: vm.Return},
	}, s("Hello there!"))
}

func TestCondJump(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 3},
		{OpCode: vm.NoOp},
		{OpCode: vm.CondJump, Op1: 5},
		{OpCode: vm.Const, Op1: 0},
		{OpCode: vm.Return},
		{OpCode: vm.Const, Op1: 1},
		{OpCode: vm.Return},
	}, s("Hello there!"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 4},
		{OpCode: vm.CondJump, Op1: 4},
		{OpCode: vm.Const, Op1: 0},
		{OpCode: vm.Return},
		{OpCode: vm.Const, Op1: 1},
		{OpCode: vm.Return},
	}, s("The first bit of data"))
}

func TestEval(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 0}, // string
		{OpCode: vm.Eval},
		{OpCode: vm.Return},
	}, s("The first bit of data"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 5}, // empty list
		{OpCode: vm.Eval},
		{OpCode: vm.Return},
	}, a.EmptyList)
}

func TestCallAndApply(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 6}, // func
		{OpCode: vm.Const, Op1: 1}, // string
		{OpCode: vm.Const, Op1: 0}, // string
		{OpCode: vm.Call, Op1: 2},
		{OpCode: vm.Return},
	}, s("The first bit of dataHello there!"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 6}, // func
		{OpCode: vm.Const, Op1: 1}, // string
		{OpCode: vm.Const, Op1: 0}, // string
		{OpCode: vm.Vector, Op1: 2},
		{OpCode: vm.Apply},
		{OpCode: vm.Return},
	}, s("The first bit of dataHello there!"))
}

func TestPanic(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectError(a.ErrStr("i blew up!"))

	testInstructions(t, []vm.Instruction{
		{OpCode: vm.Const, Op1: 2},
		{OpCode: vm.Panic},
	}, a.Nil)
}
