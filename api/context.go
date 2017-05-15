package api

const defaultContextEntries = 16

// AlreadyBound is thrown when an attempt is made to rebind a Name
const AlreadyBound = "symbol is already bound in this context: %s"

// Context represents a variable scope
type Context interface {
	Get(Name) (Value, bool)
	Put(Name, Value)
	Delete(Name)
}

type context struct {
	parent Context
	vars   Variables
}

// NewContext creates a new independent Context instance
func NewContext() Context {
	return &context{
		parent: nil,
		vars:   make(Variables, defaultContextEntries),
	}
}

// ChildContext creates a new child Context of the provided parent
func ChildContext(parent Context) Context {
	return &context{
		parent: parent,
		vars:   make(Variables, defaultContextEntries),
	}
}

// NewEvalContext creates a new Context instance that
// chains up to the UserDomain Context for special forms
func NewEvalContext() Context {
	ns := GetNamespace(UserDomain)
	c := ChildContext(ns)
	c.Put(ContextDomain, ns)
	return c
}

// Get retrieves a value from the Context chain
func (c *context) Get(n Name) (Value, bool) {
	if v, ok := c.vars[n]; ok {
		return v, true
	}
	if c.parent != nil {
		return c.parent.Get(n)
	}
	return Nil, false
}

// Put puts a Value into the immediate Context
func (c *context) Put(n Name, v Value) {
	if _, ok := c.vars[n]; ok {
		panic(Err(AlreadyBound, n))
	}
	c.vars[n] = v
}

// Delete removes a Value from the immediate Context
func (c *context) Delete(n Name) {
	delete(c.vars, n)
}
