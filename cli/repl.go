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

	"github.com/kode4food/readline"
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	p "github.com/kode4food/sputter/parser"
)

type any interface{}

type sentinel struct{}

var anyChar = regexp.MustCompile(".")

const replBuiltIns = "*repl-builtins*"

const (
	domain = cyan + "%s" + reset + " "
	prompt = domain + "[%d]> " + code
	cont   = domain + "[%d]" + dgray + "␤   " + code

	output = bold + "%s" + reset
	good   = domain + result + "[%d]= " + output
	bad    = domain + red + "[%d]! " + output
	paired = esc + "7m"
)

var nothing = &sentinel{}

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

var (
	openers = map[rune]rune{')': '(', ']': '[', '}': '{'}
	closers = map[rune]rune{'(': ')', '[': ']', '{': '}'}
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

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: getHistoryFile(),
		Painter:     repl,
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
func (r *REPL) Run() {
	defer r.rl.Close()

	fmt.Println(a.Language, a.Version)
	help(nil, a.EmptyList)
	r.setInitialPrompt()

	for {
		line, err := r.rl.Readline()
		r.buf.WriteString(line + "\n")
		fmt.Print(reset)

		if err != nil {
			emptyBuffer := isEmptyString(r.buf.String())
			if err == readline.ErrInterrupt && !emptyBuffer {
				r.reset()
				continue
			}
			break
		}

		if isEmptyString(line) {
			continue
		}

		if !r.evalBuffer() {
			r.setContinuePrompt()
			continue
		}

		r.reset()
	}
	shutdown(nil, a.EmptyList)
}

func (r *REPL) reset() {
	r.buf.Reset()
	r.idx++
	r.setInitialPrompt()
}

func (r *REPL) setInitialPrompt() {
	if a.GetContextNamespace(r.ctx) != r.ns {
		fmt.Println()
		r.ns = a.GetContextNamespace(r.ctx)
	}

	ns := r.ns.Domain()
	r.setPrompt(fmt.Sprintf(prompt, ns, r.idx))
}

func (r *REPL) setContinuePrompt() {
	r.setPrompt(fmt.Sprintf(cont, r.nsSpace(), r.idx))
}

func (r *REPL) setPrompt(s string) {
	r.rl.SetPrompt(s)
}

func (r *REPL) nsSpace() string {
	ns := string(r.ns.Domain())
	return anyChar.ReplaceAllString(ns, " ")
}

func (r *REPL) evalBuffer() (completed bool) {
	defer func() {
		if rec := recover(); rec != nil {
			if isRecoverable(rec.(string)) {
				completed = false
				return
			}
			r.outputError(rec)
			completed = true
		}
	}()

	l := p.NewLexer(a.Str(r.buf.String()))
	tr := p.NewReader(r.ctx, l)
	res := p.EvalReader(r.ctx, tr)
	r.outputResult(res)
	return true
}

func (r *REPL) outputResult(v any) {
	if v == nothing {
		return
	}
	var sv any
	if s, ok := v.(a.Value); ok {
		sv = string(s.Str())
	} else {
		sv = v
	}
	res := fmt.Sprintf(good, r.nsSpace(), r.idx, sv)
	fmt.Println(res)
}

func (r *REPL) outputError(err any) {
	res := fmt.Sprintf(bad, r.nsSpace(), r.idx, err)
	fmt.Println(res)
}

// Paint implements the Painter interface
func (r *REPL) Paint(line []rune, pos int) []rune {
	if line == nil || len(line) == 0 {
		return line
	}

	l := len(line)
	npos := pos
	if npos < 0 {
		npos = 0
	}
	if npos >= l {
		npos = l - 1
	}
	k := line[npos]
	if _, ok := openers[k]; ok {
		return markOpener(line, npos, k)
	} else if _, ok := closers[k]; ok {
		return markCloser(line, npos, k)
	}
	return line
}

func (s *sentinel) Str() a.Str {
	return ""
}

func markOpener(line []rune, pos int, c rune) []rune {
	o := openers[c]
	depth := 0
	for i := pos; i >= 0; i-- {
		if line[i] == o {
			depth--
			if depth == 0 {
				return markMatch(line, i)
			}
		} else if line[i] == c {
			depth++
		}
	}
	return line
}

func markCloser(line []rune, pos int, o rune) []rune {
	c := closers[o]
	depth := 0
	for i := pos; i < len(line); i++ {
		if line[i] == c {
			depth--
			if depth == 0 {
				return markMatch(line, i)
			}
		} else if line[i] == o {
			depth++
		}
	}
	return line
}

func markMatch(line []rune, pos int) []rune {
	m := []rune{}
	if pos > 0 {
		m = append(m, line[:pos]...)
	}
	m = append(m, []rune(paired)...)
	m = append(m, line[pos])
	m = append(m, []rune(reset+code)...)
	if pos < len(line)-1 {
		m = append(m, line[pos+1:]...)
	}
	return m
}

func isEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isRecoverable(err string) bool {
	return err == p.ListNotClosed ||
		err == p.VectorNotClosed ||
		err == p.MapNotClosed
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

func registerBuiltIns() {
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
	registerBuiltIns()
}
