package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
)

func makeStructGetters(t reflect.Type) propertyGetters {
	r := propertyGetters{}
	for i := 0; i < t.NumField(); i++ {
		fi := t.Field(i)
		if fi.PkgPath != "" {
			continue // only surface exported fields
		}
		n := kebabCase(fi.Name)
		r[n] = makeFieldGetter(i, fi)
	}
	return r
}

func makeFieldGetter(idx int, fi reflect.StructField) outMapper {
	if c, ok := convertOut[fi.Type.Kind()]; ok {
		return func(v reflect.Value) a.Value {
			return c(v.Field(idx))
		}
	}

	return func(v reflect.Value) a.Value {
		return badConvert(v.Field(idx))
	}
}

func badConvert(v reflect.Value) a.Value {
	panic(a.Err(BadConversionType, v.Type().String()))
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
		reflect.Struct:  wrapValue,
		reflect.Ptr:     wrapValue,
	}
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
