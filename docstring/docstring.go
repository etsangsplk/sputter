package docstring

// Get resolves documentation using assets produced by go-bindata
func Get(n string) string {
	return string(MustAsset(n + ".md"))
}
