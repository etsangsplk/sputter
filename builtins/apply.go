package builtins

import a "github.com/kode4food/sputter/api"

const (
	applyName   = "apply"
	partialName = "partial"
	isApplyName = "is-apply"
)

type (
	applyFunction   struct{ BaseBuiltIn }
	partialFunction struct{ BaseBuiltIn }
	isApplyFunction struct{ BaseBuiltIn }

	boundFunction struct {
		a.BaseFunction
		bound a.Applicable
		args  a.Vector
	}
)

// BoundArgsKey is the Metadata key for a Function's bound count
var BoundArgsKey = a.NewKeyword("bound-args")

func (*applyFunction) Apply(c a.Context, args a.Vector) a.Value {
	ac := a.AssertMinimumArity(args, 2)
	fn := args[0].(a.Applicable)
	if ac == 2 {
		return fn.Apply(c, a.SequenceToVector(args[1].(a.Sequence)))
	}
	last := len(args) - 1
	ls := a.SequenceToVector(args[last].(a.Sequence))
	prependedArgs := append(args[1:last], ls...)
	return fn.Apply(c, prependedArgs)
}

func (*partialFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 1)
	bound := args[0].(a.Applicable)
	values := args[1:]

	if bf, ok := bound.(*boundFunction); ok {
		return bf.rebind(values)
	}
	return bindFunction(bound, values)
}

func (*isApplyFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if _, ok := args[0].(a.Applicable); ok {
		return a.True
	}
	return a.False
}

func bindFunction(bound a.Applicable, args a.Vector) *boundFunction {
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

func (f *boundFunction) Apply(c a.Context, args a.Vector) a.Value {
	fullArgs := f.args.Concat(a.SequenceToVector(args))
	return f.bound.Apply(c, fullArgs)
}

func (f *boundFunction) rebind(args a.Vector) *boundFunction {
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

func init() {
	var apply *applyFunction
	var partial *partialFunction
	var isApply *isApplyFunction

	RegisterBuiltIn(applyName, apply)
	RegisterBuiltIn(partialName, partial)
	RegisterBuiltIn(isApplyName, isApply)
}
