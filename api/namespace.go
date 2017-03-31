package api

import "fmt"

const (
	defaultNamespaceEntries = 128
	defaultSymbolEntries    = 4096
)

const (
	// BuiltInDomain stores built-ins
	BuiltInDomain = Name("sputter")

	// UserDomain stores user defined vars
	UserDomain = Name("user")

	// LocalDomain stores local vars
	LocalDomain = Name("")

	// ContextDomain identifies the scoped domain
	ContextDomain = Name("*ns*")

	// ExpectedNamespace is thrown when a Value is not a Namespace
	ExpectedNamespace = "value is not a namespace: %s"
)

type symbolMap map[Name](Symbol)
type namespaceMap map[Name](Namespace)

var namespaces = make(namespaceMap, defaultNamespaceEntries)

// Namespace is a container where Qualified Symbols are mapped to Values
type Namespace interface {
	Context
	Domain() Name
	Intern(n Name) Symbol
}

type namespace struct {
	Context
	domain  Name
	symbols symbolMap
}

// Domain returns the Domain of the Namespace
func (ns *namespace) Domain() Name {
	return ns.domain
}

// Intern returns a Symbol based on the Name and Namespace Domain.
// This Symbol will be atomic, meaning that there will be only one
// instance, allowing the Symbols to be compared by reference
func (ns *namespace) Intern(n Name) Symbol {
	d := ns.domain
	k := qualifiedName(n, d)
	if s, ok := ns.symbols[k]; ok {
		return s
	}
	s := &symbol{name: n, domain: d}
	ns.symbols[k] = s
	return s
}

func (ns *namespace) String() string {
	return "(ns " + string(ns.domain) + ")"
}

// GetNamespace returns the Namespace for the specified domain.
func GetNamespace(n Name) Namespace {
	if ns, ok := namespaces[n]; ok {
		return ns
	}
	ns := &namespace{
		Context: NewContext(),
		domain:  n,
		symbols: make(symbolMap, defaultSymbolEntries),
	}
	namespaces[n] = ns
	return ns
}

// GetContextNamespace resolves the Namespace based on its Context
func GetContextNamespace(c Context) Namespace {
	if v, ok := c.Get(ContextDomain); ok {
		return AssertNamespace(v)
	}
	return GetNamespace(UserDomain)
}

// AssertNamespace will cast a Value to a Namespace or explode violently
func AssertNamespace(v Value) Namespace {
	if r, ok := v.(Namespace); ok {
		return r
	}
	panic(fmt.Sprintf(ExpectedNamespace, String(v)))
}

func init() {
	builtInContext := NewContext()
	userContext := ChildContext(builtInContext)

	namespaces[BuiltInDomain] = &namespace{
		Context: builtInContext,
		domain:  BuiltInDomain,
		symbols: make(symbolMap, defaultSymbolEntries),
	}

	namespaces[UserDomain] = &namespace{
		Context: userContext,
		domain:  UserDomain,
		symbols: make(symbolMap, defaultSymbolEntries),
	}
}
