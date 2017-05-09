package parser

import (
	"regexp"

	a "github.com/kode4food/sputter/api"
)

const (
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

// Lexer breaks Lisp expressions into Tokens
type Lexer interface {
	a.Sequence
}

type lexer struct {
	resolved bool
	src      string
	isSeq    bool
	first    a.Value
	rest     *lexer
}

// Token is a lexer token
type Token struct {
	Type  TokenType
	Value a.Value
}

type tokenizer func(s a.Str) *Token

type matchEntry struct {
	pattern  *regexp.Regexp
	function tokenizer
}

type mactchEntries []matchEntry

var (
	escaped    = regexp.MustCompile(`\\\\|\\"|\\[^\\"]`)
	escapedMap = map[string]string{
		`\\`: `\`,
		`\"`: `"`,
	}
	matchers mactchEntries
)

// NewLexer creates a new lexer instance
func NewLexer(src string) a.Sequence {
	l := &lexer{
		src: src,
	}

	return a.Filter(l, func(v a.Value) bool {
		return isNotWhitespace(v.(*Token))
	})
}

func (l *lexer) resolve() *lexer {
	if l.resolved {
		return l
	}

	t, src := l.matchToken()

	l.resolved = true
	l.isSeq = true
	l.first = t
	l.rest = &lexer{
		src: src,
	}
	return l
}

func (l *lexer) matchToken() (*Token, string) {
	src := l.src
	for _, s := range matchers {
		if i := s.pattern.FindStringIndex(src); i != nil {
			f := src[:i[1]]
			r := src[len(f):]
			return s.function(a.Str(f)), r
		}
	}
	// Shouldn't happen because of the patterns that are defined,
	// but is here as a safety net
	panic(UnmatchedState)
}

func (l *lexer) IsSequence() bool {
	return l.resolve().isSeq
}

func (l *lexer) First() a.Value {
	return l.resolve().first
}

func (l *lexer) Rest() a.Sequence {
	return l.resolve().rest
}

func (l *lexer) Prepend(v a.Value) a.Sequence {
	return &lexer{
		isSeq: true,
		first: v,
		rest:  l,
	}
}

func (l *lexer) Str() a.Str {
	return a.MakeSequenceStr(l)
}

// Str converts this Value into a Str
func (t *Token) Str() a.Str {
	return a.Str("")
}

func isNotWhitespace(t *Token) bool {
	return t.Type != Whitespace && t.Type != Comment
}

func makeToken(t TokenType, v a.Value) *Token {
	return &Token{
		Type:  t,
		Value: v,
	}
}

func tokenState(t TokenType) tokenizer {
	return func(s a.Str) *Token {
		return makeToken(t, a.Str(s))
	}
}

func endState(t TokenType) tokenizer {
	return func(_ a.Str) *Token {
		return makeToken(t, a.Nil)
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

func stringState(r a.Str) *Token {
	s := unescape(r[1 : len(r)-1])
	return makeToken(String, a.Str(s))
}

func ratioState(s a.Str) *Token {
	v := a.ParseNumber(s)
	return makeToken(Ratio, v)
}

func numberState(s a.Str) *Token {
	v := a.ParseNumber(s)
	return makeToken(Number, v)
}

func init() {
	pattern := func(p string, s tokenizer) matchEntry {
		return matchEntry{
			pattern:  regexp.MustCompile(p),
			function: s,
		}
	}

	matchers = mactchEntries{
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
