package cli

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	a "github.com/kode4food/sputter/api"
	r "github.com/kode4food/sputter/reader"

	"github.com/chzyer/readline"
)

var anyChar = regexp.MustCompile(".")

const replBuiltIns = "*repl-builtins*"

const (
	red    = "\033[31m"
	green  = "\033[32m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
	gray   = "\033[90m"
	bold   = "\033[1m"
	reset  = "\033[0m"

	domain = cyan + "%s" + reset + " "
	prompt = domain + "[%d]> "
	cont   = domain + "[%d]" + gray + "␤   " + reset

	output = bold + "%s" + reset
	good   = domain + green + "[%d]= " + output
	bad    = domain + red + "[%d]! " + output
)

var farewells = []string{
	"Until next time...",
	"Ciao!",
	"Bye!",
	"Bye for now!",
	"Have a wonderful day!",
	"Goodbye",
	"B'bye!",
}

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

	return repl
}

// Run will perform the Eval-Print-Loop
func (repl *REPL) Run() {
	defer repl.rl.Close()

	fmt.Println(a.Language, a.Version)
	displayStandardHelp()
	repl.setInitialPrompt()

	for {
		line, err := repl.rl.Readline()
		repl.buf.WriteString(line + "\n")

		if err != nil {
			emptyBuffer := isEmpty(repl.buf.String())
			if err == readline.ErrInterrupt && !emptyBuffer {
				repl.reset()
				continue
			}
			break
		}

		if isEmpty(line) {
			continue
		}

		if !repl.evalBuffer() {
			repl.setContinuePrompt()
			continue
		}

		repl.reset()
	}
	shutdown(nil, a.EmptyList)
}

func (repl *REPL) reset() {
	repl.buf.Reset()
	repl.idx++
	repl.setInitialPrompt()
}

func (repl *REPL) setInitialPrompt() {
	if a.GetContextNamespace(repl.ctx) != repl.ns {
		fmt.Println()
		repl.ns = a.GetContextNamespace(repl.ctx)
	}

	ns := repl.ns.Domain()
	repl.rl.SetPrompt(fmt.Sprintf(prompt, ns, repl.idx))
}

func (repl *REPL) setContinuePrompt() {
	repl.rl.SetPrompt(fmt.Sprintf(cont, repl.nsSpace(), repl.idx))
}

func (repl *REPL) nsSpace() string {
	ns := string(repl.ns.Domain())
	return anyChar.ReplaceAllString(ns, " ")
}

func (repl *REPL) evalBuffer() (completed bool) {
	defer func() {
		if rec := recover(); rec != nil {
			if isRecoverable(rec) {
				completed = false
				return
			}
			repl.error(rec)
			completed = true
		}
	}()

	l := r.NewLexer(repl.buf.String())
	tr := r.NewReader(repl.ctx, l)
	res := r.EvalReader(repl.ctx, tr)
	repl.output(res)
	return true
}

func (repl *REPL) output(v a.Value) {
	res := fmt.Sprintf(good, repl.nsSpace(), repl.idx, a.String(v))
	fmt.Println(res)
}

func (repl *REPL) error(err a.Value) {
	res := fmt.Sprintf(bad, repl.nsSpace(), repl.idx, err)
	fmt.Println(res)
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isRecoverable(err a.Value) bool {
	return err == r.ListNotClosed ||
		err == r.VectorNotClosed ||
		err == r.MapNotClosed ||
		err == r.UnexpectedEndOfFile
}

func use(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name
	ns := a.GetNamespace(n)
	c.Delete(a.ContextDomain)
	c.Put(a.ContextDomain, ns)
	return ns
}

func shutdown(c a.Context, args a.Sequence) a.Value {
	t := time.Now().UTC().UnixNano()
	rs := rand.NewSource(t)
	r := rand.New(rs)
	idx := r.Intn(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return a.Nil
}

func displayStandardHelp() {
	ns := getBuiltInsNamespace()
	v, _ := ns.Get(replBuiltIns)
	for _, e := range v.(a.Vector) {
		f := e.(*a.Function)
		fn := fmt.Sprintf("%-8s", "("+string(f.Name)+")")
		fmt.Println(yellow + fn + reset + "; " + f.Documentation())
	}
	fmt.Println()
}

func help(c a.Context, args a.Sequence) a.Value {
	if a.Count(args) == 0 {
		displayStandardHelp()
		return a.Nil
	}
	sym := a.AssertUnqualified(args.First())
	if v, ok := c.Get(sym.Name); ok {
		if d, ok := v.(a.Documented); ok {
			return d.Documentation()
		}
		panic(fmt.Sprintf("Symbol '%s' is not documented", sym))
	}
	panic(fmt.Sprintf("Could not resolve symbol '%s'", sym))
}

func getBuiltInsNamespace() a.Namespace {
	return a.GetNamespace(a.BuiltInDomain)
}

func putBuiltIn(f *a.Function) {
	ns := getBuiltInsNamespace()
	if _, ok := ns.Get(replBuiltIns); !ok {
		ns.Put(replBuiltIns, a.Vector{})
	}
	v, _ := ns.Get(replBuiltIns)
	bi := append(v.(a.Vector), f)
	ns.Delete(replBuiltIns)
	ns.Put(replBuiltIns, bi)
	ns.Put(f.Name, f)
}

func registerREPLBuiltIns() {
	putBuiltIn(&a.Function{
		Name: "use",
		Doc:  "Change namespace. Example: (use foo)",
		Exec: use,
	})
	putBuiltIn(&a.Function{
		Name: "quit",
		Doc:  "Quit the REPL",
		Exec: shutdown,
	})
	putBuiltIn(&a.Function{
		Name: "help",
		Doc:  "Display this help",
		Exec: help,
	})
}

func init() {
	registerREPLBuiltIns()
}
