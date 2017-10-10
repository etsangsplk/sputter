package cli

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chzyer/readline"
	a "github.com/kode4food/sputter/api"
	d "github.com/kode4food/sputter/docstring"
	e "github.com/kode4food/sputter/evaluator"
)

const (
	replBuiltIns = "*repl-builtins*"

	domain = cyan + "%s" + reset + " "
	prompt = domain + "[%d]> " + code
	cont   = domain + "[%d]" + dgray + "␤   " + code

	output = bold + "%s" + reset
	good   = domain + result + "[%d]= " + output
	bad    = domain + red + "[%d]! " + output
	paired = esc + "7m"
)

type (
	any      interface{}
	sentinel struct{}

	// REPL manages a Read-Eval-Print Loop
	REPL struct {
		buf bytes.Buffer
		ctx a.Context
		ns  a.Namespace
		rl  *readline.Instance
		idx int
	}
)

var (
	anyChar = regexp.MustCompile(".")

	nothing = new(sentinel)

	farewells = []string{
		"¡Adiós!",
		"Au revoir!",
		"Bye for now!",
		"Ciao!",
		"Tchau!",
		"Tschüss!",
		"Hoşçakal!",
		"Αντίο!",
		"До свидания!",
		"अलविदा!",
		"안녕!",
		"再见!",
		"じゃあね",
	}

	openers = map[rune]rune{')': '(', ']': '[', '}': '{'}
	closers = map[rune]rune{'(': ')', '[': ']', '{': '}'}
)

// NewREPL instantiates a new REPL instance
func NewREPL() *REPL {
	repl := new(REPL)

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile: getHistoryFile(),
		Painter:     repl,
	})

	if err != nil {
		panic(err)
	}

	repl.rl = rl
	repl.ctx = e.NewEvalContext()
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
		if err := toError(recover()); err != nil {
			if isRecoverable(err) {
				completed = false
				return
			}
			r.outputError(err)
			completed = true
		}
	}()

	res := e.EvalStr(r.ctx, a.Str(r.buf.String()))
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

func (r *REPL) outputError(err error) {
	msg := err.Error()
	res := fmt.Sprintf(bad, r.nsSpace(), r.idx, msg)
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

func toError(v interface{}) error {
	if v == nil {
		return nil
	}
	if e, ok := v.(error); ok {
		return e
	}
	panic(fmt.Sprintf("non-standard error: %s", v))
}

func isRecoverable(err error) bool {
	msg := err.Error()
	return msg == e.ListNotClosed ||
		msg == e.VectorNotClosed ||
		msg == e.MapNotClosed
}

func use(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 1)
	n := args.First().(a.LocalSymbol).Name()
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

func debugInfo(_ a.Context, _ a.Sequence) a.Value {
	runtime.GC()
	fmt.Println("Number of goroutines: ", runtime.NumGoroutine())
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
	sym := args.First().(a.LocalSymbol)
	if v, ok := c.Get(sym.Name()); ok {
		if vd, ok := v.(a.Documented); ok {
			docStr := vd.Documentation()
			f := formatForREPL(string(docStr))
			fmt.Println(f)
			return nothing
		}
		panic(a.ErrStr("symbol is not documented: %s", sym))
	}
	panic(a.ErrStr("could not resolve symbol: %s", sym))
}

func getBuiltInsNamespace() a.Namespace {
	return a.GetNamespace(a.BuiltInDomain)
}

func registerBuiltIn(v a.AnnotatedValue) {
	ns := getBuiltInsNamespace()
	if _, ok := ns.Get(replBuiltIns); !ok {
		ns.Put(replBuiltIns, a.Values{})
	}
	vec, _ := ns.Get(replBuiltIns)
	bi := vec.(a.Vector).Conjoin(v)
	ns.Delete(replBuiltIns)
	ns.Put(replBuiltIns, bi)

	n := v.Metadata().MustGet(a.NameKey).(a.Name)
	ns.Put(n, v)
}

func registerBuiltIns() {
	registerBuiltIn(
		a.NewExecFunction(use).WithMetadata(a.Properties{
			a.NameKey:     a.Name("use"),
			a.DocAssetKey: a.Str("repl-use"),
			a.SpecialKey:  a.True,
		}),
	)

	registerBuiltIn(
		a.NewExecFunction(shutdown).WithMetadata(a.Properties{
			a.NameKey:     a.Name("quit"),
			a.DocAssetKey: a.Str("repl-quit"),
		}),
	)

	registerBuiltIn(
		a.NewExecFunction(debugInfo).WithMetadata(a.Properties{
			a.NameKey: a.Name("debug-info"),
		}),
	)

	registerBuiltIn(
		a.NewExecFunction(cls).WithMetadata(a.Properties{
			a.NameKey:     a.Name("cls"),
			a.DocAssetKey: a.Str("repl-cls"),
		}),
	)

	registerBuiltIn(
		a.NewExecFunction(help).WithMetadata(a.Properties{
			a.NameKey:     a.Name("help"),
			a.DocAssetKey: a.Str("repl-help"),
		}),
	)

	registerBuiltIn(
		a.NewExecFunction(doc).WithMetadata(a.Properties{
			a.NameKey:     a.Name("doc"),
			a.DocAssetKey: a.Str("repl-doc"),
			a.SpecialKey:  a.True,
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
