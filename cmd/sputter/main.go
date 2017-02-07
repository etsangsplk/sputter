package main

import (
	"io/ioutil"
	"os"

	"github.com/kode4food/sputter"
	"fmt"
)

func main() {
	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		context := sputter.Builtins.Child()
		l := sputter.NewLexer(string(buffer))
		c := sputter.NewCoder(l)
		e := sputter.NewExecutor(c)
		e.Exec(context)
	} else {
		fmt.Println("File not found:", filename)
		os.Exit(-1)
	}
}
