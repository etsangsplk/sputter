package builtins

import (
	"fmt"

	a "github.com/kode4food/sputter/api"
)

// ExpectedCondResult is raised when a predicate is not paired
const ExpectedCondResult = "expected result for predicate '%s'"

type branch struct {
	predicate a.Value
	result    a.Value
}

type branches []branch

func makeCond(b branches) a.SequenceProcessor {
	return func(c a.Context, args a.Sequence) a.Value {
		for _, e := range b {
			if a.Truthy(a.Eval(c, e.predicate)) {
				return a.Eval(c, e.result)
			}
		}
		return a.Nil
	}
}

func cond(_ a.Context, form a.Sequence) a.Value {
	b := []branch{}
	i := a.Iterate(form.Rest())
	for p, ok := i.Next(); ok; p, ok = i.Next() {
		if r, ok := i.Next(); ok {
			b = append(b, branch{
				predicate: p,
				result:    r,
			})
		} else {
			panic(fmt.Sprintf(ExpectedCondResult, a.String(p)))
		}
	}
	return a.NewList(&a.Function{Exec: makeCond(b)})
}

// this will be replaced by a macro -> cond
func _if(c a.Context, args a.Sequence) a.Value {
	a.AssertArityRange(args, 2, 3)
	i := a.Iterate(args)
	condVal, _ := i.Next()
	cond := a.Eval(c, condVal)
	if !a.Truthy(cond) {
		i.Next()
	}
	result, _ := i.Next()
	return a.Eval(c, result)
}

func init() {
	registerMacro(&a.Macro{Name: "cond", Prep: cond})
	registerFunction(&a.Function{Name: "if", Exec: _if})
}
