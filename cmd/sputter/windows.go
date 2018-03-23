// +build windows

package main

const (
	esc      = ""
	black    = ""
	red      = ""
	green    = ""
	yellow   = ""
	blue     = ""
	magenta  = ""
	cyan     = ""
	lgray    = ""
	dgray    = ""
	lred     = ""
	lgreen   = ""
	lyellow  = ""
	lblue    = ""
	lmagenta = ""
	lcyan    = ""
	white    = ""
	bold     = ""
	italic   = ""
	reset    = ""
	clear    = ""
	paired   = ""
	nlMarker = "\\"

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
}

// Paint implements the Painter interface
func (r *REPL) Paint(line []rune, pos int) []rune {
	return line
}
