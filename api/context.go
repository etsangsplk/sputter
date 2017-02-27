package api

const defaultVarsSize = 16

// Context represents a variable scope
type Context interface {
	Globals() Context
	Get(n Name) (v Value, bound bool)
	Put(n Name, v Value) Context
}

// BasicContext is the most basic Context implementation
type basicContext struct {
	parent Context
	vars   Variables
	global bool
}

// NewGlobalContext creates a new global Context instance that
// chains up to a reserved Context for special forms
func NewGlobalContext(reserved Context) Context {
	return &basicContext{
		parent: reserved,
		vars:   make(Variables, defaultVarsSize),
		global: true,
	}
}

// NewContext creates a new global Context instance
func NewContext() Context {
	return NewGlobalContext(nil)
}

// ChildContext creates a new child Context of the provided parent
func ChildContext(parent Context) Context {
	return &basicContext{
		parent: parent,
		vars:   make(Variables, defaultVarsSize),
	}
}

// Globals retrieves the Root Context (one with no parent)
func (c *basicContext) Globals() Context {
	if c.global || c.parent == nil {
		return c
	}
	return c.parent.Globals()
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
func (c *basicContext) Put(n Name, v Value) Context {
	c.vars[n] = v
	return c
}

// PutFunction puts a Function into a Context by its name
func PutFunction(c Context, f *Function) Context {
	return c.Put(f.Name, f)
}
