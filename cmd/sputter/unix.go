// +build !windows

package main

const (
	esc      = "\033["
	black    = esc + "30m"
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
	italic   = esc + "3m"
	reset    = esc + "0m"
	clear    = esc + "2J" + esc + "f"
	paired   = esc + "7m"
	nlMarker = "␤"

	h1     = lyellow
	h2     = yellow
	code   = lblue
	result = green
)

var farewells = []string{
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
