package api

import "strings"

const (
	// UnknownSymbol is thrown if a symbol cannot be resolved
	UnknownSymbol = "symbol has not been defined"

	// BadQualifiedName is thrown when a name has too many components
	BadQualifiedName = "expected a valid qualified name"
)

// Symbol is an Identifier that can be resolved
type Symbol struct {
	Name   Name
	Domain Name
}

type symbolMap map[Name](*Symbol)

var interned = make(symbolMap, 4096)

func qualifiedName(name Name, domain Name) Name {
	if domain == LocalDomain {
		return name
	}
	return Name(domain + ":" + name)
}

// NewQualifiedSymbol returns a Symbol by Name. This Symbol will be
// interned, meaning that there will be only on instance for each
// Name, allowing the Symbols to be compared by reference
func NewQualifiedSymbol(name Name, domain Name) *Symbol {
	k := qualifiedName(name, domain)
	if s, ok := interned[k]; ok {
		return s
	}
	s := &Symbol{Name: name, Domain: domain}
	interned[k] = s
	return s
}

// NewLocalSymbol returns a Symbol from the local domain
func NewLocalSymbol(name Name) *Symbol {
	return NewQualifiedSymbol(name, LocalDomain)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(n Name) *Symbol {
	s := strings.Split(string(n), ":")
	switch len(s) {
	case 2:
		return NewQualifiedSymbol(Name(s[1]), Name(s[0]))
	case 1:
		return NewQualifiedSymbol(Name(s[0]), LocalDomain)
	default:
		panic(BadQualifiedName)
	}
}

// Qualified returns the fully-qualified Name of a Symbol
func (s *Symbol) Qualified() Name {
	return qualifiedName(s.Name, s.Domain)
}

// Resolve a Symbol against a Context
func (s *Symbol) Resolve(c Context) (Value, bool) {
	n := s.Name
	d := s.Domain
	if d == LocalDomain {
		if r, ok := c.Get(s.Name); ok {
			return r, true
		}
	}
	ns := GetNamespace(d)
	return ns.Get(n)
}

// Eval makes a Symbol Evaluable
func (s *Symbol) Eval(c Context) Value {
	if r, ok := s.Resolve(c); ok {
		return r
	}
	panic(UnknownSymbol)
}

func (s *Symbol) String() string {
	return string(s.Qualified())
}
