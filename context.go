package sputter

const defaultVarSize = 16

// Variables are how a closure stores key/value pairs
type Variables map[string]Value

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
func (c *Context) Get(name string) (Value, bool) {
	if value, ok := c.vars[name]; ok {
		return value, true
	} else if c.parent != nil {
		return c.parent.Get(name)
	}
	return nil, false
}

// Put puts a value into the immediate Context
func (c *Context) Put(name string, value Value) {
	c.vars[name] = value
}

func (c *Context) PutFunction(f *Function) {
	c.vars[f.name] = f
}

// Evaluable can be evaluated against a Context
type Evaluable interface {
	Evaluate(c *Context) Value
}
