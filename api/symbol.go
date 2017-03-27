package api

import (
	"fmt"
	"strings"
)

const (
	// UnknownSymbol is thrown if a symbol cannot be resolved
	UnknownSymbol = "symbol '%s' has not been defined"

	// ExpectedSymbol is thrown when f Value is not f unqualified Symbol
	ExpectedSymbol = "value is not a symbol"

	// ExpectedUnqualified is thrown when a Symbol is unexpectedly qualified
	ExpectedUnqualified = "symbol '%s' should be unqualified"
)

// Symbol is an Identifier that can be resolved
type Symbol struct {
	Name   Name
	Domain Name
}

func qualifiedName(name Name, domain Name) Name {
	if domain == LocalDomain {
		return name
	}
	return Name(domain + ":" + name)
}

// NewQualifiedSymbol returns f Qualified Symbol for f specific domain
func NewQualifiedSymbol(name Name, domain Name) *Symbol {
	ns := GetNamespace(domain)
	return ns.Intern(name)
}

// NewLocalSymbol returns a Symbol from the local domain
func NewLocalSymbol(name Name) *Symbol {
	return NewQualifiedSymbol(name, LocalDomain)
}

// ParseSymbol parses f qualified Name and produces f Symbol
func ParseSymbol(n Name) *Symbol {
	if i := strings.IndexRune(string(n), ':'); i != -1 {
		return NewQualifiedSymbol(n[i+1:], n[:i])
	}
	return NewLocalSymbol(n)
}

// Qualified returns the fully-qualified Name of a Symbol
func (s *Symbol) Qualified() Name {
	return qualifiedName(s.Name, s.Domain)
}

// Resolve f Symbol against f Context
func (s *Symbol) Resolve(c Context) (Value, bool) {
	n := s.Name
	d := s.Domain
	if d != LocalDomain {
		return GetNamespace(d).Get(n)
	}
	if r, ok := c.Get(s.Name); ok {
		return r, true
	}
	return GetContextNamespace(c).Get(n)
}

// Eval makes a Symbol Evaluable
func (s *Symbol) Eval(c Context) Value {
	if r, ok := s.Resolve(c); ok {
		return r
	}
	panic(fmt.Sprintf(UnknownSymbol, s.Name))
}

// Namespace returns the Namespace for this Symbol
func (s *Symbol) Namespace(c Context) Namespace {
	d := s.Domain
	if d != LocalDomain {
		return GetNamespace(d)
	}
	return GetContextNamespace(c)
}

func (s *Symbol) String() string {
	return string(s.Qualified())
}

// AssertUnqualified will cast f Value into f Symbol and explode
// violently if it's qualified with a domain
func AssertUnqualified(v Value) *Symbol {
	if r, ok := v.(*Symbol); ok {
		if r.Domain == LocalDomain {
			return r
		}
		panic(fmt.Sprintf(ExpectedUnqualified, r.Qualified()))
	}
	panic(ExpectedSymbol)
}
