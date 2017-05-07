package native

import (
	a "github.com/kode4food/sputter/api"
	"reflect"
)

type (
	inMapper  func(a.Value) reflect.Value
	outMapper func(reflect.Value) a.Value
)

var (
	convertIn  map[reflect.Kind]inMapper
	convertOut map[reflect.Kind]outMapper
)

func boolToNative(v a.Value) reflect.Value {
	return reflect.ValueOf(bool(a.AssertBool(v)))
}

func strToNative(v a.Value) reflect.Value {
	return reflect.ValueOf(string(a.AssertStr(v)))
}

func numberToFloat32(v a.Value) reflect.Value {
	f, _ := a.AssertNumber(v).Float64()
	return reflect.ValueOf(float32(f))
}

func numberToFloat64(v a.Value) reflect.Value {
	f, _ := a.AssertNumber(v).Float64()
	return reflect.ValueOf(f)
}

func numberToInt(v a.Value) reflect.Value {
	i := a.AssertInteger(v)
	return reflect.ValueOf(i)
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
	convertIn = map[reflect.Kind]inMapper{
		reflect.Bool:    boolToNative,
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
		reflect.String:  strToNative,
	}

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
