package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

// Wrapped is the interface for wrapped Go values
type Wrapped interface {
	a.Value
	a.Annotated
	a.Getter
	Wrapped() interface{}
}

type wrapped struct {
	value    reflect.Value
	typeInfo *typeInfo
	meta     a.Metadata
}

type typeInfo struct {
	name    a.Name
	typ     reflect.Type
	getters propertyGetters
	meta    a.Metadata
}

type propertyGetters map[string]outMapper

var types = u.NewCache()

// Wrap wraps a non-primitive Go value
func Wrap(v reflect.Value) a.Value {
	ti := getTypeInfo(v.Type())

	return &wrapped{
		value:    v,
		typeInfo: ti,
		meta:     ti.meta,
	}
}

// Wrapped returns the wrapped Go value
func (n *wrapped) Wrapped() interface{} {
	return n.value.Interface()
}

// Metadata makes wrapped Annotated
func (n *wrapped) Metadata() a.Metadata {
	return n.meta
}

// WithMetadata copies the wrapped with new Metadata
func (n *wrapped) WithMetadata(md a.Metadata) a.Annotated {
	return &wrapped{
		value:    n.value,
		typeInfo: n.typeInfo,
		meta:     n.meta.Merge(md),
	}
}

func (n *wrapped) Get(key a.Value) (a.Value, bool) {
	name := string(key.Str())
	if g, ok := n.typeInfo.getters[name]; ok {
		return g(n.value), true
	}
	return a.Nil, false
}

func (n *wrapped) Type() a.Name {
	return n.meta[a.MetaType].(a.Name)
}

func (n *wrapped) Eval(_ a.Context) a.Value {
	return n
}

func (n *wrapped) Str() a.Str {
	return a.MakeDumpStr(n)
}

func getTypeInfo(t reflect.Type) *typeInfo {
	tn := typeName(t)
	return types.Get(tn, func() u.Any {
		return &typeInfo{
			name: tn,
			meta: a.Metadata{
				a.MetaType: tn,
			},
			getters: makePropertyGetters(t),
		}
	}).(*typeInfo)
}

func makePropertyGetters(t reflect.Type) propertyGetters {
	mg := makeMethodGetters(t)
	var r propertyGetters

	switch t.Kind() {
	case reflect.Ptr:
		r = makePointerGetters(t)
	case reflect.Struct:
		r = makeStructGetters(t)
	default:
		return mg
	}

	// methods always win
	for k, v := range mg {
		r[k] = v
	}
	return r
}

func makePointerGetters(t reflect.Type) propertyGetters {
	pg := makePropertyGetters(t.Elem())
	r := make(propertyGetters, len(pg))
	for k, v := range pg {
		r[k] = makeIndirectGetter(v)
	}
	return r
}

func makeIndirectGetter(g outMapper) outMapper {
	return func(v reflect.Value) a.Value {
		return g(v.Elem())
	}
}
