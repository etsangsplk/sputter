package builtins

import (
	a "github.com/kode4food/sputter/api"
)

const partialName = "partial"

type (
	partialFunction struct{ BaseBuiltIn }

	boundFunction struct {
		a.BaseFunction
		bound a.Applicable
		args  a.Values
	}
)

// BoundKey is the Metadata key for a bound Function
var BoundKey = a.NewKeyword("bound")

// BindFunction binds a set of arguments to an Applicable
func BindFunction(bound a.Applicable, args a.Values) *boundFunction {
	var md a.Object
	if an, ok := bound.(a.Annotated); ok {
		md = an.Metadata()
	} else {
		md = a.DefaultBaseFunction.Metadata()
	}

	md = md.Child(a.Properties{
		BoundKey: a.True,
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

func (f *boundFunction) Rebind(args a.Values) *boundFunction {
	return &boundFunction{
		BaseFunction: f.BaseFunction,
		bound:        f.bound,
		args:         f.args.Concat(args),
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
		return bf.Rebind(values)
	} else {
		return BindFunction(bound, values)
	}
}

func init() {
	var partial *partialFunction

	RegisterBuiltIn(partialName, partial)
}
