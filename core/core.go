package core

import (
	"sort"
	"strings"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assets"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

const prefix = "core/"

func init() {
	names := assets.AssetNames()
	sort.Strings(names)
	for _, name := range names {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		src := a.Str(assets.MustAsset(name))
		e.EvalStr(b.Namespace, src)
	}
}
