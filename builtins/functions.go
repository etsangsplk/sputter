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
	applyName         = "apply"
	isApplicableName  = "apply?"
	isSpecialFormName = "special-form?"
)

type (
	lambdaFunction struct{ BaseBuiltIn }
	applyFunction  struct{ BaseBuiltIn }

	argProcessor func(a.Context, a.Sequence) (a.Context, bool)

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

	blockFunction struct {
		a.BaseFunction
		body a.Block
	}
)

var (
	emptyMetadata = a.Properties{}
	restMarker    = a.Name("&")
)

func (f *singleFunction) Apply(c a.Context, args a.Sequence) a.Value {
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

func (f *multiFunction) Apply(c a.Context, args a.Sequence) a.Value {
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

// NewBlockFunction creates a new simple Function based on a Block
func NewBlockFunction(args a.Sequence) a.Function {
	return &blockFunction{
		BaseFunction: a.DefaultBaseFunction,
		body:         a.MakeBlock(args),
	}
}

func (f *blockFunction) Apply(c a.Context, _ a.Sequence) a.Value {
	return f.body.Eval(c)
}

func (f *blockFunction) WithMetadata(md a.Object) a.AnnotatedValue {
	return &blockFunction{
		BaseFunction: f.Extend(md),
		body:         f.body,
	}
}

func makeArgProcessor(cl a.Context, s a.Sequence) argProcessor {
	an := []a.Name{}
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

	return func(_ a.Context, args a.Sequence) (a.Context, bool) {
		if _, ok := ac(args); !ok {
			return nil, false
		}
		l := a.ChildContext(cl)
		i := args
		var t a.Value
		for _, n := range an {
			t, i, _ = i.Split()
			l.Put(n, t)
		}
		l.Put(rn, a.SequenceToList(i))
		return l, true
	}
}

func makeFixedArgProcessor(cl a.Context, an []a.Name) argProcessor {
	ac := a.MakeArityChecker(len(an))

	return func(_ a.Context, args a.Sequence) (a.Context, bool) {
		if _, ok := ac(args); !ok {
			return nil, false
		}
		l := a.ChildContext(cl)
		i := args
		var t a.Value
		for _, n := range an {
			t, i, _ = i.Split()
			l.Put(n, t)
		}
		return l, true
	}
}

func optionalMetadata(args a.Sequence) (a.Object, a.Sequence) {
	r := args
	var md a.Object
	if s, ok := r.First().(a.Str); ok {
		md = a.Properties{a.DocKey: s}
		r = r.Rest()
	} else {
		md = emptyMetadata
	}

	if m, ok := r.First().(a.MappedSequence); ok {
		md = md.Child(toProperties(m))
		r = r.Rest()
	}
	return md, r
}

func optionalName(args a.Sequence) (a.Name, a.Sequence) {
	f, r, _ := args.Split()
	if s, ok := f.(a.Symbol); ok {
		ls := s.(a.LocalSymbol)
		return ls.Name(), r
	}
	return a.DefaultFunctionName, args
}

func parseNamedFunction(args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	f, r, _ := args.Split()
	fn := f.(a.LocalSymbol).Name()
	return parseFunctionRest(fn, r)
}

func parseFunction(args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 2)
	fn, r := optionalName(args)
	return parseFunctionRest(fn, r)
}

func parseFunctionRest(fn a.Name, r a.Sequence) *functionDefinition {
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

	nc := a.ChildContextVars(c, a.Variables{
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

	nc := a.ChildContextVars(c, a.Variables{
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

func (*lambdaFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fd := parseFunction(args)
	return makeFunction(c, fd)
}

func (*applyFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	f, r, _ := args.Split()
	fn := f.(a.Applicable)
	s := r.First().(a.Sequence)
	return fn.Apply(c, s)
}

func isApplicable(v a.Value) bool {
	if _, ok := v.(a.Applicable); ok {
		return true
	}
	return false
}

func isSpecialForm(v a.Value) bool {
	if ap, ok := v.(a.Applicable); ok {
		return a.IsSpecialForm(ap)
	}
	return false
}

func init() {
	var lambda *lambdaFunction
	var apply *applyFunction

	RegisterBuiltIn(lambdaName, lambda)
	RegisterBuiltIn(applyName, apply)

	RegisterSequencePredicate(isApplicableName, isApplicable)
	RegisterSequencePredicate(isSpecialFormName, isSpecialForm)
}
