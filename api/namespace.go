package api

type namespaceMap map[Name](Context)

var namespaces = make(namespaceMap, 128)

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

// GetNamespace returns the Context for the specified domain.
func GetNamespace(domain Name) Context {
	if ns, ok := namespaces[domain]; ok {
		return ns
	}
	ns := NewContext()
	namespaces[domain] = ns
	return ns
}

// GetContextNamespace resolves the Namespace based on its Context
func GetContextNamespace(c Context) Context {
	if v, ok := c.Get(ContextDomain); ok {
		ns := AssertName(v)
		return GetNamespace(ns)
	}
	return GetNamespace(UserDomain)
}

func init() {
	builtInNamespace := NewContext()
	userNamespace := ChildContext(builtInNamespace)

	namespaces[BuiltInDomain] = builtInNamespace
	namespaces[UserDomain] = userNamespace
}
