package api_test

import (
	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) *a.Number {
	return a.NewFloat(f)
}

func expectError(as *assert.Wrapper, err string) {
	if rec := recover(); rec != nil {
		as.String(err, rec)
		return
	}
	as.Fail("error not raised")
}
