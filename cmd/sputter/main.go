package main

import (
	"os"

	_ "github.com/kode4food/sputter/builtins"
	_ "github.com/kode4food/sputter/core"
)

func main() {
	if isStdInPiped() {
		EvaluateStdIn()
	} else if len(os.Args) < 2 {
		NewREPL().Run()
	} else {
		EvaluateFile()
	}
}

func isStdInPiped() bool {
	s, _ := os.Stdin.Stat()
	return (s.Mode() & os.ModeCharDevice) == 0
}
