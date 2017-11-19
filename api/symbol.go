package api

import (
	"fmt"
	"strings"
	"sync/atomic"
)

const (
	// UnknownSymbol is thrown if a symbol cannot be resolved
	UnknownSymbol = "symbol has not been defined: %s"

	genSymTemplate = "x-%s-gensym-%d"
)

type (
	// Symbol is a qualified identifier that can be resolved
	Symbol interface {
		Value
		Evaluable
		Name() Name
		Domain() Name
		Qualified() Name
		Namespace(Context) Namespace
		Resolve(Context) (Value, bool)
		SymbolType()
	}

	// QualifiedSymbol represents a domain-qualified symbol
	QualifiedSymbol interface {
		Symbol
		QualifiedSymbolType()
	}

	// LocalSymbol represents a locally-resolved symbol (unqualified)
	LocalSymbol interface {
		Symbol
		LocalSymbolType()
	}

	localSymbol struct {
		name Name
	}

	qualifiedSymbol struct {
		name   Name
		domain Name
	}
)

var genSymIncrement uint64

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Name, domain Name) Symbol {
	ns := GetNamespace(domain)
	return ns.Intern(name)
}

// NewBuiltInSymbol returns a Symbol qualified to the BuiltIn domain
func NewBuiltInSymbol(name Name) Symbol {
	return NewQualifiedSymbol(name, BuiltInDomain)
}

// NewLocalSymbol returns a Symbol from the local domain
func NewLocalSymbol(name Name) Symbol {
	return NewQualifiedSymbol(name, LocalDomain)
}

// NewGeneratedSymbol creates a generated Symbol
func NewGeneratedSymbol(name Name) Symbol {
	idx := atomic.AddUint64(&genSymIncrement, 1)
	q := fmt.Sprintf(genSymTemplate, name, idx)
	return NewLocalSymbol(Name(q))
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(n Name) Symbol {
	if i := strings.IndexRune(string(n), ':'); i != -1 {
		return NewQualifiedSymbol(n[i+1:], n[:i])
	}
	return NewLocalSymbol(n)
}

func resolveNamespace(_ Context, domain Name) Namespace {
	return GetNamespace(domain)
}

func (s *qualifiedSymbol) Name() Name {
	return s.name
}

func (s *qualifiedSymbol) Domain() Name {
	return s.domain
}

func (s *qualifiedSymbol) Qualified() Name {
	return Name(s.domain + ":" + s.name)
}

func (s *qualifiedSymbol) Namespace(c Context) Namespace {
	return GetNamespace(s.domain)
}

func (s *qualifiedSymbol) Resolve(c Context) (Value, bool) {
	return resolveNamespace(c, s.domain).Get(s.name)
}

func (s *qualifiedSymbol) Eval(c Context) Value {
	if r, ok := s.Resolve(c); ok {
		return r
	}
	panic(ErrStr(UnknownSymbol, s.Qualified()))
}

func (s *qualifiedSymbol) Str() Str {
	return Str(s.Qualified())
}

func (s *qualifiedSymbol) SymbolType()          {}
func (s *qualifiedSymbol) QualifiedSymbolType() {}

func (s *localSymbol) Name() Name {
	return s.name
}

func (s *localSymbol) Domain() Name {
	return LocalDomain
}

func (s *localSymbol) Qualified() Name {
	return s.name
}

func (s *localSymbol) Namespace(c Context) Namespace {
	return GetContextNamespace(c)
}

func (s *localSymbol) Resolve(c Context) (Value, bool) {
	n := s.name
	if r, ok := c.Get(n); ok {
		return r, true
	}
	return GetContextNamespace(c).Get(n)
}

func (s *localSymbol) Eval(c Context) Value {
	if r, ok := s.Resolve(c); ok {
		return r
	}
	panic(ErrStr(UnknownSymbol, s.Qualified()))
}

func (s *localSymbol) Str() Str {
	return Str(s.Qualified())
}

func (s *localSymbol) SymbolType()      {}
func (s *localSymbol) LocalSymbolType() {}
