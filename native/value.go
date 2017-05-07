package native

import (
	"reflect"
	"regexp"
	"strings"

	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

// BadConversionType is raised if a value type can't be converted
const BadConversionType = "can't convert between type: %s"

// Value is the interface for wrapped Go values
type Value interface {
	a.Value
	a.Annotated
	a.Getter
	Wrapped() interface{}
}

type value struct {
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

type (
	inMapper        func(a.Value) reflect.Value
	propertyGetters map[string]outMapper
)

var (
	types      = u.NewCache()
	camelCase  = regexp.MustCompile("[a-z][A-Z]")
	convertOut map[reflect.Kind]outMapper
)

// NewValue wraps a value value using Go's reflection API
func NewValue(a interface{}) Value {
	return wrapValue(reflect.ValueOf(a)).(Value)
}

func wrapValue(v reflect.Value) a.Value {
	t := v.Type()
	ti := getTypeInfo(t)

	return &value{
		value:    v,
		typeInfo: ti,
		meta:     ti.meta,
	}
}

// Wrapped returns the wrapped Go value
func (n *value) Wrapped() interface{} {
	return n.value.Interface()
}

// Metadata makes Value Annotated
func (n *value) Metadata() a.Metadata {
	return n.meta
}

// WithMetadata copies the Value with new Metadata
func (n *value) WithMetadata(md a.Metadata) a.Annotated {
	return &value{
		value:    n.value,
		typeInfo: n.typeInfo,
		meta:     n.meta.Merge(md),
	}
}

func (n *value) Get(key a.Value) (a.Value, bool) {
	name := string(key.Str())
	if g, ok := n.typeInfo.getters[name]; ok {
		return g(n.value), true
	}
	return a.Nil, false
}

// Type returns the type name of the wrapped value
func (n *value) Type() a.Name {
	return n.meta[a.MetaType].(a.Name)
}

// Str converts this Value into a Str
func (n *value) Str() a.Str {
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
	g := propertyGetters{}
	switch t.Kind() {
	case reflect.Ptr:
		for k, v := range makePropertyGetters(t.Elem()) {
			g[k] = makeIndirectGetter(v)
		}
	case reflect.Struct:
		for k, v := range makeStructGetters(t) {
			g[k] = v
		}
	}
	for k, v := range makeMethodGetters(t) {
		g[k] = v
	}

	return g
}

func makeIndirectGetter(g outMapper) outMapper {
	return func(v reflect.Value) a.Value {
		return g(v.Elem())
	}
}

func typeName(t reflect.Type) a.Name {
	return a.Name(kebabCase(t.String()))
}

func kebabCase(n string) string {
	r := camelCase.ReplaceAllStringFunc(n, func(s string) string {
		return s[:1] + "-" + s[1:]
	})
	return strings.ToLower(r)
}
