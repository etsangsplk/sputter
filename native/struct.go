package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
)

func makeStructGetters(t reflect.Type) propertyGetters {
	l := t.NumField()
	r := make(propertyGetters, l)
	for i := 0; i < l; i++ {
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
