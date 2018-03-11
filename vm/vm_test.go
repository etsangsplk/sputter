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
		VarsCount:    16,
		Data:         vmTestData,
		Instructions: inst,
	}

	r1 := m1.Apply(a.Variables{}, vmTestArgs)
	as.Equal(expect, r1)

	m2 := m1.WithMetadata(a.Properties{
		a.NameKey: a.Name("newMetadata"),
	}).(*vm.Module)
	r2 := m2.Apply(a.Variables{}, vmTestArgs)
	as.Equal(expect, r2)
}

func testMapped(t *testing.T, o vm.OpCode, expect a.Value) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(o, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, expect)
}

func testBinary(t *testing.T, o vm.OpCode, expect a.Value) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 7, vm.Vars+0),
		vm.MakeInst(vm.Const, 8, vm.Vars+1),
		vm.MakeInst(o, vm.Vars+0, vm.Vars+1, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
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
		vm.MakeInst(vm.Const, 3, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.True)

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 1, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("Hello there!"))
}

func TestSequences(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.IsSeq, vm.Args, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.True)

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 3),
		vm.MakeInst(vm.IsSeq),
		vm.MakeInst(vm.Return),
	}, a.False)

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.First, vm.Args, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("first"))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Rest, vm.Args, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s(`("second" "third")`))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 1, vm.Vars+0),
		vm.MakeInst(vm.Prepend, vm.Vars+0, vm.Args, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s(`("Hello there!" "first" "second" "third")`))
}

func TestSequenceSplit(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 3, vm.Vars+0), // boolean
		vm.MakeInst(vm.Split, vm.Vars+0, vm.Vars+0, vm.Vars+1, vm.Vars+2),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.False)

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 5, vm.Vars+0), // empty list
		vm.MakeInst(vm.Split, vm.Vars+0, vm.Vars+0, vm.Vars+1, vm.Vars+2),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.False)

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Split, vm.Args, vm.Vars+0, vm.Vars+1, vm.Vars+2),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.True)

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Split, vm.Args, vm.Vars+0, vm.Vars+1, vm.Vars+2),
		vm.MakeInst(vm.CondJump, 4, vm.Vars+0),
		vm.MakeInst(vm.Const, 2, vm.Vars+0), // error
		vm.MakeInst(vm.Panic, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+1),
	}, s("first"))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Split, vm.Args, vm.Vars+0, vm.Vars+1, vm.Vars+2),
		vm.MakeInst(vm.CondJump, 4, vm.Vars+0),
		vm.MakeInst(vm.Const, 2, vm.Vars+0), // error
		vm.MakeInst(vm.Panic, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+2),
	}, s(`("second" "third")`))
}

func TestDup(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Dup, vm.Args, vm.Vars+0),
		vm.MakeInst(vm.First, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s(`("first" "second" "third")`))
}

func TestIncDec(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.One, vm.Vars+0),
		vm.MakeInst(vm.Inc, vm.Vars+0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, f(2))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.One, vm.Vars+0),
		vm.MakeInst(vm.Dec, vm.Vars+0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.Zero)
}

func TestNil(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Nil, vm.Args),
		vm.MakeInst(vm.Return, vm.Args),
	}, a.Nil)
}

func TestJump(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 1, vm.Vars+0),
		vm.MakeInst(vm.Jump, 3),
		vm.MakeInst(vm.Const, 0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("Hello there!"))
}

func TestCondJump(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 3, vm.Vars+0),
		vm.MakeInst(vm.CondJump, 4),
		vm.MakeInst(vm.Const, 0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
		vm.MakeInst(vm.Const, 1, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("Hello there!"))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 4, vm.Vars+0),
		vm.MakeInst(vm.CondJump, 4, vm.Vars+0),
		vm.MakeInst(vm.Const, 0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
		vm.MakeInst(vm.Const, 1, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("The first bit of data"))
}

func TestEval(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 0, vm.Vars+0), // string
		vm.MakeInst(vm.Eval, vm.Vars+0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("The first bit of data"))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 5, vm.Vars+0), // empty list
		vm.MakeInst(vm.Eval, vm.Vars+0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, a.EmptyList)
}

func TestCallAndApply(t *testing.T) {
	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 6, vm.Vars+0), // func
		vm.MakeInst(vm.Const, 0, vm.Vars+1), // string
		vm.MakeInst(vm.Const, 1, vm.Vars+2), // string
		vm.MakeInst(vm.Vector, vm.Vars+1, vm.Vars+3, vm.Vars+1),
		vm.MakeInst(vm.Apply, vm.Vars+1, vm.Vars+0, vm.Vars+0),
		vm.MakeInst(vm.Return, vm.Vars+0),
	}, s("The first bit of dataHello there!"))
}

func TestPanic(t *testing.T) {
	as := assert.New(t)

	defer as.ExpectError(a.ErrStr("i blew up!"))

	testInstructions(t, []vm.Instruction{
		vm.MakeInst(vm.Const, 2, vm.Vars+0),
		vm.MakeInst(vm.Panic, vm.Vars+0),
	}, a.Nil)
}
