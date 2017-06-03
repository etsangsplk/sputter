package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
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

func makeEvaluator(c a.Context, src a.Str) a.Sequence {
	l := e.NewLexer(src)
	r := e.NewReader(c, l)
	e := e.Expand(c, r).(a.Sequence)
	return e
}

func evalBuffer(src []byte) a.Value {
	c := e.NewEvalContext()
	s := makeEvaluator(c, a.Str(src))
	return s.Eval(a.ChildContext(c))
}

func exitWithError() {
	if rec := recover(); rec != nil {
		fmt.Println(rec)
		os.Exit(-2)
	}
}
