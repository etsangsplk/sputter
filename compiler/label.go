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

type labelMap map[uint]uint

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
	for i, elem := range inst {
		if elem.OpCode == vm.JumpLabel || elem.OpCode == vm.CondJumpLabel {
			if addr, ok := m[elem.Op1]; ok {
				e2 := elem
				if elem.OpCode == vm.JumpLabel {
					e2.OpCode = vm.Jump
				} else {
					e2.OpCode = vm.CondJump
				}
				e2.Op1 = addr
				elem = e2
			} else {
				panic(a.ErrStr(UnknownLabel, elem.Op1))
			}
		}
		r[i] = elem
	}
	return r
}
