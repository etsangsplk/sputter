package api_test

import "github.com/kode4food/sputter/assert"

func expectError(as *assert.Wrapper, err string) {
	if rec := recover(); rec != nil {
		as.Equal(err, rec)
		return
	}
	as.Fail("error not raised")
}
