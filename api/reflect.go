package api

import "reflect"

// Native is the interface for wrapped Go values
type Native interface {
	Value
	Annotated
	WrappedValue() interface{}
}

type native struct {
	value reflect.Value
	meta  Metadata
}

var typeMeta = map[reflect.Type]Metadata{}

// NewNative wraps a native value using Go's reflection API
func NewNative(a interface{}) Native {
	v := reflect.ValueOf(a)
	t := v.Type()
	md, ok := typeMeta[t]
	if !ok {
		md = makeMetadata(t)
		typeMeta[t] = md
	}

	return &native{
		value: v,
		meta:  md,
	}
}

func makeMetadata(t reflect.Type) Metadata {
	return Metadata{
		MetaType: Name(typeName(t)),
	}
}

func typeName(t reflect.Type) string {
	return t.String()
}

// WrappedValue returns the wrapped Go value
func (n *native) WrappedValue() interface{} {
	return n.value.Interface()
}

// Metadata makes Native Annotated
func (n *native) Metadata() Metadata {
	return n.meta
}

// WithMetadata copies the Native with new Metadata
func (n *native) WithMetadata(md Metadata) Annotated {
	return &native{
		value: n.value,
		meta:  n.meta.Merge(md),
	}
}

// Type returns the type name of the wrapped native
func (n *native) Type() Name {
	return n.meta[MetaType].(Name)
}

// Str converts this Value into a Str
func (n *native) Str() Str {
	return MakeDumpStr(n)
}
