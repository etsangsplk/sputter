package main

import (
	"fmt"
	"io/ioutil"
	"os"

	s "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No file specified")
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
		c := r.NewCoder(b.Context, l)
		r.EvalCoder(s.NewContext(), c)
	} else {
		fmt.Println("File not found:", filename)
		os.Exit(-1)
	}
}
