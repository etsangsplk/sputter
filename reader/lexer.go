package interpreter

import (
	"math/big"
	"regexp"

	a "github.com/kode4food/sputter/api"
)

const (
	// UnexpectedEndOfFile is the error returned when EOF is unexpectedly reached
	UnexpectedEndOfFile = "end of file reached unexpectedly"

	// UnmatchedState is the error returned when the lexing state is invalid
	UnmatchedState = "unmatched lexing state"
)

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
	VectorStart
	VectorEnd
	DataMarker
	EndOfFile
	Whitespace
	Comment
)

// EOFToken marks the end of a Reader stream
var EOFToken = &Token{EndOfFile, ""}

// Token is a Lexer token
type Token struct {
	Type  TokenType
	Value a.Value
}

// Reader defines the one method that a Token processor must provide
type Reader interface {
	Next() *Token
}

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
		{regexp.MustCompile(`^;[^\n]*[\n]`), tokenState(Comment)},
		{regexp.MustCompile(`^\s+`), tokenState(Whitespace)},
		{regexp.MustCompile(`^\(`), tokenState(ListStart)},
		{regexp.MustCompile(`^\[`), tokenState(VectorStart)},
		{regexp.MustCompile(`^\)`), tokenState(ListEnd)},
		{regexp.MustCompile(`^]`), tokenState(VectorEnd)},
		{regexp.MustCompile(`^'`), tokenState(DataMarker)},

		{regexp.MustCompile(`^"(\\.|[^"])*"`), stringState},
		{regexp.MustCompile(`^[+-]?[1-9]\d*/[1-9]\d*`), ratioState},
		{regexp.MustCompile(`^[+-]?(0|[1-9]\d*(\.\d+)?([eE][+-]?\d+)?)`), numberState},

		{regexp.MustCompile(`^[^()\[\]\s]+`), tokenState(Identifier)},
		{regexp.MustCompile(`^.`), endState(Error)},
	}
}

// NewLexer instantiates a new Lexer instance
func NewLexer(src string) *Lexer {
	l := &Lexer{
		input:  src,
		tokens: make(chan *Token),
	}

	go l.run()
	return l
}

// Next returns the next Token from the lexer's Token channel
func (l *Lexer) Next() *Token {
	for {
		t, ok := <-l.tokens
		if !ok {
			panic(UnexpectedEndOfFile)
		}
		if t.Type != Whitespace {
			return t
		}
	}
}

func (l *Lexer) run() {
	for s := initState; s != nil; {
		s = s(l)
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
	src := l.input[l.pos:]
	for _, s := range states {
		if i := s.pattern.FindStringIndex(src); i != nil {
			r := src[i[0]:i[1]]
			l.pos += len(r)
			return s.function
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
