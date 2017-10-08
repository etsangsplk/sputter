package compiler

import (
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/vm"
)

const (
	// UnknownLabel is raised when a jump is attempted to an unknown label
	UnknownLabel = "unknown label: %s"

	// DuplicatedLabel is raised if a label is defined twice
	DuplicatedLabel = "duplicated label: %s"
)

type (
	// InstructionsProcessor transforms a set of Instructions
	InstructionsProcessor func(vm.Instructions) vm.Instructions

	labelMap map[uint]uint
)

// RemoveLabels processes Instructions, turning Labels into Jump indexes
func RemoveLabels(inst vm.Instructions) vm.Instructions {
	if m, ok := gatherLabels(inst); ok {
		inst = stripLabels(inst, m)
		return rewriteJumps(inst, m)
	}
	return inst
}

func gatherLabels(inst vm.Instructions) (labelMap, bool) {
	m := labelMap{}
	i := uint(0)
	for _, e := range inst {
		if e.OpCode != vm.Label {
			i++
			continue
		}
		if _, ok := m[e.Op1]; ok {
			panic(a.ErrStr(DuplicatedLabel, e.Op1))
		}
		m[e.Op1] = i
	}
	return m, len(m) > 0
}

func stripLabels(inst vm.Instructions, m labelMap) vm.Instructions {
	r := make(vm.Instructions, len(inst)-len(m))
	i := 0
	for _, e := range inst {
		if e.OpCode == vm.Label {
			continue
		}
		r[i] = e
		i++
	}
	return r
}

func rewriteJumps(inst vm.Instructions, m labelMap) vm.Instructions {
	r := make(vm.Instructions, len(inst))
	for i, e := range inst {
		if e.OpCode == vm.Jump || e.OpCode == vm.CondJump {
			if addr, ok := m[e.Op1]; ok {
				e = vm.Instruction{
					OpCode: e.OpCode,
					Op1:    addr,
				}
			} else {
				panic(a.ErrStr(UnknownLabel, e.Op1))
			}
		}
		r[i] = e
	}
	return r
}
