package api_test

import "github.com/stretchr/testify/assert"

func expectError(as *assert.Assertions, err string) {
	if rec := recover(); rec != nil {
		as.Equal(err, rec, "error raised")
		return
	}
	as.Fail("error not raised")
}
