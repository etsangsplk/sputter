package core

import (
	"fmt"
	"os"
	"strings"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assets"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

const prefix = "core/"

func init() {
	var filename string

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Fprintf(os.Stderr, "\nBootstrap Error\n\n")
			fmt.Fprintf(os.Stderr, "  %s: %s\n\n", filename, rec)
			os.Exit(-1)
		}
	}()

	for _, filename = range assets.AssetNames() {
		if !strings.HasPrefix(filename, prefix) {
			continue
		}
		src := a.Str(assets.MustGet(filename))
		e.EvalStr(b.Namespace, src)
	}
}
