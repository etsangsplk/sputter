package api

import u "github.com/kode4food/sputter/util"

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

var namespaces = u.NewCache()

// Namespace is a container where Qualified Symbols are mapped to Values
type Namespace interface {
	Value
	Context
	Domain() Name
	Intern(Name) Symbol
}

type namespace struct {
	Context
	domain  Name
	symbols u.Cache
}

type withNamespace struct {
	Context
	ns Namespace
}

func (ns *namespace) Domain() Name {
	return ns.domain
}

// Intern returns a IsSymbol based on the Name and Namespace Domain.
// This IsSymbol will be atomic, meaning that there will be only one
// instance, allowing the Symbols to be compared by reference
func (ns *namespace) Intern(n Name) Symbol {
	d := ns.domain
	k := qualifiedName(n, d)
	return ns.symbols.Get(k, func() u.Any {
		return &symbol{name: n, domain: d}
	}).(Symbol)
}

func (ns *namespace) Str() Str {
	return "(ns " + Str(ns.domain) + ")"
}

// GetNamespace returns the Namespace for the specified domain.
func GetNamespace(n Name) Namespace {
	return namespaces.Get(n, func() u.Any {
		ns := &namespace{
			Context: NewContext(),
			domain:  n,
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	}).(Namespace)
}

// GetContextNamespace resolves the Namespace based on its Context
func GetContextNamespace(c Context) Namespace {
	if v, ok := c.Get(ContextDomain); ok {
		return AssertNamespace(v)
	}
	return GetNamespace(UserDomain)
}

// WithNamespace creates a child Context that performs a Namespace lookup
// before checking the Context's parent
func WithNamespace(c Context, ns Namespace) Context {
	return ChildContext(&withNamespace{
		Context: c,
		ns:      ns,
	})
}

// ElementAt retrieves a value from the Context chain
func (w *withNamespace) Get(n Name) (v Value, bound bool) {
	if v, ok := w.ns.Get(n); ok {
		return v, true
	}
	return w.Context.Get(n)
}

// AssertNamespace will cast a Value to a Namespace or explode violently
func AssertNamespace(v Value) Namespace {
	if r, ok := v.(Namespace); ok {
		return r
	}
	panic(Err(ExpectedNamespace, v))
}

func init() {
	builtInContext := NewContext()
	userContext := ChildContext(builtInContext)

	namespaces.Get(BuiltInDomain, func() u.Any {
		ns := &namespace{
			Context: builtInContext,
			domain:  BuiltInDomain,
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	})

	namespaces.Get(UserDomain, func() u.Any {
		ns := &namespace{
			Context: userContext,
			domain:  UserDomain,
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	})
}
