package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
)

type (
	inMapper  func(a.Value) reflect.Value
	outMapper func(reflect.Value) a.Value
)

var (
	convertIn   map[reflect.Kind]inMapper
	convertOut  map[reflect.Kind]outMapper
	valueHolder a.Value
	valueType   = reflect.TypeOf(&valueHolder).Elem()
)

func valueToReflect(v a.Value) reflect.Value {
	return reflect.ValueOf(v)
}

func reflectToValue(v reflect.Value) a.Value {
	return v.Interface().(a.Value)
}

func boolToReflect(v a.Value) reflect.Value {
	return reflect.ValueOf(bool(a.AssertBool(v)))
}

func reflectToBool(v reflect.Value) a.Value {
	return a.Bool(v.Bool())
}

func strToReflect(v a.Value) reflect.Value {
	return reflect.ValueOf(string(a.AssertStr(v)))
}

func reflectToStr(v reflect.Value) a.Value {
	return a.Str(v.String())
}

func numberToFloat32(v a.Value) reflect.Value {
	f, _ := a.AssertNumber(v).Float64()
	return reflect.ValueOf(float32(f))
}

func numberToFloat64(v a.Value) reflect.Value {
	f, _ := a.AssertNumber(v).Float64()
	return reflect.ValueOf(f)
}

func floatToNumber(v reflect.Value) a.Value {
	return a.NewFloat(v.Float())
}

func numberToInt(v a.Value) reflect.Value {
	i := a.AssertInteger(v)
	return reflect.ValueOf(i)
}

func intToNumber(v reflect.Value) a.Value {
	return a.NewFloat(float64(v.Int()))
}

func getConvertIn(t reflect.Type) inMapper {
	if t.Implements(valueType) {
		return valueToReflect
	}
	if c, ok := convertIn[t.Kind()]; ok {
		return c
	}
	return func(_ a.Value) reflect.Value {
		panic(a.ErrStr(BadConversionType, t))
	}
}

func getConvertOut(t reflect.Type) outMapper {
	if t.Implements(valueType) {
		return reflectToValue
	}
	if c, ok := convertOut[t.Kind()]; ok {
		return c
	}
	return func(_ reflect.Value) a.Value {
		panic(a.ErrStr(BadConversionType, t))
	}
}

func init() {
	convertIn = map[reflect.Kind]inMapper{
		reflect.Bool:    boolToReflect,
		reflect.Int:     numberToInt,
		reflect.Int8:    numberToInt,
		reflect.Int16:   numberToInt,
		reflect.Int32:   numberToInt,
		reflect.Int64:   numberToInt,
		reflect.Uint:    numberToInt,
		reflect.Uint8:   numberToInt,
		reflect.Uint16:  numberToInt,
		reflect.Uint32:  numberToInt,
		reflect.Uint64:  numberToInt,
		reflect.Float32: numberToFloat32,
		reflect.Float64: numberToFloat64,
		reflect.String:  strToReflect,
	}

	convertOut = map[reflect.Kind]outMapper{
		reflect.Bool:    reflectToBool,
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
		reflect.String:  reflectToStr,
		reflect.Func:    funcToApplicable,
		reflect.Struct:  Wrap,
		reflect.Ptr:     Wrap,
	}
}
