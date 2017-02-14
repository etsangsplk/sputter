package api

const defaultVarSize = 16

// Variables are how a closure stores name/value pairs
type Variables map[Name]Value

// Context represents a functional closure
type Context struct {
	parent *Context
	vars   Variables
}

// NewContext creates a new Context instance
func NewContext() *Context {
	return &Context{nil, make(Variables, defaultVarSize)}
}

// Child creates a child Context instance
func (c *Context) Child() *Context {
	return &Context{c, make(Variables, defaultVarSize)}
}

// Get retrieves a value from the Context chain
func (c *Context) Get(n Name) (v Value, bound bool) {
	if v, ok := c.vars[n]; ok {
		return v, true
	} else if c.parent != nil {
		return c.parent.Get(n)
	}
	return EmptyList, false
}

// Globals retrieves the Root Context (one with no parent)
func (c *Context) Globals() *Context {
	t := c
	for t.parent != nil {
		t = t.parent
	}
	return t
}

// Put puts a Value into the immediate Context
func (c *Context) Put(n Name, v Value) *Context {
	c.vars[n] = v
	return c
}

// PutFunction puts a Function into the immediate Context by its name
func (c *Context) PutFunction(f *Function) *Context {
	c.vars[f.Name] = f
	return c
}

// Evaluable can be evaluated against a Context
type Evaluable interface {
	Evaluate(c *Context) Value
}
