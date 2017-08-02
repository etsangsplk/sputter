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
		BaseFunction: a.NewBaseFunction(f.Metadata().Child(md.Flatten())),
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
		BaseFunction: a.NewBaseFunction(f.Metadata().Child(md.Flatten())),
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
		BaseFunction: a.NewBaseFunction(f.Metadata().Child(md.Flatten())),
		body:         f.body,
	}
}

func makeArgProcessor(cl a.Context, s a.Sequence) argProcessor {
	an := []a.Name{}
	for i := s; i.IsSequence(); i = i.Rest() {
		n := a.AssertUnqualified(i.First()).Name()
		if n == restMarker {
			rn := parseRestArg(i)
			return makeRestArgProcessor(cl, an, rn)
		}
		an = append(an, n)
	}
	return makeFixedArgProcessor(cl, an)
}

func parseRestArg(s a.Sequence) a.Name {
	r := s.Rest()
	if r.IsSequence() {
		n := a.AssertUnqualified(r.First()).Name()
		if n != restMarker && !r.Rest().IsSequence() {
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
		for _, n := range an {
			l.Put(n, i.First())
			i = i.Rest()
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
		for _, n := range an {
			l.Put(n, i.First())
			i = i.Rest()
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
	f := args.First()
	if s, ok := f.(a.Symbol); ok {
		if s.Domain() == a.LocalDomain {
			return s.Name(), args.Rest()
		}
		panic(a.ErrStr(a.ExpectedUnqualified, s.Qualified()))
	}
	return a.DefaultFunctionName, args
}

func parseNamedFunction(args a.Sequence) *functionDefinition {
	a.AssertMinimumArity(args, 3)
	fn := a.AssertUnqualified(args.First()).Name()
	return parseFunctionRest(fn, args.Rest())
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

func parseFunctionSignatures(r a.Sequence) functionSignatures {
	if args, ok := r.First().(a.Vector); ok {
		return functionSignatures{
			{args: args, body: r.Rest()},
		}
	}
	res := functionSignatures{}
	for i := r; i.IsSequence(); i = i.Rest() {
		l := a.AssertList(i.First())
		res = append(res, &functionSignature{
			args: a.AssertVector(l.First()),
			body: l.Rest(),
		})
	}
	return res
}

func makeFunction(c a.Context, d *functionDefinition) a.Function {
	var res a.Function

	if len(d.sigs) > 1 {
		res = makeMultiFunction(c, d.sigs)
	} else {
		res = makeSingleFunction(c, d.sigs[0])
	}
	return res.WithMetadata(d.meta).(a.Function)
}

func makeSingleFunction(c a.Context, s *functionSignature) a.Function {
	ex := a.MacroExpandAll(c, s.body).(a.Sequence)

	return &singleFunction{
		BaseFunction: a.DefaultBaseFunction,
		argProcessor: makeArgProcessor(c, s.args),
		args:         string(s.args.Str()),
		body:         a.MakeBlock(ex),
	}
}

func makeMultiFunction(c a.Context, sigs functionSignatures) a.Function {
	ls := len(sigs)
	matchers := make([]argProcessorMatch, ls)
	args := make([]string, ls)

	for i, s := range sigs {
		ex := a.MacroExpandAll(c, s.body).(a.Sequence)
		matchers[i] = argProcessorMatch{
			args: makeArgProcessor(c, s.args),
			body: a.MakeBlock(ex),
		}
		args[i] = string(s.args.Str())
	}

	return &multiFunction{
		BaseFunction: a.DefaultBaseFunction,
		args:         strings.Join(args, " or "),
		matchers:     matchers,
	}
}

func (f *lambdaFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fd := parseFunction(args)
	return makeFunction(c, fd)
}

func (f *applyFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 2)
	fn := a.AssertApplicable(args.First())
	s := a.AssertSequence(args.Rest().First())
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

	RegisterBuiltIn("lambda", lambda)
	RegisterBuiltIn("apply", apply)

	RegisterSequencePredicate("apply?", isApplicable)
	RegisterSequencePredicate("special-form?", isSpecialForm)
}
