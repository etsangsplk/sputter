package api

import (
	"reflect"
	"regexp"
	"strings"

	u "github.com/kode4food/sputter/util"
)

// BadConversionType is raised if a native type can't be converted
const BadConversionType = "can't convert native type: %s"

// Native is the interface for wrapped Go values
type Native interface {
	Value
	Annotated
	Getter
	NativeValue() interface{}
}

type native struct {
	value    reflect.Value
	typeInfo *typeInfo
	meta     Metadata
}

type typeInfo struct {
	name    Name
	typ     reflect.Type
	getters propertyGetters
	meta    Metadata
}

type (
	nativeMapper    func(v reflect.Value) Value
	propertyGetters map[string]nativeMapper
)

var (
	types      = u.NewCache()
	camelCase  = regexp.MustCompile("[a-z][A-Z]")
	converters map[reflect.Kind]nativeMapper
)

// NewNative wraps a native value using Go's reflection API
func NewNative(a interface{}) Native {
	return toNative(reflect.ValueOf(a)).(Native)
}

func getTypeInfo(t reflect.Type) *typeInfo {
	tn := typeName(t)
	return types.Get(tn, func() u.Any {
		return &typeInfo{
			name: tn,
			meta: Metadata{
				MetaType: tn,
			},
			getters: makePropertyGetters(t),
		}
	}).(*typeInfo)
}

func kebabCase(n string) string {
	r := camelCase.ReplaceAllStringFunc(n, func(s string) string {
		return s[:1] + "-" + s[1:]
	})
	return strings.ToLower(r)
}

func typeName(t reflect.Type) Name {
	return Name(kebabCase(t.String()))
}

// NativeValue returns the wrapped Go value
func (n *native) NativeValue() interface{} {
	return n.value.Interface()
}

// Metadata makes Native Annotated
func (n *native) Metadata() Metadata {
	return n.meta
}

// WithMetadata copies the Native with new Metadata
func (n *native) WithMetadata(md Metadata) Annotated {
	return &native{
		value:    n.value,
		typeInfo: n.typeInfo,
		meta:     n.meta.Merge(md),
	}
}

func (n *native) Get(key Value) (Value, bool) {
	name := string(key.Str())
	if g, ok := n.typeInfo.getters[name]; ok {
		return g(n.value), true
	}
	return Nil, false
}

// Type returns the type name of the wrapped native
func (n *native) Type() Name {
	return n.meta[MetaType].(Name)
}

// Str converts this Value into a Str
func (n *native) Str() Str {
	return MakeDumpStr(n)
}

func makePropertyGetters(t reflect.Type) propertyGetters {
	g := propertyGetters{}
	switch t.Kind() {
	case reflect.Ptr:
		for n, v := range makePropertyGetters(t.Elem()) {
			g[n] = makeIndirectGetter(v)
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)
			if fi.PkgPath != "" {
				// we only surface exported fields
				continue
			}
			n := kebabCase(fi.Name)
			g[n] = makeFieldGetter(i, fi)
		}
	}
	for i := 0; i < t.NumMethod(); i++ {
		mi := t.Method(i)
		n := kebabCase(mi.Name)
		g[n] = makeMethodGetter(mi)
	}
	return g
}

func makeIndirectGetter(g nativeMapper) nativeMapper {
	return func(v reflect.Value) Value {
		return g(v.Elem())
	}
}

func makeFieldGetter(idx int, fi reflect.StructField) nativeMapper {
	if c, ok := converters[fi.Type.Kind()]; ok {
		return func(v reflect.Value) Value {
			return c(v.Field(idx))
		}
	}

	return func(v reflect.Value) Value {
		return badConvert(v.Field(idx))
	}
}

func nativeToBool(v reflect.Value) Value {
	return Bool(v.Bool())
}

func nativeToStr(v reflect.Value) Value {
	return Str(v.String())
}

func toNative(v reflect.Value) Value {
	t := v.Type()
	ti := getTypeInfo(t)

	return &native{
		value:    v,
		typeInfo: ti,
		meta:     ti.meta,
	}
}

func floatToNumber(v reflect.Value) Value {
	return NewFloat(v.Float())
}

func intToNumber(v reflect.Value) Value {
	return NewFloat(float64(v.Int()))
}

func badConvert(v reflect.Value) Value {
	panic(Err(BadConversionType, v.Type().String()))
}

func makeMethodGetter(mi reflect.Method) nativeMapper {
	return func(v reflect.Value) Value {
		return Nil
	}
}

func init() {
	converters = map[reflect.Kind]nativeMapper{
		reflect.Bool:    nativeToBool,
		reflect.Int:     intToNumber,
		reflect.Int8:    intToNumber,
		reflect.Int16:   intToNumber,
		reflect.Int32:   intToNumber,
		reflect.Int64:   intToNumber,
		reflect.Uint:    intToNumber,
		reflect.Uint8:   intToNumber,
		reflect.Uint16:  intToNumber,
		reflect.Uint32:  intToNumber,
		reflect.Uint64:  intToNumber,
		reflect.Float32: floatToNumber,
		reflect.Float64: floatToNumber,
		reflect.String:  nativeToStr,
		reflect.Struct:  toNative,
		reflect.Ptr:     toNative,
	}
}
