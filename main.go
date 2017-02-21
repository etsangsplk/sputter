package main

import (
	"fmt"
	"io/ioutil"
	"os"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
)

const (
	noFileSpecified = "No file specified"
	fileNotFound = "File not found: %s"
)

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
		l := r.NewLexer(string(buffer))
		tr := r.NewReader(b.Context, l)
		r.EvalReader(s.NewContext(), tr)
	} else {
		fmt.Println(fmt.Sprintf(fileNotFound, filename))
		os.Exit(-1)
	}
}
