package docstring

import a "github.com/kode4food/sputter/api"

// Get resolves documentation using assets produced by go-bindata
func Get(n string) a.Str {
	return a.Str(MustAsset(filename(n)))
}

// Exists returns whether or not a specific docstring exists
func Exists(n string) bool {
	_, err := AssetInfo(filename(n))
	return err == nil
}

func filename(n string) string {
	return n + ".md"
}
