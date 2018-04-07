package api

import (
	"sync"

	u "github.com/kode4food/sputter/util"
)

const (
	// BuiltInDomain stores built-ins
	BuiltInDomain = Name("sputter")

	// UserDomain stores user defined vars
	UserDomain = Name("user")

	// LocalDomain interns local vars
	LocalDomain = Name("")

	// ContextDomain identifies the scoped domain
	ContextDomain = Name("*ns*")
)

type (
	// Namespace is a container where Qualified Symbols are mapped to Vector
	Namespace interface {
		Context
		Domain() Name
		Intern(Name) Symbol
		NamespaceType()
	}

	namespace struct {
		Context
		sync.RWMutex
		symbols u.Cache
	}

	qualifiedNamespace struct {
		namespace
		domain Name
	}

	childNamespace struct {
		Namespace
		parent Namespace
	}

	withNamespace struct {
		Context
		ns Namespace
	}
)

var namespaces = u.NewCache()

// Get retrieves a variable by name
func (ns *namespace) Get(n Name) (Value, bool) {
	ns.RLock()
	defer ns.RUnlock()
	return ns.Context.Get(n)
}

// Has checks for the existence of a variable and returns its context
func (ns *namespace) Has(n Name) (Context, bool) {
	ns.RLock()
	defer ns.RUnlock()
	if _, ok := ns.Context.Has(n); ok {
		return ns, true
	}
	return ns, false
}

// Delete removes a variable by name
func (ns *namespace) Delete(n Name) {
	ns.Lock()
	defer ns.Unlock()
	ns.Context.Delete(n)
}

// Put sets a variable by name if it hasn't already been set
func (ns *namespace) Put(n Name, v Value) {
	ns.Lock()
	defer ns.Unlock()
	if _, ok := ns.Context.Has(n); ok {
		panic(ErrStr(AlreadyBound, n))
	}
	ns.Context.Put(n, v)
}

func (ns *namespace) Intern(n Name) Symbol {
	return ns.symbols.Get(n, func() u.Any {
		return localSymbol(n)
	}).(Symbol)
}

func (ns *namespace) Domain() Name {
	return LocalDomain
}

func (ns *namespace) Str() Str {
	return "(ns *local*)"
}

func (*namespace) NamespaceType() {}

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

func (ns *qualifiedNamespace) Domain() Name {
	return ns.domain
}

func (ns *qualifiedNamespace) Str() Str {
	return "(ns " + Str(ns.domain) + ")"
}

// Get retrieves a variable by name
func (ns *childNamespace) Get(n Name) (Value, bool) {
	if v, ok := ns.Namespace.Get(n); ok {
		return v, true
	}
	return ns.parent.Get(n)
}

// Has checks for the existence of a variable and returns its context
func (ns *childNamespace) Has(n Name) (Context, bool) {
	if c, ok := ns.Namespace.Has(n); ok {
		return c, true
	}
	return ns.parent.Has(n)
}

// GetNamespace returns the Namespace for the specified domain.
func GetNamespace(n Name) Namespace {
	return namespaces.Get(n, func() u.Any {
		ns := &qualifiedNamespace{
			namespace: namespace{
				Context: Variables{},
				symbols: u.NewCache(),
			},
			domain: n,
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
	return ChildLocals(&withNamespace{
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
	namespaces.Get(LocalDomain, func() u.Any {
		ns := &namespace{
			Context: Variables{},
			symbols: u.NewCache(),
		}
		ns.Put(ContextDomain, ns)
		return ns
	})

	bi := namespaces.Get(BuiltInDomain, func() u.Any {
		ns := &qualifiedNamespace{
			namespace: namespace{
				Context: Variables{},
				symbols: u.NewCache(),
			},
			domain: BuiltInDomain,
		}
		ns.Put(ContextDomain, ns)
		return ns
	}).(Namespace)

	namespaces.Get(UserDomain, func() u.Any {
		ns := &childNamespace{
			Namespace: &qualifiedNamespace{
				namespace: namespace{
					Context: Variables{},
					symbols: u.NewCache(),
				},
				domain: UserDomain,
			},
			parent: bi,
		}
		ns.Put(ContextDomain, ns)
		return ns
	})
}
