package main

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	e "github.com/kode4food/sputter/evaluator"
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

func readSource(c a.Context, src a.Str) a.Sequence {
	l := e.Scan(src)
	return e.Read(c, l)
}

func evalBuffer(src []byte) a.Value {
	c := e.NewEvalContext()
	r := readSource(c, a.Str(src))
	return a.Eval(c, r)
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
