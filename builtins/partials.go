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

func (f *boundFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fullArgs := f.args.Concat(a.SequenceToValues(args))
	return f.bound.Apply(c, fullArgs)
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
		return &boundFunction{
			BaseFunction: a.DefaultBaseFunction,
			bound:        bf.bound,
			args:         bf.args.Concat(values),
		}
	}
	return &boundFunction{
		BaseFunction: a.DefaultBaseFunction,
		bound:        bound,
		args:         values,
	}
}

func init() {
	var partial *partialFunction

	RegisterBuiltIn(partialName, partial)
}
