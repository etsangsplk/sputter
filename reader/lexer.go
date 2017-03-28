package reader

import (
	"regexp"

	a "github.com/kode4food/sputter/api"
)

const (
	// UnexpectedEndOfFile is the error returned when EOF is unexpectedly reached
	UnexpectedEndOfFile = "end of file reached unexpectedly"

	// UnmatchedState is the error returned when the lexer state is invalid
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
	MapStart
	MapEnd
	QuoteMarker
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
	re := regexp.MustCompile
	states = stateMap{
		{re(`^$`), endState(EndOfFile)},
		{re(`^;[^\n]*[\n]`), tokenState(Comment)},
		{re(`^[\s,]+`), tokenState(Whitespace)},
		{re(`^\(`), tokenState(ListStart)},
		{re(`^\[`), tokenState(VectorStart)},
		{re(`^{`), tokenState(MapStart)},
		{re(`^\)`), tokenState(ListEnd)},
		{re(`^]`), tokenState(VectorEnd)},
		{re(`^}`), tokenState(MapEnd)},
		{re(`^'`), tokenState(QuoteMarker)},

		{re(`^"(\\.|[^"])*"`), stringState},

		{re(`^[+-]?[1-9]\d*/[1-9]\d*`), ratioState},
		{re(`^[+-]?((0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?)`), numberState},

		{re(`^[^(){}\[\]\s,'";]+`), tokenState(Identifier)},

		{re(`^.`), endState(Error)},
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

func isNotWhitespace(t *Token) bool {
	return t.Type != Whitespace && t.Type != Comment
}

// Next returns the next Token from the lexer's Token channel
func (l *lispLexer) Next() *Token {
	for {
		t, ok := <-l.tokens
		if !ok {
			panic(UnexpectedEndOfFile)
		}
		if isNotWhitespace(t) {
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
			r := src[:i[1]]
			l.pos += len(r)
			return s.function
		}
	}
	// Can't happen because of the patterns that are defined,
	// but is here as a safety net
	panic(UnmatchedState)
}

func tokenState(t TokenType) stateFunc {
	return func(l *lispLexer) stateFunc {
		l.emit(t)
		return initState
	}
}

func initState(l *lispLexer) stateFunc {
	state := l.matchState()
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
	v := a.ParseNumber(l.currentToken())
	l.emitValue(Ratio, v)
	return initState
}

func numberState(l *lispLexer) stateFunc {
	v := a.ParseNumber(l.currentToken())
	l.emitValue(Number, v)
	return initState
}
