package core

import (
	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	e "github.com/kode4food/sputter/evaluator"
)

// Load resolves built-ins using assets produced by go-bindata
func Load(n string) {
	src := a.Str(MustAsset(n + ".lisp"))
	e.EvalStr(b.Namespace, src)
}

func init() {
	Load("core")
}
