package builtins

import a "github.com/kode4food/sputter/api"

const ifName = "if"

type ifFunction struct{ BaseBuiltIn }

func (*ifFunction) Apply(c a.Context, args a.Vector) a.Value {
	i := a.AssertArityRange(args, 2, 3)
	if a.Truthy(a.Eval(c, args[0])) {
		return a.Eval(c, args[1])
	}
	if i == 3 {
		return a.Eval(c, args[2])
	}
	return a.Nil
}

func init() {
	var _if *ifFunction

	RegisterBuiltIn(ifName, _if)
}
