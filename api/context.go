package api

const defaultContextEntries = 16

// Context represents a variable scope
type Context interface {
	Get(n Name) (v Value, bound bool)
	Put(n Name, v Value)
	Delete(n Name)
}

// basicContext is the most basic Context implementation
type basicContext struct {
	parent Context
	vars   Variables
}

// NewContext creates a new independent Context instance
func NewContext() Context {
	return &basicContext{
		parent: nil,
		vars:   make(Variables, defaultContextEntries),
	}
}

// ChildContext creates a new child Context of the provided parent
func ChildContext(parent Context) Context {
	return &basicContext{
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
func (c *basicContext) Get(n Name) (Value, bool) {
	if v, ok := c.vars[n]; ok {
		return v, true
	}
	if c.parent != nil {
		return c.parent.Get(n)
	}
	return Nil, false
}

// Put puts a Value into the immediate Context
func (c *basicContext) Put(n Name, v Value) {
	c.vars[n] = v
}

// Delete removes a Value from the immediate Context
func (c *basicContext) Delete(n Name) {
	delete(c.vars, n)
}
