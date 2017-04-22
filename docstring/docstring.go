package docstring

import a "github.com/kode4food/sputter/api"

// Get resolves documentation using assets produced by go-bindata
func Get(n string) a.Str {
	return a.Str(MustAsset(n + ".md"))
}
