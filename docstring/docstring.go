package docstring

import (
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assets"
)

// Get resolves documentation using assets produced by go-bindata
func Get(n string) a.Str {
	return a.Str(assets.MustAsset(filename(n)))
}

// Exists returns whether or not a specific docstring exists
func Exists(n string) bool {
	_, err := assets.AssetInfo(filename(n))
	return err == nil
}

func filename(n string) string {
	return "docstring/" + n + ".md"
}
