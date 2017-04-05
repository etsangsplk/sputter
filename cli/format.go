package cli

import (
	"bytes"
	"regexp"
	"strings"

	term "github.com/wayneashleyberry/terminal-dimensions"
)

const (
	esc      = "\033["
	red      = esc + "31m"
	green    = esc + "32m"
	yellow   = esc + "33m"
	blue     = esc + "34m"
	magenta  = esc + "35m"
	cyan     = esc + "36m"
	lgray    = esc + "37m"
	dgray    = esc + "90m"
	lred     = esc + "91m"
	lgreen   = esc + "92m"
	lyellow  = esc + "93m"
	lblue    = esc + "94m"
	lmagenta = esc + "95m"
	lcyan    = esc + "96m"
	white    = esc + "97m"
	bold     = esc + "1m"
	reset    = esc + "0m"
	clear    = esc + "2J" + esc + "f"
)

const (
	h1   = lyellow
	h2   = yellow
	code = lblue
)

// This is *not* a full-featured markdown formatter, or even a compliant
// one for that matter.  It only supports the productions that are
// currently used by documentation strings, and will likely not evolve
// much beyond that

type formatter func(string) string

var markdownFormatters = map[*regexp.Regexp]formatter{
	regexp.MustCompile("^# .*$"):     formatHeader1,
	regexp.MustCompile("^## .*$"):    formatHeader2,
	regexp.MustCompile("^  .*$"):     formatIndent,
	regexp.MustCompile("(`.*`)"):     formatCode,
	regexp.MustCompile("([*].*[*])"): formatBold,
}

func getWidth() int {
	if w, err := term.Width(); err == nil {
		return int(w) - 4
	}
	return 76
}

func wrapLine(s string, w int) []string {
	r := []string{}
	var b bytes.Buffer
	for _, e := range strings.Split(s, " ") {
		l := b.Len()
		if l > 0 {
			if l+len(e)+1 >= w {
				r = append(r, b.String())
				b.Reset()
			} else {
				b.WriteString(" ")
			}
		}
		b.WriteString(e)
	}
	r = append(r, b.String())
	return r
}

func wrapLines(s []string) []string {
	w := getWidth()
	r := []string{}
	for _, e := range s {
		r = append(r, wrapLine(e, w)...)
	}
	return r
}

func formatMarkdown(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var out []string
	for _, l := range wrapLines(lines) {
		for r, f := range markdownFormatters {
			l = r.ReplaceAllStringFunc(l, f)
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

func formatHeader1(s string) string {
	return h1 + s[2:] + reset
}

func formatHeader2(s string) string {
	return h2 + s[3:] + reset
}

func formatIndent(s string) string {
	return code + s + reset
}

func formatCode(s string) string {
	return code + s[1:len(s)-1] + reset
}

func formatBold(s string) string {
	return bold + s[1:len(s)-1] + reset
}
