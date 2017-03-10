package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	a "github.com/kode4food/sputter/api"
	_ "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"

	"github.com/chzyer/readline"
)

const (
	fileNotFound = "File not found: %s"
)

func isReadable(src string) (ok bool) {
	defer func() {
		if rec := recover(); rec != nil {
			ok = false
		}
	}()

	c := a.NewContext()
	l := r.NewLexer(src)
	tr := r.NewReader(c, l)
	for v := tr.Next(); v != r.EndOfReader; v = tr.Next() {
	}
	return true
}

func repl() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "> ",
		HistoryFile:            "/tmp/sputter-history",
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	c := a.NewEvalContext()
	var b bytes.Buffer

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		b.WriteString(line)
		b.WriteString("\n")
		src := b.String()

		if !isReadable(src) {
			rl.SetPrompt(">>> ")
			continue
		}

		l := r.NewLexer(src)
		tr := r.NewReader(c, l)
		v := r.EvalReader(a.ChildContext(c), tr)
		fmt.Println(v)

		b.Reset()
		rl.SetPrompt("> ")
	}
}

func evalString(src string) a.Value {
	l := r.NewLexer(src)
	c := a.NewEvalContext()
	tr := r.NewReader(c, l)
	return r.EvalReader(c, tr)
}

func cli() {
	filename := os.Args[1]
	if buffer, err := ioutil.ReadFile(filename); err == nil {
		evalString(string(buffer))
	} else {
		fmt.Println(fmt.Sprintf(fileNotFound, filename))
		os.Exit(-1)
	}
}

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			os.Exit(-3)
		}
	}()

	if len(os.Args) < 2 {
		repl()
	} else {
		cli()
	}
}
