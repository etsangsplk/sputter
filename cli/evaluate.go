package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"
)

const fileNotFound = "File not found: %s"

// EvaluateStdIn reads from StdIn and evaluates it
func EvaluateStdIn() {
	defer exitWithError()

	buffer, _ := ioutil.ReadAll(os.Stdin)
	evalBuffer(buffer)
}

// EvaluateFile reads the specific source file and evaluates it
func EvaluateFile() {
	defer exitWithError()

	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		evalBuffer(buffer)
	} else {
		fmt.Println(fmt.Sprintf(fileNotFound, filename))
		os.Exit(-1)
	}
}

func evalBuffer(src []byte) a.Value {
	c := a.NewEvalContext()
	l := r.NewLexer(string(src))
	tr := r.NewReader(c, l)
	return r.EvalReader(a.ChildContext(c), tr)
}

func exitWithError() {
	if rec := recover(); rec != nil {
		fmt.Println(rec)
		os.Exit(-2)
	}
}
