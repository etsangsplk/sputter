package builtins

import a "github.com/kode4food/sputter/api"

type branch struct {
	condition a.Value
	steps     a.Vector
}

type branches []branch

func makeBranch(e a.Value) branch {
	b := a.AssertSequence(e)
	a.AssertMinimumArity(b, 2)
	i := a.Iterate(b)

	c, _ := i.Next()
	s := make(a.Vector, 0)
	for n, ok := i.Next(); ok; n, ok = i.Next() {
		s = append(s, n)
	}

	return branch{
		condition: c,
		steps:     s,
	}
}

func makeCond(b branches) a.SequenceProcessor {
	return func(c a.Context, args a.Sequence) a.Value {
		for _, e := range b {
			if a.Truthy(a.Eval(c, e.condition)) {
				return a.EvalSequence(c, e.steps)
			}
		}
		return a.Nil
	}
}

func cond(_ a.Context, form a.Sequence) a.Value {
	i := a.Iterate(form)
	i.Next() // we're already here

	b := []branch{}
	for e, ok := i.Next(); ok; e, ok = i.Next() {
		b = append(b, makeBranch(e))
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
