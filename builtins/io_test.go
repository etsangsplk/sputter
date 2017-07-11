package builtins_test

import (
	"bytes"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

const stdoutName = "*stdout*"

func testOutput(t *testing.T, src string, expected string) {
	as := assert.New(t)

	ns := a.GetNamespace(a.BuiltInDomain)
	o, _ := ns.Get(stdoutName)
	ns.Delete(stdoutName)
	b := bytes.NewBufferString("")
	w := a.NewWriter(b, a.StrOutput)
	ns.Put(stdoutName, w)

	runCode(src)

	ns.Delete(stdoutName)
	ns.Put(stdoutName, o)
	as.String(expected, b.String())
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
