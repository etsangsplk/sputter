package builtins_test

import (
	"bytes"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func testOutput(t *testing.T, src string, expected string) {
	as := assert.New(t)

	ns := a.GetNamespace(a.UserDomain)
	o, _ := ns.Get("*stdout*")
	ns.Delete("*stdout*")
	w := bytes.NewBufferString("")
	ns.Put("*stdout*", a.NewNative(w))

	runCode(src)

	ns.Delete("*stdout*")
	ns.Put("*stdout*", o)
	as.String(expected, w.String())
}

func TestIO(t *testing.T) {
	testOutput(t, `
		(println "hello" "there")
	`, "hello there\n")

	testOutput(t, `
		(print "hello" "there")
	`, "hello there")

	testOutput(t, `
		(print "hello" 99)
	`, "hello 99")

	testOutput(t, `
		(prn "hello" "there")
	`, "\"hello\" \"there\"\n")

	testOutput(t, `
		(pr "hello" "there")
	`, "\"hello\" \"there\"")

	testOutput(t, `
		(pr "hello" 99)
	`, "\"hello\" 99")
}
