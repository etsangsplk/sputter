package native

import (
	"reflect"

	a "github.com/kode4food/sputter/api"
	u "github.com/kode4food/sputter/util"
)

type funcInfo struct {
	fn  funcResolver
	in  argSetters
	out resultGetters
}

type (
	funcResolver   func(v reflect.Value) reflect.Value
	argumentReader func(c a.Context, args a.Sequence) []reflect.Value
	resultGetters  []outMapper
	argSetters     []inMapper
)

var funcs = u.NewCache()

func funcToApplicable(v reflect.Value) a.Value {
	ft := v.Type()

	om := funcs.Get(ft, func() u.Any {
		fi := makeFuncInfo(ft)
		return fi.makeFuncMapper()
	}).(outMapper)

	return om(v)
}

func makeFuncInfo(ft reflect.Type) funcInfo {
	in := makeFuncArgSetters(ft)
	out := makeResultGetters(ft)

	return funcInfo{
		fn:  defaultFuncResolver,
		in:  in,
		out: out,
	}
}

func defaultFuncResolver(v reflect.Value) reflect.Value {
	return v
}

func makeFuncArgSetters(ft reflect.Type) argSetters {
	return makeArgSetters(ft, 0)
}

func makeArgSetters(ft reflect.Type, startAt int) argSetters {
	r := argSetters{}
	for i := startAt; i < ft.NumIn(); i++ {
		r = append(r, getConvertIn(ft.In(i)))
	}
	return r
}

func makeResultGetters(ft reflect.Type) resultGetters {
	r := resultGetters{}
	for i := 0; i < ft.NumOut(); i++ {
		t := ft.Out(i)
		c := getConvertOut(t)
		r = append(r, c)
	}
	return r
}

func (fi funcInfo) makeFuncMapper() outMapper {
	olen := len(fi.out)
	if olen == 0 {
		return fi.makeVoidFuncMapper()
	} else if olen == 1 {
		return fi.makeSingularFuncMapper()
	} else {
		return fi.makePluralFuncMapper()
	}
}

func (fi funcInfo) makeVoidFuncMapper() outMapper {
	prepareArgs := fi.makeArgPreparer()

	return func(v reflect.Value) a.Value {
		fn := fi.fn(v)
		return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			fin := prepareArgs(c, args)
			fn.Call(fin)
			return a.Nil
		})
	}
}

func (fi funcInfo) makeSingularFuncMapper() outMapper {
	prepareArgs := fi.makeArgPreparer()
	out := fi.out

	return func(v reflect.Value) a.Value {
		fn := fi.fn(v)
		return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			fin := prepareArgs(c, args)
			r := fn.Call(fin)
			return out[0](r[0])
		})
	}
}

func (fi funcInfo) makePluralFuncMapper() outMapper {
	prepareArgs := fi.makeArgPreparer()
	out := fi.out
	olen := len(out)

	return func(v reflect.Value) a.Value {
		fn := fi.fn(v)
		return a.NewFunction(func(c a.Context, args a.Sequence) a.Value {
			fin := prepareArgs(c, args)
			r := fn.Call(fin)
			fout := make([]a.Value, olen)
			for i := 0; i < olen; i++ {
				fout[i] = out[i](r[i])
			}
			return a.NewVector(fout...)
		})
	}
}

func (fi funcInfo) makeArgPreparer() argumentReader {
	in := fi.in
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
