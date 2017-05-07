package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
)

type methodInfo struct {
	idx int
	in  argSetters
	out resultGetters
}

type (
	argumentReader func(c a.Context, args a.Sequence) []reflect.Value
	inMapper       func(a.Value) reflect.Value
	resultGetters  []outMapper
	argSetters     []inMapper
)

var convertIn map[reflect.Kind]inMapper

func makeMethodGetters(t reflect.Type) propertyGetters {
	r := propertyGetters{}
	for i := 0; i < t.NumMethod(); i++ {
		mi := t.Method(i)
		if mi.PkgPath != "" {
			continue // only surface exported methods
		}
		n := kebabCase(mi.Name)
		r[n] = makeMethodInvoker(mi)
	}
	return r
}

func makeMethodInvoker(m reflect.Method) outMapper {
	mi := makeMethodInfo(m)
	return mi.makeInvoker()
}

func makeMethodInfo(m reflect.Method) *methodInfo {
	ft := m.Func.Type()

	in := argSetters{}
	for i := 1; i < ft.NumIn(); i++ {
		in = append(in, makeArgSetter(ft.In(i)))
	}

	out := resultGetters{}
	for i := 0; i < ft.NumOut(); i++ {
		t := ft.Out(i)
		if c, ok := convertOut[t.Kind()]; ok {
			out = append(out, c)
		}
	}

	return &methodInfo{
		idx: m.Index,
		in:  in,
		out: out,
	}
}

func (mi *methodInfo) makeInvoker() outMapper {
	olen := len(mi.out)
	if olen == 0 {
		return mi.makeVoidInvoker()
	} else if olen == 1 {
		return mi.makeSingularInvoker()
	} else {
		return mi.makePluralInvoker()
	}
}

func (mi *methodInfo) makeArgPreparer() argumentReader {
	in := mi.in
	ilen := len(in)

	return func(c a.Context, args a.Sequence) []reflect.Value {
		a.AssertArity(args, ilen)
		fin := make([]reflect.Value, ilen)
		e := args
		for i := 0; i < ilen; i++ {
			fin[i] = in[i](a.Eval(c, e.First()))
			e = e.Rest()
		}
		return fin
	}
}

func (mi *methodInfo) makeVoidInvoker() outMapper {
	prepareArgs := mi.makeArgPreparer()

	return func(v reflect.Value) a.Value {
		fn := v.Method(mi.idx)
		return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			fin := prepareArgs(c, args)
			fn.Call(fin)
			return a.Nil
		})
	}
}

func (mi *methodInfo) makeSingularInvoker() outMapper {
	prepareArgs := mi.makeArgPreparer()
	out := mi.out

	return func(v reflect.Value) a.Value {
		fn := v.Method(mi.idx)
		return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			fin := prepareArgs(c, args)
			r := fn.Call(fin)
			return out[0](r[0])
		})
	}
}

func (mi *methodInfo) makePluralInvoker() outMapper {
	prepareArgs := mi.makeArgPreparer()
	out := mi.out
	olen := len(out)

	return func(v reflect.Value) a.Value {
		fn := v.Method(mi.idx)
		return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			fin := prepareArgs(c, args)
			r := fn.Call(fin)
			fout := make(a.Vector, olen)
			for i := 0; i < olen; i++ {
				fout[i] = out[i](r[i])
			}
			return fout
		})
	}
}

func makeArgSetter(t reflect.Type) inMapper {
	if c, ok := convertIn[t.Kind()]; ok {
		return c
	}

	return func(_ a.Value) reflect.Value {
		panic(a.Err(BadConversionType, t))
	}
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
}

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
