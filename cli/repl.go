package cli

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	a "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"

	"github.com/chzyer/readline"
)

var any = regexp.MustCompile(".")

const (
	red   = "\033[31m"
	green = "\033[32m"
	cyan  = "\033[36m"
	bold  = "\033[1m"
	reset = "\033[0m"

	domain = cyan + "%s" + reset + " "
	prompt = domain + "[%d]> "
	cont   = domain + "[%d]>   "

	output = domain + "[%d]%s"
	good   = green + "= " + bold
	bad    = red + "! " + bold
)

// REPL instantiates a Read-Evaluate-Print Loop
func REPL() {
	fmt.Println(a.Language, a.Version)
	rl, err := readline.NewEx(&readline.Config{})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var b bytes.Buffer
	c := a.NewEvalContext()
	ns := a.GetContextNamespace(c)
	i := 1

	c.Put("use", &a.Function{Name: "use", Apply: use})
	rl.SetPrompt(fmt.Sprintf(prompt, ns.Domain(), i))
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		b.WriteString(line + "\n")
		src := b.String()

		if !isReadable(src) {
			rl.SetPrompt(fmt.Sprintf(cont, nsSpace(ns), i))
			continue
		}

		res := evalLine(c, src)
		fmt.Println(fmt.Sprintf(output, nsSpace(ns), i, res))
		b.Reset()

		if a.GetContextNamespace(c) != ns {
			fmt.Println()
			ns = a.GetContextNamespace(c)
		}

		i++
		rl.SetPrompt(fmt.Sprintf(prompt, ns.Domain(), i))
	}
}

func use(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name
	ns := a.GetNamespace(n)
	c.Delete(a.ContextDomain)
	c.Put(a.ContextDomain, ns)
	return ns
}

func nsSpace(ns a.Namespace) string {
	return any.ReplaceAllString(string(ns.Domain()), " ")
}

func isReadable(src string) (ok bool) {
	defer func() {
		if rec := recover(); rec != nil {
			if isRecoverable(rec) {
				ok = false
				return
			}
			ok = true
		}
	}()

	c := a.NewContext()
	l := r.NewLexer(src)
	tr := r.NewReader(c, l)
	for v := tr.Next(); v != r.EndOfReader; v = tr.Next() {
	}
	return true
}

func isRecoverable(err a.Value) bool {
	return err == r.ListNotClosed ||
		err == r.VectorNotClosed ||
		err == r.UnexpectedEndOfFile
}

func evalLine(c a.Context, src string) (result string) {
	defer func() {
		if rec := recover(); rec != nil {
			result = fmt.Sprint(bad, rec, reset)
		}
	}()

	l := r.NewLexer(src)
	tr := r.NewReader(c, l)
	return fmt.Sprint(good, r.EvalReader(c, tr), reset)
}
