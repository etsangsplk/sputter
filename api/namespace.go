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
)

type (
	// Namespace is a container where Qualified Symbols are mapped to Values
	Namespace interface {
		Context
		Domain() Name
		Intern(Name) Symbol
		NamespaceType()
	}

	qualifiedNamespace struct {
		Context
		domain  Name
		symbols u.Cache
	}

	localNamespace struct {
		Context
		domain  Name
		symbols u.Cache
	}

	withNamespace struct {
		Context
		ns Namespace
	}
)

var namespaces = u.NewCache()

func (ns *qualifiedNamespace) Domain() Name {
	return ns.domain
}

func (ns *qualifiedNamespace) Intern(n Name) Symbol {
	d := ns.domain
	k := Name(d + ":" + n)
	return ns.symbols.Get(k, func() u.Any {
		return &qualifiedSymbol{
			name:   n,
			domain: d,
		}
	}).(Symbol)
}

func (ns *qualifiedNamespace) Str() Str {
	return "(ns " + Str(ns.domain) + ")"
}

func (ns *qualifiedNamespace) NamespaceType() {}

func (ns *localNamespace) Domain() Name {
	return LocalDomain
}

func (ns *localNamespace) Intern(n Name) Symbol {
	return ns.symbols.Get(n, func() u.Any {
		return &localSymbol{
			name: n,
		}
	}).(Symbol)
}

func (ns *localNamespace) Str() Str {
	return "(ns *local*)"
}

func (ns *localNamespace) NamespaceType() {}

// GetNamespace returns the Namespace for the specified domain.
func GetNamespace(n Name) Namespace {
	return namespaces.Get(n, func() u.Any {
		ns := &qualifiedNamespace{
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
		return v.(Namespace)
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

// Get retrieves a value from the Context chain
func (w *withNamespace) Get(n Name) (Value, bool) {
	if v, ok := w.ns.Get(n); ok {
		return v, true
	}
	return w.Context.Get(n)
}

func init() {
	builtInContext := NewContext()
	userContext := ChildContext(builtInContext)

	namespaces.Get(LocalDomain, func() u.Any {
		ns := &localNamespace{
			Context: NewContext(),
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	})

	namespaces.Get(BuiltInDomain, func() u.Any {
		ns := &qualifiedNamespace{
			Context: builtInContext,
			domain:  BuiltInDomain,
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	})

	namespaces.Get(UserDomain, func() u.Any {
		ns := &qualifiedNamespace{
			Context: userContext,
			domain:  UserDomain,
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	})
}
