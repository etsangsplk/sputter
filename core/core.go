package core

import (
	"strings"

	"fmt"
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assets"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

const prefix = "core/"

func init() {
	var name string
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println("Couldn't load:", name)
			fmt.Println(rec)
		}
	}()

	for _, name = range assets.AssetNames() {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		src := a.Str(assets.MustGet(name))
		e.EvalStr(b.Namespace, src)
	}
}
