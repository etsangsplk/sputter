package builtins

import a "github.com/kode4food/sputter/api"

const ifName = "if"

type ifFunction struct{ BaseBuiltIn }

func (*ifFunction) Apply(c a.Context, args a.Sequence) a.Value {
	i := a.AssertArityRange(args, 2, 3)
	f, r, _ := args.Split()
	if a.Truthy(a.Eval(c, f)) {
		return a.Eval(c, r.First())
	}
	if i == 3 {
		return a.Eval(c, r.Rest().First())
	}
	return a.Nil
}

func init() {
	var _if *ifFunction

	RegisterBuiltIn(ifName, _if)
}
