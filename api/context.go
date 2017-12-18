package api

const defaultContextEntries = 16

type (
	// Context represents a mutable variable scope
	Context interface {
		Value
		Get(Name) (Value, bool)
		Has(Name) (Context, bool)
		Put(Name, Value)
		Delete(Name)
	}

	childContext struct {
		Context
		parent Context
	}
)

// NewContext creates a new independent Context instance
func NewContext() Context {
	return &WriteOnceVariables{
		Variables: make(Variables, defaultContextEntries),
	}
}

// ChildContext creates a new child Context of the provided parent
func ChildContext(parent Context) Context {
	return &childContext{
		Context: NewContext(),
		parent:  parent,
	}
}

// ChildLocals creates a new child Context for local variables
func ChildLocals(parent Context) Context {
	return &childContext{
		Context: Variables{},
		parent:  parent,
	}
}

// ChildVariables creates a new child Context with Variables
func ChildVariables(parent Context, vars Variables) Context {
	return &childContext{
		Context: vars,
		parent:  parent,
	}
}

// NewClosure copies variables into an isolated Context
func NewClosure(parent Context, names Names) Context {
	ns := GetContextNamespace(parent)
	vars := make(Variables, len(names))
	for _, n := range names {
		if v, ok := parent.Get(n); ok {
			vars[n] = v
		}
	}
	return &childContext{
		Context: vars,
		parent:  ns,
	}
}

// Get retrieves a value from the Context chain
func (c *childContext) Get(n Name) (Value, bool) {
	if v, ok := c.Context.Get(n); ok {
		return v, true
	}
	return c.parent.Get(n)
}

// Has looks up the Context in which a value exists
func (c *childContext) Has(n Name) (Context, bool) {
	if _, ok := c.Context.Has(n); ok {
		return c, true
	}
	return c.parent.Has(n)
}
