package util_test

import (
	"testing"

	"github.com/kode4food/sputter/assert"
	u "github.com/kode4food/sputter/util"
)

func TestCache(t *testing.T) {
	as := assert.New(t)
	c := u.NewCache()
	count := 0
	res := func() u.Any {
		count++
		return "value"
	}

	r1 := c.Get("blah", res)
	as.String("value", r1)
	as.Number(1, count)

	r2 := c.Get("blah", res)
	as.String("value", r2)
	as.Number(1, count)

	r3 := c.Get("nah", res)
	as.String("value", r3)
	as.Number(2, count)

	r4 := c.Get("nah", res)
	as.String("value", r4)
	as.Number(2, count)

	r5 := c.Get("blah", res)
	as.String("value", r5)
	as.Number(2, count)
}
