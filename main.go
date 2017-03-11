package main

import (
	"os"

	_ "github.com/kode4food/sputter/builtins"
	c "github.com/kode4food/sputter/cli"
)

func main() {
	if len(os.Args) < 2 {
		c.NewREPL().Run()
	} else {
		c.Exec()
	}
}
