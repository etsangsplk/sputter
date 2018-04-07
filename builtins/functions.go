package builtins

import (
	"strings"

	a "github.com/kode4food/sputter/api"
)

const (
	// InvalidRestArgument is thrown if you include more than one rest argument
	InvalidRestArgument = "rest-argument not well-formed: %s"

	// ExpectedArguments is thrown if argument patterns don't match
	ExpectedArguments = "expected arguments of the form: %s"

	lambdaName        = "fn"
	isSpecialFormName = "is-special-form"
)

type (
	lambdaFunction        struct{ BaseBuiltIn }
	isSpecialFormFunction struct{ BaseBuiltIn }

	functionDefinition struct {
		name       a.Name
		signatures funcSignatures
		meta       a.Object
	}

	funcSignature struct {
		args a.Vector
		body a.Sequence
	}

	funcSignatures []funcSignature

	singleFunction struct {
		a.BaseFunction
		args    string
		variant funcVariant
	}

	multiFunction struct {
		a.BaseFunction
		args     string
		variants funcVariants
	}

	funcVariant  func(a.Context, a.Vector) (a.Value, bool)
	funcVariants []funcVariant
)

var (
	emptyMetadata = a.Properties{}
	restMarker    = a.Name("&")
)

func (f *singleFunction) Apply(c a.Context, args a.Vector) a.Value {
	if r, ok := f.variant(c, args); ok {
		return r
	}
	panic(a.ErrStr(ExpectedArguments, f.args))
}

func (f *singleFunction) WithMetadata(md a.Object) a.AnnotatedValue {
	return &singleFunction{
		BaseFunction: f.Extend(md),
		args:         f.args,
		variant:      f.variant,
	}
}

func (f *multiFunction) Apply(c a.Context, args a.Vector) a.Value {
	for _, p := range f.variants {
		if r, ok := p(c, args); ok {
			return r
		}
	}
	panic(a.ErrStr(ExpectedArguments, f.args))
}

func (f *multiFunction) WithMetadata(md a.Object) a.AnnotatedValue {
	return &multiFunction{
		BaseFunction: f.Extend(md),
		args:         f.args,
		variants:     f.variants,
	}
}

func makeVariant(cl a.Context, s funcSignature) funcVariant {
	var an a.Names
	for f, r, ok := s.args.Split(); ok; f, r, ok = r.Split() {
		n := f.(a.LocalSymbol).Name()
		if n == restMarker {
			rn := parseRestArg(r)
			return makeRestVariant(cl, s, append(an, rn))
		}
		an = append(an, n)
	}
	return makeFixedVariant(cl, s, an)
}

func parseRestArg(s a.Sequence) a.Name {
	if f, r, ok := s.Split(); ok {
		n := f.(a.LocalSymbol).Name()
		if n != restMarker && !r.IsSequence() {
			return n
		}
	}
	panic(a.ErrStr(InvalidRestArgument, s))
}

func makeRestVariant(cl a.Context, s funcSignature, an a.Names) funcVariant {
	ex := a.MacroExpandAll(cl, s.body).(a.Sequence)
	body := a.MakeBlock(ex)

	nl := len(an)
	al := nl - 1

	rn := an[al]
	an = an[0:al]

	return func(_ a.Context, args a.Vector) (a.Value, bool) {
		if c := len(args); c < al {
			return nil, false
		}
		l := make(a.Variables, nl)
		for i, n := range an {
			l[n] = args[i]
		}
		// TODO: Investigate why we need a List
		l[rn] = a.NewList(args[len(an):]...)
		return body.Eval(a.ChildContext(cl, l)), true
	}
}

func makeFixedVariant(cl a.Context, s funcSignature, an a.Names) funcVariant {
	ex := a.MacroExpandAll(cl, s.body).(a.Sequence)
	body := a.MakeBlock(ex)
	al := len(an)

	return func(_ a.Context, args a.Vector) (a.Value, bool) {
		if c := len(args); c != al {
			return nil, false
		}
		l := make(a.Variables, al)
		for i, n := range an {
			l[n] = args[i]
		}
		return body.Eval(a.ChildContext(cl, l)), true
	}
}

func optionalMetadata(args a.Vector) (a.Object, a.Vector) {
	r := args
	var md a.Object
	if s, ok := r[0].(a.Str); ok {
		md = a.Properties{a.DocKey: s}
		r = r[1:]
	} else {
		md = emptyMetadata
	}

	if m, ok := r[0].(a.MappedSequence); ok {
		md = md.Child(toProperties(m))
		r = r[1:]
	}
	return md, r
}

func optionalName(args a.Vector) (a.Name, a.Vector) {
	if s, ok := args[0].(a.Symbol); ok {
		ls := s.(a.LocalSymbol)
		return ls.Name(), args[1:]
	}
	return a.DefaultFunctionName, args
}

func parseNamedFunction(args a.Vector) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	fn := args[0].(a.LocalSymbol).Name()
	return parseFunctionRest(fn, args[1:])
}

func parseFunction(args a.Vector) *functionDefinition {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	return parseFunctionRest(fn, r)
}

func parseFunctionRest(fn a.Name, r a.Vector) *functionDefinition {
	md, r := optionalMetadata(r)
	s := parseFunctionSignatures(r)
	md = md.Child(a.Properties{
		a.NameKey: fn,
	})

	return &functionDefinition{
		name:       fn,
		signatures: s,
		meta:       md,
	}
}

func parseFunctionSignatures(s a.Sequence) funcSignatures {
	f, r, ok := s.Split()
	if ar, ok := f.(a.Vector); ok {
		return funcSignatures{
			{args: ar, body: r},
		}
	}
	res := funcSignatures{}
	for ; ok; f, r, ok = r.Split() {
		lf, lr, _ := f.(*a.List).Split()
		res = append(res, funcSignature{
			args: lf.(a.Vector),
			body: lr,
		})
	}
	return res
}

func makeFunction(c a.Context, d *functionDefinition) a.Function {
	if len(d.signatures) > 1 {
		return makeMultiFunction(c, d)
	}
	return makeSingleFunction(c, d)
}

func makeSingleFunction(c a.Context, d *functionDefinition) a.Function {
	s := d.signatures[0]

	f := &singleFunction{
		BaseFunction: a.DefaultBaseFunction.Extend(d.meta),
		args:         string(s.args.Str()),
	}

	closure := a.ChildContext(c, a.Variables{
		d.name: f,
	})

	f.variant = makeVariant(closure, s)
	return f
}

func makeMultiFunction(c a.Context, d *functionDefinition) a.Function {
	s := d.signatures
	ls := len(s)

	f := &multiFunction{
		BaseFunction: a.DefaultBaseFunction.Extend(d.meta),
		variants:     make(funcVariants, ls),
	}

	closure := a.ChildContext(c, a.Variables{
		d.name: f,
	})

	ar := make([]string, ls)
	for i, s := range s {
		f.variants[i] = makeVariant(closure, s)
		ar[i] = string(s.args.Str())
	}

	f.args = strings.Join(ar, " or ")
	return f
}

func (*lambdaFunction) Apply(c a.Context, args a.Vector) a.Value {
	fd := parseFunction(args)
	return makeFunction(c, fd)
}

func (*isSpecialFormFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if ap, ok := args[0].(a.Applicable); ok && a.IsSpecialForm(ap) {
		return a.True
	}
	return a.False
}

func init() {
	var lambda *lambdaFunction
	var isSpecialForm *isSpecialFormFunction

	RegisterBuiltIn(lambdaName, lambda)
	RegisterBuiltIn(isSpecialFormName, isSpecialForm)
}
