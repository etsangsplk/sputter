package main

import (
	"os"

	_ "github.com/kode4food/sputter/builtins"
	c "github.com/kode4food/sputter/cli"
	_ "github.com/kode4food/sputter/core"
)

func main() {
	if isStdInPiped() {
		c.EvaluateStdIn()
	} else if len(os.Args) < 2 {
		c.NewREPL().Run()
	} else {
		c.EvaluateFile()
	}
}

func isStdInPiped() bool {
	s, _ := os.Stdin.Stat()
	return (s.Mode() & os.ModeCharDevice) == 0
}
