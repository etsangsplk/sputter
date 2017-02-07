package main

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	b "github.com/kode4food/sputter/builtins"
	i "github.com/kode4food/sputter/interpreter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No file specified")
		os.Exit(-2)
	}
	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		context := a.NewContext()
		l := i.NewLexer(string(buffer))
		c := i.NewCoder(b.BuiltIns, l)
		i.EvaluateCoder(context, c)
	} else {
		fmt.Println("File not found:", filename)
		os.Exit(-1)
	}
}
