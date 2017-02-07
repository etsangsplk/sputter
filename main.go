package main

import (
	"io/ioutil"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No file specified")
		os.Exit(-2)
	}
	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		context := NewContext()
		l := NewLexer(string(buffer))
		c := NewCoder(BuiltIns, l)
		EvaluateCoder(context, c)
	} else {
		fmt.Println("File not found:", filename)
		os.Exit(-1)
	}
}
