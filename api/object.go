package api

import "bytes"

// ValueNotFound is raised when a Value can't be retrieved by Name
const ValueNotFound = "value not found in object: %s"

type (
	// Object represents the basic element of the universal design pattern
	Object interface {
		Value
		Mapped
		MustGet(Value) Value
		Child(Properties) Object
		Flatten() Properties
	}

	// Properties maps Values to Values for UDP Objects
	Properties map[Value]Value

	object struct {
		props  Properties
		parent Object
	}
)

var objectPrototype = Properties{
	TypeKey: Name("object"),
}

// NewObject instantiates a new UDP Object with Variables
func NewObject(props Properties) Object {
	return objectPrototype.Child(props)
}

// Get attempts to retrieve a Value from an Object
func (p Properties) Get(k Value) (Value, bool) {
	if v, ok := p[k]; ok {
		return v, ok
	}
	return Nil, false
}

// MustGet retrieves a Value from an Object or explodes
func (p Properties) MustGet(k Value) Value {
	if v, ok := p.Get(k); ok {
		return v
	}
	panic(ErrStr(ValueNotFound, k))
}

// Child instantiates a new UDP child Object with Properties
func (p Properties) Child(props Properties) Object {
	return &object{
		props:  props,
		parent: p,
	}
}

// Flatten against Properties just returns the Properties (which are flat)
func (p Properties) Flatten() Properties {
	return p
}

// Str converts this Value into a Str
func (p Properties) Str() Str {
	var buf bytes.Buffer
	buf.WriteString("{")
	idx := 0
	for k, v := range p {
		if idx > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(string(k.Str()))
		buf.WriteString(" ")
		buf.WriteString(string(v.Str()))
		idx++
	}
	buf.WriteString("}")
	return Str(buf.String())
}

func (o *object) Get(k Value) (Value, bool) {
	if v, ok := o.props[k]; ok {
		return v, true
	}
	return o.parent.Get(k)
}

func (o *object) MustGet(k Value) Value {
	if v, ok := o.Get(k); ok {
		return v
	}
	panic(ErrStr(ValueNotFound, k))
}

// Child instantiates a new UDP child Object with Properties
func (o *object) Child(props Properties) Object {
	return &object{
		props:  props,
		parent: o,
	}
}

func (o *object) Flatten() Properties {
	pf := o.parent.Flatten()
	sf := o.props
	r := make(Properties, len(pf)+len(sf))
	for k, v := range pf {
		r[k] = v
	}
	for k, v := range sf {
		r[k] = v
	}
	return r
}

func (o *object) Str() Str {
	return o.Flatten().Str()
}
