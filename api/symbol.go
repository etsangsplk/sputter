package api

import (
	"fmt"
	"strings"
)

const (
	// UnknownSymbol is thrown if a symbol cannot be resolved
	UnknownSymbol = "symbol has not been defined: %s"

	// ExpectedSymbol is thrown when a Value is not a unqualified Symbol
	ExpectedSymbol = "value is not a symbol: %s"

	// ExpectedUnqualified is thrown when a Symbol is unexpectedly qualified
	ExpectedUnqualified = "symbol should be unqualified: %s"
)

// Symbol is a qualified identifier that can be resolved
type Symbol interface {
	Evaluable
	Named
	Domain() Name
	Qualified() Name
	Namespace(Context) Namespace
	Resolve(Context) (Value, bool)
}

type symbol struct {
	name   Name
	domain Name
}

func qualifiedName(name Name, domain Name) Name {
	if domain == LocalDomain {
		return name
	}
	return Name(domain + ":" + name)
}

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Name, domain Name) Symbol {
	ns := GetNamespace(domain)
	return ns.Intern(name)
}

// NewLocalSymbol returns a Symbol from the local domain
func NewLocalSymbol(name Name) Symbol {
	return NewQualifiedSymbol(name, LocalDomain)
}

// ParseSymbol parses a qualified Name and produces a Symbol
func ParseSymbol(n Name) Symbol {
	if i := strings.IndexRune(string(n), ':'); i != -1 {
		return NewQualifiedSymbol(n[i+1:], n[:i])
	}
	return NewLocalSymbol(n)
}

// Name returns the Name portion of the Symbol
func (s symbol) Name() Name {
	return s.name
}

// Domain returns the domain portion of the Symbol
func (s symbol) Domain() Name {
	return s.domain
}

// Qualified returns the fully-qualified Name of a Symbol
func (s symbol) Qualified() Name {
	return qualifiedName(s.name, s.domain)
}

// Namespace returns the Namespace for this Symbol
func (s symbol) Namespace(c Context) Namespace {
	d := s.domain
	if d != LocalDomain {
		return GetNamespace(d)
	}
	return GetContextNamespace(c)
}

// Resolve a Symbol against a Context
func (s symbol) Resolve(c Context) (Value, bool) {
	n := s.name
	d := s.domain
	if d != LocalDomain {
		return GetNamespace(d).Get(n)
	}
	if r, ok := c.Get(s.name); ok {
		return r, true
	}
	return GetContextNamespace(c).Get(n)
}

// Eval makes a Symbol Evaluable
func (s symbol) Eval(c Context) Value {
	if r, ok := s.Resolve(c); ok {
		return r
	}
	panic(fmt.Sprintf(UnknownSymbol, s.name))
}

func (s symbol) String() string {
	return string(s.Qualified())
}

// AssertUnqualified will cast a Value into a Symbol and explode
// violently if it's qualified with a domain
func AssertUnqualified(v Value) Symbol {
	if r, ok := v.(Symbol); ok {
		if r.Domain() == LocalDomain {
			return r
		}
		panic(fmt.Sprintf(ExpectedUnqualified, r.Qualified()))
	}
	panic(fmt.Sprintf(ExpectedSymbol, String(v)))
}
