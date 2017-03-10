package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	a "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"
)

const fileNotFound = "File not found: %s"

// Exec reads the specific source file and evaluates it
func Exec() {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			os.Exit(-2)
		}
	}()

	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		evalSource(string(buffer))
	} else {
		fmt.Println(fmt.Sprintf(fileNotFound, filename))
		os.Exit(-1)
	}
}

func evalSource(src string) a.Value {
	c := a.NewEvalContext()
	l := r.NewLexer(src)
	tr := r.NewReader(c, l)
	return r.EvalReader(a.ChildContext(c), tr)
}
