package main

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/evaluator"
	"github.com/kode4food/sputter/reader"
)

const fileNotFound = "file not found: %s"

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
	ctx := evaluator.NewEvalContext()
	r := reader.ReadStr(a.Str(src))
	e := evaluator.Evaluate(ctx, r)
	if v, ok := a.Last(e); ok {
		return v
	}
	return a.Nil
}

func exitWithError() {
	if rec := recover(); rec != nil {
		if ev, ok := rec.(error); ok {
			fmt.Println(ev.Error())
		} else {
			fmt.Println(rec)
		}
		os.Exit(-2)
	}
}
