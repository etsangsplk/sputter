package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

type (
	// GoValue is the interface for wrapped Go values
	GoValue struct {
		value    reflect.Value
		typeInfo *typeInfo
		meta     a.Object
	}

	typeInfo struct {
		name    a.Name
		typ     reflect.Type
		getters propertyGetters
		meta    a.Object
	}

	propertyGetters map[string]outMapper
)

var types = u.NewCache()

// Wrap wraps a non-primitive Go value
func Wrap(v reflect.Value) a.Value {
	ti := getTypeInfo(v.Type())

	return &GoValue{
		value:    v,
		typeInfo: ti,
		meta:     ti.meta,
	}
}

// Wrapped returns the wrapped Go value
func (v *GoValue) Wrapped() interface{} {
	return v.value.Interface()
}

// Metadata makes wrapped Annotated
func (v *GoValue) Metadata() a.Object {
	return v.meta
}

// WithMetadata copies the wrapped with new Metadata
func (v *GoValue) WithMetadata(md a.Object) a.AnnotatedValue {
	return &GoValue{
		value:    v.value,
		typeInfo: v.typeInfo,
		meta:     v.meta.Child(md.Flatten()),
	}
}

// Get returns a property by key from the wrapped Go value
func (v *GoValue) Get(key a.Value) (a.Value, bool) {
	name := string(key.Str())
	if g, ok := v.typeInfo.getters[name]; ok {
		return g(v.value), true
	}
	return a.Nil, false
}

// Type returns the name of the wrapped Go value's type
func (v *GoValue) Type() a.Name {
	t, _ := v.meta.Get(a.TypeKey)
	return t.(a.Name)
}

// Str converts this Go value to a Str
func (v *GoValue) Str() a.Str {
	return a.MakeDumpStr(v)
}

func getTypeInfo(t reflect.Type) *typeInfo {
	tn := typeName(t)
	return types.Get(tn, func() u.Any {
		return &typeInfo{
			name:    tn,
			meta:    a.Properties{a.TypeKey: tn},
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
