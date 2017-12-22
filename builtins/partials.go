package builtins

import a "github.com/kode4food/sputter/api"

const partialName = "partial"

type (
	partialFunction struct{ BaseBuiltIn }

	boundFunction struct {
		a.BaseFunction
		bound a.Applicable
		args  a.Values
	}
)

// BoundArgsKey is the Metadata key for a Function's bound count
var BoundArgsKey = a.NewKeyword("bound-args")

func bindFunction(bound a.Applicable, args a.Values) *boundFunction {
	var md a.Object
	if an, ok := bound.(a.Annotated); ok {
		md = an.Metadata()
	} else {
		md = a.DefaultBaseFunction.Metadata()
	}

	md = md.Child(a.Properties{
		BoundArgsKey: a.NewFloat(float64(len(args))),
	})

	return &boundFunction{
		BaseFunction: a.DefaultBaseFunction.Extend(md),
		bound:        bound,
		args:         args,
	}
}

func (f *boundFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fullArgs := f.args.Concat(a.SequenceToValues(args))
	return f.bound.Apply(c, fullArgs)
}

func (f *boundFunction) rebind(args a.Values) *boundFunction {
	newArgs := f.args.Concat(args)
	md := a.Properties{
		BoundArgsKey: a.NewFloat(float64(len(newArgs))),
	}
	return &boundFunction{
		BaseFunction: f.BaseFunction.Extend(md),
		bound:        f.bound,
		args:         newArgs,
	}
}

func (f *boundFunction) WithMetadata(md a.Object) a.AnnotatedValue {
	return &boundFunction{
		BaseFunction: f.Extend(md),
		bound:        f.bound,
		args:         f.args,
	}
}

func (*partialFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 1)
	bound := args.First().(a.Applicable)
	values := a.SequenceToValues(args.Rest())

	if bf, ok := bound.(*boundFunction); ok {
		return bf.rebind(values)
	}
	return bindFunction(bound, values)
}

func init() {
	var partial *partialFunction

	RegisterBuiltIn(partialName, partial)
}