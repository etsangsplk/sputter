package cli

import (
	"regexp"
	"strings"
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

type formatter func(string) string

var markdownFormatters = map[*regexp.Regexp]formatter{
	regexp.MustCompile("^# .*$"):     formatHeader1,
	regexp.MustCompile("^## .*$"):    formatHeader2,
	regexp.MustCompile("^  .*$"):     formatIndent,
	regexp.MustCompile("(`.*`)"):     formatCode,
	regexp.MustCompile("([*].*[*])"): formatBold,
}

func formatMarkdown(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var out []string
	for _, l := range lines {
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
