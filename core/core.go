package core

import (
	"sort"

	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

func init() {
	names := AssetNames()
	sort.Strings(names)
	for _, name := range names {
		src := a.Str(MustAsset(name))
		e.EvalStr(b.Namespace, src)
	}
}
