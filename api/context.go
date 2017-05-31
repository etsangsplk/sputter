package api

import "sync"

const defaultContextEntries = 16

// AlreadyBound is thrown when an attempt is made to rebind a Name
const AlreadyBound = "symbol is already bound in this context: %s"

// Context represents a mutable variable scope
type Context interface {
	Get(Name) (Value, bool)
	Put(Name, Value)
	Delete(Name)
}

type context struct {
	sync.RWMutex
	parent Context
	vars   Variables
}

type rootContext struct {
	*context
}

// NewContext creates a new independent Context instance
func NewContext() Context {
	return &rootContext{
		context: &context{
			vars: make(Variables, defaultContextEntries),
		},
	}
}

// ChildContext creates a new child Context of the provided parent
func ChildContext(parent Context) Context {
	return &context{
		parent: parent,
		vars:   make(Variables, defaultContextEntries),
	}
}

// ChildContextVars creates a new child Context with Variables
func ChildContextVars(parent Context, vars Variables) Context {
	return &context{
		parent: parent,
		vars:   vars,
	}
}

// Get retrieves a value from the Context
func (c *rootContext) Get(n Name) (Value, bool) {
	c.RLock()
	defer c.RUnlock()
	if v, ok := c.vars[n]; ok {
		return v, ok
	}
	return Nil, false
}

// Get retrieves a value from the Context chain
func (c *context) Get(n Name) (Value, bool) {
	c.RLock()
	defer c.RUnlock()
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
	c.Lock()
	defer c.Unlock()
	if _, ok := c.vars[n]; ok {
		panic(Err(AlreadyBound, n))
	}
	c.vars[n] = v
}

// Delete removes a Value from the immediate Context
func (c *context) Delete(n Name) {
	c.Lock()
	defer c.Unlock()
	delete(c.vars, n)
}
