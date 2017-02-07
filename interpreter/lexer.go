package interpreter

import (
	"math/big"
	"regexp"

	a "github.com/kode4food/sputter/api"
)

// UnexpectedEndOfFile is the error returned when EOF is unexpectedly reached
const UnexpectedEndOfFile = "End of file reached unexpectedly"

// UnmatchedState is the error returned when the lexing state is invalid
const UnmatchedState = "Unmatched lexing state"

// TokenType is an opaque type for lexer tokens
type TokenType int

// Token Types
const (
	Error TokenType = iota
	Identifier
	String
	Number
	Ratio
	ListStart
	ListEnd
	ArgsStart
	ArgsEnd
	LiteralMarker
	EndOfFile
	Whitespace
)

// Token is a Lexer token
type Token struct {
	Type  TokenType
	Value a.Value
}

// TokenReader defines the one method that a Token processor must provide
type TokenReader interface {
	Next() *Token
}

// EOFToken marks the end of a TokenReader stream
var EOFToken = &Token{EndOfFile, ""}

// Lexer is the lexer interface
type Lexer struct {
	input  string
	start  int
	pos    int
	tokens chan *Token
}

type stateFunc func(*Lexer) stateFunc

type stateMapEntry struct {
	pattern  *regexp.Regexp
	function stateFunc
}

type stateMap []stateMapEntry

var states stateMap

func init() {
	states = stateMap{
		{regexp.MustCompile(`^$`), endState(EndOfFile)},
		{regexp.MustCompile(`^\s+`), tokenState(Whitespace)},
		{regexp.MustCompile(`^\(`), tokenState(ListStart)},
		{regexp.MustCompile(`^\)`), tokenState(ListEnd)},
		{regexp.MustCompile(`^\[`), tokenState(ArgsStart)},
		{regexp.MustCompile(`^]`), tokenState(ArgsEnd)},
		{regexp.MustCompile(`^'`), tokenState(LiteralMarker)},

		{regexp.MustCompile(`^"(\\.|[^"])*"`), stringState},
		{regexp.MustCompile(`^[1-9]\d*/[1-9]\d*`), ratioState},
		{regexp.MustCompile(`^(0|[1-9]\d*(\.\d+)?([eE][+-]?\d+)?)`), numberState},

		{regexp.MustCompile(`^[^()\[\]\s]+`), tokenState(Identifier)},
		{regexp.MustCompile(`^.`), endState(Error)},
	}
}

// NewLexer instantiates a new Lexer instance
func NewLexer(source string) *Lexer {
	lexer := &Lexer{
		input:  source,
		tokens: make(chan *Token),
	}

	go lexer.run()
	return lexer
}

// Next returns the next Token from the lexer's Token channel
func (l *Lexer) Next() *Token {
	for {
		token, ok := <-l.tokens
		if !ok {
			panic(UnexpectedEndOfFile)
		}
		if token.Type != Whitespace {
			return token
		}
	}
}

func (l *Lexer) run() {
	for state := initState; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *Lexer) emitValue(t TokenType, v a.Value) {
	l.tokens <- &Token{t, v}
	l.skip()
}

func (l *Lexer) emit(t TokenType) {
	l.emitValue(t, l.currentToken())
}

func (l *Lexer) currentToken() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) skip() {
	l.start = l.pos
}

func (l *Lexer) matchState() stateFunc {
	str := l.input[l.pos:]
	for _, entry := range states {
		if match := entry.pattern.FindStringIndex(str); match != nil {
			result := str[match[0]:match[1]]
			l.pos += len(result)
			return entry.function
		}
	}
	panic(UnmatchedState)
}

func tokenState(t TokenType) stateFunc {
	return func(l *Lexer) stateFunc {
		l.emit(t)
		return initState
	}
}

func initState(l *Lexer) stateFunc {
	var state = l.matchState()
	return state(l)
}

func endState(t TokenType) stateFunc {
	return func(l *Lexer) stateFunc {
		l.emit(t)
		return nil
	}
}

func stringState(l *Lexer) stateFunc {
	v := l.currentToken()
	l.emitValue(String, v[1:len(v)-1])
	return initState
}

func ratioState(l *Lexer) stateFunc {
	v := big.NewRat(1, 1)
	v.SetString(l.currentToken())
	l.emitValue(Ratio, v)
	return initState
}

func numberState(l *Lexer) stateFunc {
	v := big.NewFloat(0)
	v.SetString(l.currentToken())
	l.emitValue(Number, v)
	return initState
}
