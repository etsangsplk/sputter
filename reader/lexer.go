package reader

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

// EOFToken marks the end of a Lexer stream
var EOFToken = &Token{EndOfFile, ""}

// Token is a Lexer token
type Token struct {
	Type  TokenType
	Value a.Value
}

// Lexer is the lexer interface
type Lexer interface {
	Next() *Token
}

type lispLexer struct {
	input  string
	start  int
	pos    int
	tokens chan *Token
}

type stateFunc func(*lispLexer) stateFunc

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

// NewLexer instantiates a new Lisp Lexer instance
func NewLexer(src string) Lexer {
	l := &lispLexer{
		input:  src,
		tokens: make(chan *Token),
	}

	go l.run()
	return l
}

// Next returns the next Token from the lexer's Token channel
func (l *lispLexer) Next() *Token {
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

func (l *lispLexer) run() {
	for s := initState; s != nil; {
		s = s(l)
	}
	close(l.tokens)
}

func (l *lispLexer) emitValue(t TokenType, v a.Value) {
	l.tokens <- &Token{t, v}
	l.skip()
}

func (l *lispLexer) emit(t TokenType) {
	l.emitValue(t, l.currentToken())
}

func (l *lispLexer) currentToken() string {
	return l.input[l.start:l.pos]
}

func (l *lispLexer) skip() {
	l.start = l.pos
}

func (l *lispLexer) matchState() stateFunc {
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
	return func(l *lispLexer) stateFunc {
		l.emit(t)
		return initState
	}
}

func initState(l *lispLexer) stateFunc {
	var state = l.matchState()
	return state(l)
}

func endState(t TokenType) stateFunc {
	return func(l *lispLexer) stateFunc {
		l.emit(t)
		return nil
	}
}

func stringState(l *lispLexer) stateFunc {
	v := l.currentToken()
	l.emitValue(String, v[1:len(v)-1])
	return initState
}

func ratioState(l *lispLexer) stateFunc {
	v := big.NewRat(1, 1)
	v.SetString(l.currentToken())
	l.emitValue(Ratio, v)
	return initState
}

func numberState(l *lispLexer) stateFunc {
	v := big.NewFloat(0)
	v.SetString(l.currentToken())
	l.emitValue(Number, v)
	return initState
}
