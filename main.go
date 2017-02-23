package main

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
)

const (
	noFileSpecified = "No file specified"
	fileNotFound    = "File not found: %s"
)

func evalString(src string) a.Value {
	l := r.NewLexer(string(src))
	tr := r.NewReader(b.Context, l)
	return r.EvalReader(a.NewContext(), tr)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(noFileSpecified)
		os.Exit(-2)
	}

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			os.Exit(-3)
			return
		}
	}()

	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		evalString(string(buffer))
	} else {
		fmt.Println(fmt.Sprintf(fileNotFound, filename))
		os.Exit(-1)
	}
}
