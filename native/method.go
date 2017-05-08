package native

import "reflect"

func makeMethodGetters(t reflect.Type) propertyGetters {
	l := t.NumMethod()
	r := make(propertyGetters, l)
	for i := 0; i < l; i++ {
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
	return mi.makeFuncMapper()
}

func makeMethodInfo(m reflect.Method) funcInfo {
	ft := m.Func.Type()
	in := makeMethodArgSetters(ft)
	out := makeResultGetters(ft)

	midx := m.Index
	methodResolver := func(v reflect.Value) reflect.Value {
		return v.Method(midx)
	}

	return funcInfo{
		fn:  methodResolver,
		in:  in,
		out: out,
	}
}

func makeMethodArgSetters(mt reflect.Type) argSetters {
	return makeArgSetters(mt, 1)
}
