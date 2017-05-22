package docstring

import (
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assets"
)

// Get resolves documentation using assets produced by go-bindata
func Get(n string) a.Str {
	return a.Str(assets.MustGet(filename(n)))
}

// Exists returns whether or not a specific docstring exists
func Exists(n string) bool {
	fn := filename(n)
	for _, e := range assets.AssetNames() {
		if fn == e {
			return true
		}
	}
	return false
}

func filename(n string) string {
	return "docstring/" + n + ".md"
}
