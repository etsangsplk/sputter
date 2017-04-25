package cli

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/chzyer/readline"
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	p "github.com/kode4food/sputter/parser"
)

var anyChar = regexp.MustCompile(".")

const replBuiltIns = "*repl-builtins*"

const (
	domain = cyan + "%s" + reset + " "
	prompt = domain + "[%d]> " + code
	cont   = domain + "[%d]" + dgray + "␤   " + code

	output = bold + "%s" + reset
	good   = domain + result + "[%d]= " + output
	bad    = domain + red + "[%d]! " + output
)

type any interface{}

type empty struct{}

var nothing = a.Atom("nothing")

var farewells = []string{
	"Adiós!",
	"Au revoir!",
	"Auf Wiedersehen",
	"B'bye!",
	"Bye!",
	"Bye for now!",
	"Ciao!",
	"Have a wonderful day!",
	"再见!",
	"じゃあね",
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

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: getHistoryFile(),
	})

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
	help(nil, a.EmptyList)
	repl.setInitialPrompt()

	for {
		line, err := repl.rl.Readline()
		repl.buf.WriteString(line + "\n")
		fmt.Print(reset)

		if err != nil {
			emptyBuffer := isEmptyString(repl.buf.String())
			if err == readline.ErrInterrupt && !emptyBuffer {
				repl.reset()
				continue
			}
			break
		}

		if isEmptyString(line) {
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
	repl.setPrompt(fmt.Sprintf(prompt, ns, repl.idx))
}

func (repl *REPL) setContinuePrompt() {
	repl.setPrompt(fmt.Sprintf(cont, repl.nsSpace(), repl.idx))
}

func (repl *REPL) setPrompt(s string) {
	repl.rl.SetPrompt(s)	
}

func (repl *REPL) nsSpace() string {
	ns := string(repl.ns.Domain())
	return anyChar.ReplaceAllString(ns, " ")
}

func (repl *REPL) evalBuffer() (completed bool) {
	defer func() {
		if rec := recover(); rec != nil {
			if isRecoverable(rec.(string)) {
				completed = false
				return
			}
			repl.outputError(rec)
			completed = true
		}
	}()

	l := p.NewLexer(repl.buf.String())
	tr := p.NewReader(repl.ctx, l)
	res := p.EvalReader(repl.ctx, tr)
	repl.outputResult(res)
	return true
}

func (repl *REPL) outputResult(v any) {
	if v == nothing {
		return
	}
	var sv any
	if s, ok := v.(a.Value); ok {
		sv = string(s.Str())
	} else {
		sv = v
	}
	res := fmt.Sprintf(good, repl.nsSpace(), repl.idx, sv)
	fmt.Println(res)
}

func (repl *REPL) outputError(err any) {
	res := fmt.Sprintf(bad, repl.nsSpace(), repl.idx, err)
	fmt.Println(res)
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isRecoverable(err string) bool {
	return err == p.ListNotClosed ||
		err == p.VectorNotClosed ||
		err == p.MapNotClosed ||
		err == p.UnexpectedEndOfFile
}

func use(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := a.AssertUnqualified(args.First()).Name()
	ns := a.GetNamespace(n)
	c.Delete(a.ContextDomain)
	c.Put(a.ContextDomain, ns)
	return ns
}

func shutdown(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	t := time.Now().UTC().UnixNano()
	rs := rand.NewSource(t)
	rg := rand.New(rs)
	idx := rg.Intn(len(farewells))
	fmt.Println(farewells[idx])
	os.Exit(0)
	return nothing
}

func cls(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	fmt.Println(clear)
	return nothing
}

func formatForREPL(s string) string {
	md := formatMarkdown(s)
	lines := strings.Split(md, "\n")
	var out []string
	out = append(out, "")
	for _, l := range lines {
		if isEmptyString(l) {
			out = append(out, l)
		} else {
			out = append(out, "  "+l)
		}
	}
	out = append(out, "")
	return strings.Join(out, "\n")
}

func help(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 0)
	md := string(d.Get("repl-help"))
	fmt.Println(formatForREPL(md))
	return nothing
}

func doc(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	sym := a.AssertUnqualified(args.First())
	if v, ok := c.Get(sym.Name()); ok {
		if an, ok := v.(a.Annotated); ok {
			md := an.Metadata()
			doc := string(md[a.MetaDoc].(a.Str))
			f := formatForREPL(doc)
			fmt.Println(f)
			return nothing
		}
		panic(a.Err("Symbol is not documented: %s", sym))
	}
	panic(a.Err("Could not resolve symbol: %s", sym))
}

func getBuiltInsNamespace() a.Namespace {
	return a.GetNamespace(a.BuiltInDomain)
}

func registerBuiltIn(v a.Annotated) {
	ns := getBuiltInsNamespace()
	if _, ok := ns.Get(replBuiltIns); !ok {
		ns.Put(replBuiltIns, a.Vector{})
	}
	vec, _ := ns.Get(replBuiltIns)
	bi := append(vec.(a.Vector), v.(a.Value))
	ns.Delete(replBuiltIns)
	ns.Put(replBuiltIns, bi)

	n := v.Metadata()[a.MetaName].(a.Name)
	ns.Put(n, v.(a.Value))
}

func registerREPLBuiltIns() {
	registerBuiltIn(
		a.NewFunction(use).WithMetadata(a.Metadata{
			a.MetaName: a.Name("use"),
			a.MetaDoc:  d.Get("repl-use"),
		}),
	)

	registerBuiltIn(
		a.NewFunction(shutdown).WithMetadata(a.Metadata{
			a.MetaName: a.Name("quit"),
			a.MetaDoc:  d.Get("repl-quit"),
		}),
	)

	registerBuiltIn(
		a.NewFunction(cls).WithMetadata(a.Metadata{
			a.MetaName: a.Name("cls"),
			a.MetaDoc:  d.Get("repl-cls"),
		}),
	)

	registerBuiltIn(
		a.NewFunction(help).WithMetadata(a.Metadata{
			a.MetaName: a.Name("help"),
			a.MetaDoc:  d.Get("repl-help"),
		}),
	)

	registerBuiltIn(
		a.NewFunction(doc).WithMetadata(a.Metadata{
			a.MetaName: a.Name("doc"),
			a.MetaDoc:  d.Get("repl-doc"),
		}),
	)
}

func getScreenWidth() int {
	return readline.GetScreenWidth()
}

func getHistoryFile() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return path.Join(usr.HomeDir, ".sputter-history")
}

func init() {
	registerREPLBuiltIns()
}
