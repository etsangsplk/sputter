package api

import "strings"

const (
	// UnknownSymbol is thrown if a symbol cannot be resolved
	UnknownSymbol = "symbol has not been defined: %s"

	// ExpectedSymbol is thrown when a Value is not a unqualified Symbol
	ExpectedSymbol = "value is not a symbol: %s"

	// ExpectedUnqualified is thrown when a Symbol is unexpectedly qualified
	ExpectedUnqualified = "symbol should be unqualified: %s"
)

type (
	// Symbol is a qualified identifier that can be resolved
	Symbol interface {
		Value
		Evaluable
		IsSymbol() bool
		Name() Name
		Domain() Name
		Qualified() Name
		Namespace(Context) Namespace
		Resolve(Context) (Value, bool)
	}

	symbol struct {
		name   Name
		domain Name
	}
)

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

// NewBuiltInSymbol returns a Symbol qualified to the BuiltIn domain
func NewBuiltInSymbol(name Name) Symbol {
	return NewQualifiedSymbol(name, BuiltInDomain)
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

func (s *symbol) IsSymbol() bool {
	return true
}

func (s *symbol) Name() Name {
	return s.name
}

func (s *symbol) Domain() Name {
	return s.domain
}

func (s *symbol) Qualified() Name {
	return qualifiedName(s.name, s.domain)
}

func (s *symbol) Namespace(c Context) Namespace {
	d := s.domain
	if d != LocalDomain {
		return GetNamespace(d)
	}
	return GetContextNamespace(c)
}

func resolveNamespace(_ Context, domain Name) Namespace {
	return GetNamespace(domain)
}

func (s *symbol) Resolve(c Context) (Value, bool) {
	d := s.domain
	if d != LocalDomain {
		return resolveNamespace(c, d).Get(s.name)
	}
	n := s.name
	if r, ok := c.Get(n); ok {
		return r, true
	}
	return GetContextNamespace(c).Get(n)
}

func (s *symbol) Eval(c Context) Value {
	if r, ok := s.Resolve(c); ok {
		return r
	}
	panic(ErrStr(UnknownSymbol, s.Qualified()))
}

func (s *symbol) Str() Str {
	return Str(s.Qualified())
}

// AssertUnqualified will cast a Value into a Symbol and explode
// violently if it's qualified with a domain
func AssertUnqualified(v Value) Symbol {
	if r, ok := v.(Symbol); ok {
		if r.Domain() == LocalDomain {
			return r
		}
		panic(ErrStr(ExpectedUnqualified, r.Qualified()))
	}
	panic(ErrStr(ExpectedSymbol, v))
}
