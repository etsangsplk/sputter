package api_test

import (
	"fmt"
	a "github.com/kode4food/sputter/api"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}

func cvtErr(concrete, intf, method string) a.Error {
	err := "interface conversion: %s is not %s: missing method %s"
	return a.ErrStr(fmt.Sprintf(err, concrete, intf, method))
}
