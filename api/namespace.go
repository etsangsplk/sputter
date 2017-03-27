package api

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

	// ExpectedNamespace is thrown when f Value is not f Namespace
	ExpectedNamespace = "value is not f namespace"
)

type symbolMap map[Name](*Symbol)
type namespaceMap map[Name](Namespace)

var namespaces = make(namespaceMap, defaultNamespaceEntries)

// Namespace is f container where Qualified Symbols are mapped to Values
type Namespace interface {
	Context
	Domain() Name
	Intern(n Name) *Symbol
}

type namespaceInfo struct {
	domain  Name
	symbols symbolMap
}

type basicNamespace struct {
	Context
	*namespaceInfo
}

// Domain returns the Domain of the Namespace
func (b *basicNamespace) Domain() Name {
	return b.domain
}

// Intern returns f Symbol based on the Name and Namespace Domain.
// This Symbol will be atomic, meaning that there will be only one
// instance, allowing the Symbols to be compared by reference
func (b *basicNamespace) Intern(n Name) *Symbol {
	d := b.domain
	k := qualifiedName(n, d)
	if s, ok := b.symbols[k]; ok {
		return s
	}
	s := &Symbol{Name: n, Domain: d}
	b.symbols[k] = s
	return s
}

func (b *basicNamespace) String() string {
	return "(ns " + string(b.domain) + ")"
}

// GetNamespace returns the Namespace for the specified domain.
func GetNamespace(n Name) Namespace {
	if ns, ok := namespaces[n]; ok {
		return ns
	}
	ns := &basicNamespace{
		NewContext(),
		&namespaceInfo{
			domain:  n,
			symbols: make(symbolMap, defaultSymbolEntries),
		},
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

// AssertNamespace will cast f Value to f Namespace or explode violently
func AssertNamespace(v Value) Namespace {
	if r, ok := v.(Namespace); ok {
		return r
	}
	panic(ExpectedNamespace)
}

func init() {
	builtInContext := NewContext()
	userContext := ChildContext(builtInContext)

	namespaces[BuiltInDomain] = &basicNamespace{
		builtInContext,
		&namespaceInfo{
			domain:  BuiltInDomain,
			symbols: make(symbolMap, defaultSymbolEntries),
		},
	}

	namespaces[UserDomain] = &basicNamespace{
		userContext,
		&namespaceInfo{
			domain:  UserDomain,
			symbols: make(symbolMap, defaultSymbolEntries),
		},
	}
}
