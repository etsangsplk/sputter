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

	output = bold + "%s" + reset
	good   = domain + green + "[%d]= " + output
	bad    = domain + red + "[%d]! " + output
)

// REPL manages a Read-Eval-Print Loop
type REPL struct {
	buf bytes.Buffer
	ctx a.Context
	ns  a.Namespace
	rl  *readline.Instance
	idx int
}

// NewREPL instantiates a new REPL instance
func NewREPL() *REPL {
	repl := &REPL{}

	rl, err := readline.NewEx(&readline.Config{})
	if err != nil {
		panic(err)
	}

	repl.rl = rl
	repl.ctx = a.NewEvalContext()
	repl.ns = a.GetContextNamespace(repl.ctx)
	repl.idx = 1

	repl.registerREPLBuiltIns()
	return repl
}

// Run will perform the Eval-Print-Loop
func (repl *REPL) Run() {
	defer repl.rl.Close()

	fmt.Println(a.Language, a.Version)
	repl.setInitialPrompt()

	for {
		line, err := repl.rl.Readline()
		if err != nil {
			break
		}

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		repl.buf.WriteString(line + "\n")
		if !repl.isReadable() {
			repl.setContinuePrompt()
			continue
		}

		repl.evalLine()
		repl.buf.Reset()

		if a.GetContextNamespace(repl.ctx) != repl.ns {
			fmt.Println()
			repl.ns = a.GetContextNamespace(repl.ctx)
		}

		repl.idx++
		repl.setInitialPrompt()
	}
}

func (repl *REPL) setInitialPrompt() {
	ns := repl.ns.Domain()
	repl.rl.SetPrompt(fmt.Sprintf(prompt, ns, repl.idx))
}

func (repl *REPL) setContinuePrompt() {
	repl.rl.SetPrompt(fmt.Sprintf(cont, repl.nsSpace(), repl.idx))
}

func (repl *REPL) nsSpace() string {
	ns := repl.ns.Domain()
	return any.ReplaceAllString(string(ns), " ")
}

func (repl *REPL) isReadable() (ok bool) {
	defer func() {
		if rec := recover(); rec != nil {
			if isRecoverable(rec) {
				ok = false
				return
			}
			ok = true
		}
	}()

	l := r.NewLexer(repl.buf.String())
	tr := r.NewReader(a.ChildContext(repl.ctx), l)
	for v := tr.Next(); v != r.EndOfReader; v = tr.Next() {
	}
	return true
}

func (repl *REPL) evalLine() {
	defer func() {
		if rec := recover(); rec != nil {
			repl.output(bad, rec)
		}
	}()

	l := r.NewLexer(repl.buf.String())
	tr := r.NewReader(repl.ctx, l)
	res := r.EvalReader(repl.ctx, tr)
	repl.output(good, res)
}

func (repl *REPL) output(prefix string, v a.Value) {
	res := a.String(v)
	fmt.Println(fmt.Sprintf(prefix, repl.nsSpace(), repl.idx, res))
}

func (repl *REPL) registerREPLBuiltIns() {
	repl.ctx.Put("use", &a.Function{Name: "use", Apply: use})
}

func use(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name
	ns := a.GetNamespace(n)
	c.Delete(a.ContextDomain)
	c.Put(a.ContextDomain, ns)
	return ns
}

func isRecoverable(err a.Value) bool {
	return err == r.ListNotClosed ||
		err == r.VectorNotClosed ||
		err == r.UnexpectedEndOfFile
}
