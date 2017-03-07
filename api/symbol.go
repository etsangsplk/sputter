package api

import "strings"

const (
	// UnknownSymbol is thrown if a symbol cannot be resolved
	UnknownSymbol = "symbol has not been defined"

	// BadQualifiedName is thrown when a name has too many components
	BadQualifiedName = "expected a valid qualified name"

	// ExpectedUnqualified is thrown when a Value is not a unqualified Symbol
	ExpectedUnqualified = "value is not a symbol"
)

var keywords = make(Variables, 4096)

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

// NewQualifiedSymbol returns a Qualified Symbol for a specific domain
func NewQualifiedSymbol(name Name, domain Name) *Symbol {
	ns := GetNamespace(domain)
	return ns.Intern(name)
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
		return NewLocalSymbol(Name(s[0]))
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
	panic(UnknownSymbol)
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

// Keyword is an Atom that represents a Name for mapping purposes
type Keyword struct {
	name Name
}

// NewKeyword returns an interned instance of a Keyword
func NewKeyword(n Name) *Keyword {
	if r, ok := keywords[n]; ok {
		return r.(*Keyword)
	}
	r := &Keyword{name: n}
	keywords[n] = r
	return r
}

// Eval makes Keyword Evaluable
func (k *Keyword) Eval(c Context) Value {
	return k
}

func (k *Keyword) String() string {
	return ":" + string(k.name)
}

// AssertUnqualified will cast a Value into a Symbol and explode
// violently if it's qualified with a domain
func AssertUnqualified(v Value) *Symbol {
	if r, ok := v.(*Symbol); ok {
		if r.Domain == LocalDomain {
			return r
		}
	}
	panic(ExpectedUnqualified)
}
