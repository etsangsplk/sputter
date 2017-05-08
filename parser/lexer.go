package parser

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

// Token is a Lexer token
type Token struct {
	Type  TokenType
	Value a.Value
}

func (t Token) Str() a.Str {
	return a.Str("")
}

type lexer struct {
	input  string
	start  int
	pos    int
	tokens chan a.Value
}

type stateFunc func(*lexer) stateFunc

type stateMapEntry struct {
	pattern  *regexp.Regexp
	function stateFunc
}

type stateMap []stateMapEntry

var (
	escaped    = regexp.MustCompile(`\\\\|\\"|\\[^\\"]`)
	escapedMap = map[string]string{
		`\\`: `\`,
		`\"`: `"`,
	}
	states stateMap
)

// NewLexer instantiates a new Lisp Lexer Sequence
func NewLexer(src string) a.Sequence {
	l := &lexer{
		input:  src,
		tokens: make(chan a.Value),
	}
	go l.run()

	return a.Filter(a.NewChannelSequence(l.tokens), func(v a.Value) bool {
		return isNotWhitespace(v.(Token))
	})
}

func isNotWhitespace(t Token) bool {
	return t.Type != Whitespace && t.Type != Comment
}

func (l *lexer) run() {
	for s := initState; s != nil; {
		s = s(l)
	}
	close(l.tokens)
}

func (l *lexer) emitValue(t TokenType, v a.Value) {
	l.tokens <- Token{
		Type:  t,
		Value: v,
	}
	l.skip()
}

func (l *lexer) emit(t TokenType) {
	l.emitValue(t, l.currentToken())
}

func (l *lexer) currentToken() a.Str {
	return a.Str(l.input[l.start:l.pos])
}

func (l *lexer) skip() {
	l.start = l.pos
}

func (l *lexer) matchState() stateFunc {
	src := l.input[l.pos:]
	for _, s := range states {
		if i := s.pattern.FindStringIndex(src); i != nil {
			r := src[:i[1]]
			l.pos += len(r)
			return s.function
		}
	}
	// Shouldn't happen because of the patterns that are defined,
	// but is here as a safety net
	panic(UnmatchedState)
}

func tokenState(t TokenType) stateFunc {
	return func(l *lexer) stateFunc {
		l.emit(t)
		return initState
	}
}

func initState(l *lexer) stateFunc {
	state := l.matchState()
	return state(l)
}

func endState(t TokenType) stateFunc {
	return func(l *lexer) stateFunc {
		l.emit(t)
		return nil
	}
}

func unescape(s a.Str) a.Str {
	r := escaped.ReplaceAllStringFunc(string(s), func(e string) string {
		if r, ok := escapedMap[e]; ok {
			return r
		}
		return e
	})
	return a.Str(r)
}

func stringState(l *lexer) stateFunc {
	v := l.currentToken()
	s := unescape(v[1 : len(v)-1])
	l.emitValue(String, a.Str(s))
	return initState
}

func ratioState(l *lexer) stateFunc {
	v := a.ParseNumber(l.currentToken())
	l.emitValue(Ratio, v)
	return initState
}

func numberState(l *lexer) stateFunc {
	v := a.ParseNumber(l.currentToken())
	l.emitValue(Number, v)
	return initState
}

func init() {
	pattern := func(p string, s stateFunc) stateMapEntry {
		return stateMapEntry{
			pattern:  regexp.MustCompile(p),
			function: s,
		}
	}

	states = stateMap{
		pattern(`^$`, endState(EndOfFile)),
		pattern(`^;[^\n]*[\n]`, tokenState(Comment)),
		pattern(`^[\s,]+`, tokenState(Whitespace)),
		pattern(`^\(`, tokenState(ListStart)),
		pattern(`^\[`, tokenState(VectorStart)),
		pattern(`^{`, tokenState(MapStart)),
		pattern(`^\)`, tokenState(ListEnd)),
		pattern(`^]`, tokenState(VectorEnd)),
		pattern(`^}`, tokenState(MapEnd)),
		pattern(`^'`, tokenState(QuoteMarker)),

		pattern(`^"(\\\\|\\"|\\[^\\"]|[^"\\])*"`, stringState),

		pattern(`^[+-]?[1-9]\d*/[1-9]\d*`, ratioState),
		pattern(`^[+-]?((0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?)`, numberState),

		pattern(`^[^(){}\[\]\s,'";]+`, tokenState(Identifier)),

		pattern(`^.`, endState(Error)),
	}
}
