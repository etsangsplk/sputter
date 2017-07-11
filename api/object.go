package api

import "bytes"

// ValueNotFound is raised when a Value can't be retrieved by Name
const ValueNotFound = "value not found in object: %s"

type (
	// Object represents the basic element of the universal design pattern
	Object interface {
		Value
		Mapped
		GetValue(Value) Value
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
	MetaType: Name("object"),
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

// GetValue retrieves a Value from an Object or explodes
func (p Properties) GetValue(k Value) Value {
	if v, ok := p.Get(k); ok {
		return v
	}
	panic(Err(ValueNotFound, k))
}

// Child instantiates a new UDP child Object with Properties
func (p Properties) Child(props Properties) Object {
	return &object{
		props:  props,
		parent: p,
	}
}

// Flatten converts the Properties of an Object graph into a single set
func (p Properties) Flatten() Properties {
	r := make(Properties, len(p))
	for k, v := range p {
		r[k] = v
	}
	return r
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

func (o *object) GetValue(k Value) Value {
	if v, ok := o.Get(k); ok {
		return v
	}
	panic(Err(ValueNotFound, k))
}

// Child instantiates a new UDP child Object with Properties
func (o *object) Child(props Properties) Object {
	return &object{
		props:  props,
		parent: o,
	}
}

func (o *object) Flatten() Properties {
	r := o.parent.Flatten()
	for k, v := range o.props {
		r[k] = v
	}
	return r
}

func (o *object) Str() Str {
	return o.Flatten().Str()
}
