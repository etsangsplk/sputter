package api_test

import a "github.com/kode4food/sputter/api"

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}
