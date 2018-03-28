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

	lambdaName        = "lambda"
	isSpecialFormName = "is-special-form"
)

type (
	lambdaFunction        struct{ BaseBuiltIn }
	isSpecialFormFunction struct{ BaseBuiltIn }

	argProcessor func(a.Context, a.Values) (a.Context, bool)

	functionSignature struct {
		args a.Vector
		body a.Sequence
	}

	functionSignatures []*functionSignature

	functionDefinition struct {
		name a.Name
		sigs functionSignatures
		meta a.Object
	}

	argProcessorMatch struct {
		args argProcessor
		body a.Block
	}

	singleFunction struct {
		a.BaseFunction
		args         string
		argProcessor argProcessor
		body         a.Block
	}

	multiFunction struct {
		a.BaseFunction
		args     string
		matchers []argProcessorMatch
	}
)

var (
	emptyMetadata = a.Properties{}
	restMarker    = a.Name("&")
)

func (f *singleFunction) Apply(c a.Context, args a.Values) a.Value {
	if l, ok := f.argProcessor(c, args); ok {
		return f.body.Eval(l)
	}
	panic(a.ErrStr(ExpectedArguments, f.args))
}

func (f *singleFunction) WithMetadata(md a.Object) a.AnnotatedValue {
	return &singleFunction{
		BaseFunction: f.Extend(md),
		args:         f.args,
		argProcessor: f.argProcessor,
		body:         f.body,
	}
}

func (f *multiFunction) Apply(c a.Context, args a.Values) a.Value {
	for _, m := range f.matchers {
		if l, ok := m.args(c, args); ok {
			return m.body.Eval(l)
		}
	}
	panic(a.ErrStr(ExpectedArguments, f.args))
}

func (f *multiFunction) WithMetadata(md a.Object) a.AnnotatedValue {
	return &multiFunction{
		BaseFunction: f.Extend(md),
		args:         f.args,
		matchers:     f.matchers,
	}
}

func makeArgProcessor(cl a.Context, s a.Sequence) argProcessor {
	var an []a.Name
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		n := f.(a.LocalSymbol).Name()
		if n == restMarker {
			rn := parseRestArg(r)
			return makeRestArgProcessor(cl, an, rn)
		}
		an = append(an, n)
	}
	return makeFixedArgProcessor(cl, an)
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

func makeRestArgProcessor(cl a.Context, an []a.Name, rn a.Name) argProcessor {
	ac := a.MakeMinimumArityChecker(len(an))

	return func(_ a.Context, args a.Values) (a.Context, bool) {
		if _, ok := ac(args); !ok {
			return nil, false
		}
		l := a.ChildLocals(cl)
		for i, n := range an {
			l.Put(n, args[i])
		}
		l.Put(rn, a.SequenceToList(args[len(an):]))
		return l, true
	}
}

func makeFixedArgProcessor(cl a.Context, an []a.Name) argProcessor {
	ac := a.MakeArityChecker(len(an))

	return func(_ a.Context, args a.Values) (a.Context, bool) {
		if _, ok := ac(args); !ok {
			return nil, false
		}
		l := a.ChildLocals(cl)
		for i, n := range an {
			l.Put(n, args[i])
		}
		return l, true
	}
}

func optionalMetadata(args a.Values) (a.Object, a.Values) {
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

func optionalName(args a.Values) (a.Name, a.Values) {
	if s, ok := args[0].(a.Symbol); ok {
		ls := s.(a.LocalSymbol)
		return ls.Name(), args[1:]
	}
	return a.DefaultFunctionName, args
}

func parseNamedFunction(args a.Values) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	fn := args[0].(a.LocalSymbol).Name()
	return parseFunctionRest(fn, args[1:])
}

func parseFunction(args a.Values) *functionDefinition {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	return parseFunctionRest(fn, r)
}

func parseFunctionRest(fn a.Name, r a.Values) *functionDefinition {
	md, r := optionalMetadata(r)
	sigs := parseFunctionSignatures(r)
	md = md.Child(a.Properties{
		a.NameKey: fn,
	})

	return &functionDefinition{
		name: fn,
		sigs: sigs,
		meta: md,
	}
}

func parseFunctionSignatures(s a.Sequence) functionSignatures {
	f, r, ok := s.Split()
	if ar, ok := f.(a.Vector); ok {
		return functionSignatures{
			{args: ar, body: r},
		}
	}
	res := functionSignatures{}
	for ; ok; f, r, ok = r.Split() {
		lf, lr, _ := f.(a.List).Split()
		res = append(res, &functionSignature{
			args: lf.(a.Vector),
			body: lr,
		})
	}
	return res
}

func makeFunction(c a.Context, d *functionDefinition) a.Function {
	if len(d.sigs) > 1 {
		return makeMultiFunction(c, d)
	}
	return makeSingleFunction(c, d)
}

func makeSingleFunction(c a.Context, d *functionDefinition) a.Function {
	s := d.sigs[0]

	f := &singleFunction{
		BaseFunction: a.DefaultBaseFunction.Extend(d.meta),
		args:         string(s.args.Str()),
	}

	nc := a.ChildContext(c, a.Variables{
		d.name: f,
	})

	ex := a.MacroExpandAll(nc, s.body).(a.Sequence)
	f.argProcessor = makeArgProcessor(nc, s.args)
	f.body = a.MakeBlock(ex)
	return f
}

func makeMultiFunction(c a.Context, d *functionDefinition) a.Function {
	sigs := d.sigs
	ls := len(sigs)

	f := &multiFunction{
		BaseFunction: a.DefaultBaseFunction.Extend(d.meta),
	}

	nc := a.ChildContext(c, a.Variables{
		d.name: f,
	})

	matchers := make([]argProcessorMatch, ls)
	ar := make([]string, ls)
	for i, s := range sigs {
		ex := a.MacroExpandAll(nc, s.body).(a.Sequence)
		matchers[i] = argProcessorMatch{
			args: makeArgProcessor(nc, s.args),
			body: a.MakeBlock(ex),
		}
		ar[i] = string(s.args.Str())
	}

	f.matchers = matchers
	f.args = strings.Join(ar, " or ")
	return f
}

func (*lambdaFunction) Apply(c a.Context, args a.Values) a.Value {
	fd := parseFunction(args)
	return makeFunction(c, fd)
}

func (*isSpecialFormFunction) Apply(_ a.Context, args a.Values) a.Value {
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
