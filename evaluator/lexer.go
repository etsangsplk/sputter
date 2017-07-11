package evaluator

import (
	"regexp"

	a "github.com/kode4food/sputter/api"
)

// UnmatchedState is the error returned when the lexer state is invalid
const UnmatchedState = "unmatched lexing state"

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
	SyntaxMarker
	UnquoteMarker
	SpliceMarker
	Whitespace
	Comment
	endOfFile
)

type (
	// TokenType is an opaque type for lexer tokens
	TokenType int

	// Token is a lexer value
	Token struct {
		Type  TokenType
		Value a.Value
	}

	tokenizer func(s string) *Token

	matchEntry struct {
		pattern  *regexp.Regexp
		function tokenizer
	}

	matchEntries []matchEntry
)

var (
	escaped    = regexp.MustCompile(`\\\\|\\"|\\[^\\"]`)
	escapedMap = map[string]string{
		`\\`: `\`,
		`\"`: `"`,
	}
	matchers matchEntries
)

// Scan creates a new lexer Sequence
func Scan(src a.Str) a.Sequence {
	var resolver a.Resolver
	s := string(src)

	resolver = func() (a.Value, bool, a.Resolver) {
		if t, rs := matchToken(s); t.Type != endOfFile {
			s = rs
			return t, true, resolver
		}
		return nil, false, nil
	}

	l := a.NewLazySequence(resolver)

	return a.Filter(l, func(v a.Value) bool {
		return isNotWhitespace(v.(*Token))
	})
}

func matchToken(src string) (*Token, string) {
	for _, s := range matchers {
		if i := s.pattern.FindStringIndex(src); i != nil {
			f := src[:i[1]]
			r := src[len(f):]
			return s.function(f), r
		}
	}
	// Shouldn't happen because of the patterns that are defined,
	// but is here as a safety net
	panic(UnmatchedState)
}

// Str converts this Value into a Str
func (t *Token) Str() a.Str {
	return a.NewVector(a.NewFloat(float64(t.Type)), t.Value).Str()
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
	return func(s string) *Token {
		return makeToken(t, a.Str(s))
	}
}

func endState(t TokenType) tokenizer {
	return func(_ string) *Token {
		return makeToken(t, a.Nil)
	}
}

func unescape(s string) string {
	r := escaped.ReplaceAllStringFunc(s, func(e string) string {
		if r, ok := escapedMap[e]; ok {
			return r
		}
		return e
	})
	return r
}

func stringState(r string) *Token {
	s := unescape(r[1 : len(r)-1])
	return makeToken(String, a.Str(s))
}

func ratioState(s string) *Token {
	v := a.ParseNumber(a.Str(s))
	return makeToken(Ratio, v)
}

func numberState(s string) *Token {
	v := a.ParseNumber(a.Str(s))
	return makeToken(Number, v)
}

func init() {
	pattern := func(p string, s tokenizer) matchEntry {
		return matchEntry{
			pattern:  regexp.MustCompile(p),
			function: s,
		}
	}

	matchers = matchEntries{
		pattern(`^$`, endState(endOfFile)),
		pattern(`^;[^\n]*([\n]|$)`, tokenState(Comment)),
		pattern(`^[\s,]+`, tokenState(Whitespace)),
		pattern(`^\(`, tokenState(ListStart)),
		pattern(`^\[`, tokenState(VectorStart)),
		pattern(`^{`, tokenState(MapStart)),
		pattern(`^\)`, tokenState(ListEnd)),
		pattern(`^]`, tokenState(VectorEnd)),
		pattern(`^}`, tokenState(MapEnd)),
		pattern(`^'`, tokenState(QuoteMarker)),
		pattern("^`", tokenState(SyntaxMarker)),
		pattern(`^~@`, tokenState(SpliceMarker)),
		pattern(`^~`, tokenState(UnquoteMarker)),

		pattern(`^"(\\\\|\\"|\\[^\\"]|[^"\\])*"`, stringState),

		pattern(`^[+-]?[1-9]\d*/[1-9]\d*`, ratioState),
		pattern(`^[+-]?((0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?)`, numberState),

		pattern(`^[^(){}\[\]\s,'~@";]+`, tokenState(Identifier)),

		pattern(`^.`, endState(Error)),
	}
}
