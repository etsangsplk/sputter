package native

import (
	"reflect"
	"regexp"
	"strings"

	a "github.com/kode4food/sputter/api"
)

// BadConversionType is raised if a type can't be converted
const BadConversionType = "can't convert between type: %s"

type outMapper func(reflect.Value) a.Value

var (
	camelCase  = regexp.MustCompile("[a-z][A-Z]")
	convertOut map[reflect.Kind]outMapper
)

// New wraps a native wrapped using Go's reflection API
func New(i interface{}) a.Value {
	v := reflect.ValueOf(i)
	t := v.Type()
	if c, ok := convertOut[t.Kind()]; ok {
		return c(v)
	}
	panic(a.Err(BadConversionType, t))
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

func valueToBool(v reflect.Value) a.Value {
	return a.Bool(v.Bool())
}

func valueToStr(v reflect.Value) a.Value {
	return a.Str(v.String())
}

func floatValueToNumber(v reflect.Value) a.Value {
	return a.NewFloat(v.Float())
}

func intValueToNumber(v reflect.Value) a.Value {
	return a.NewFloat(float64(v.Int()))
}

func init() {
	convertOut = map[reflect.Kind]outMapper{
		reflect.Bool:    valueToBool,
		reflect.Int:     intValueToNumber,
		reflect.Int8:    intValueToNumber,
		reflect.Int16:   intValueToNumber,
		reflect.Int32:   intValueToNumber,
		reflect.Int64:   intValueToNumber,
		reflect.Uint:    intValueToNumber,
		reflect.Uint8:   intValueToNumber,
		reflect.Uint16:  intValueToNumber,
		reflect.Uint32:  intValueToNumber,
		reflect.Uint64:  intValueToNumber,
		reflect.Float32: floatValueToNumber,
		reflect.Float64: floatValueToNumber,
		reflect.String:  valueToStr,
		reflect.Struct:  Wrap,
		reflect.Ptr:     Wrap,
	}
}
