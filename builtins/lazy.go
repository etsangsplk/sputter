package builtins

import a "github.com/kode4food/sputter/api"

const (
	lazySequenceName = "make-lazy-seq"
	concatName       = "concat"
	filterName       = "filter"
	mapName          = "map"
	reduceName       = "reduce"
	takeName         = "take"
	dropName         = "drop"
	rangeName        = "make-range"
	forEachName      = "for-each"
)

type (
	lazySequenceFunction struct{ BaseBuiltIn }
	concatFunction       struct{ BaseBuiltIn }
	filterFunction       struct{ BaseBuiltIn }
	mapFunction          struct{ BaseBuiltIn }
	reduceFunction       struct{ BaseBuiltIn }
	takeFunction         struct{ BaseBuiltIn }
	dropFunction         struct{ BaseBuiltIn }
	rangeFunction        struct{ BaseBuiltIn }
	forEachFunction      struct{ BaseBuiltIn }
)

func makeLazyResolver(c a.Context, f a.Applicable) a.LazyResolver {
	return func() (a.Value, a.Sequence, bool) {
		r := f.Apply(c, a.EmptyList)
		if r != a.Nil {
			s := r.(a.Sequence)
			if sf, sr, ok := s.Split(); ok {
				return sf, sr, true
			}
		}
		return a.Nil, a.EmptyList, false
	}
}

func (*lazySequenceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	fn := a.NewBlockFunction(args)
	return a.NewLazySequence(makeLazyResolver(c, fn))
}

func (*concatFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	if a.AssertMinimumArity(args, 1) == 1 {
		return args.First().(a.Sequence)
	}
	return a.Concat(args)
}

func (*filterFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	fn := f.(a.Applicable)
	s := a.Concat(r)
	return a.Filter(c, s, fn)
}

func (*mapFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	fn := f.(a.Applicable)
	s := a.Concat(r)
	return a.Map(c, s, fn)
}

func (*reduceFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	fn := f.(a.Applicable)
	s := a.Concat(r)
	return a.Reduce(c, s, fn)
}

func (*takeFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	n := a.AssertInteger(f)
	s := a.Concat(r)
	return a.Take(s, n)
}

func (*dropFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	n := a.AssertInteger(f)
	s := a.Concat(r)
	return a.Drop(s, n)
}

func (*rangeFunction) Apply(_ a.Context, args a.Sequence) a.Value {
	a.AssertArity(args, 3)
	f, r, _ := args.Split()
	low := f.(a.Number)
	rf, rr, _ := r.Split()
	high := rf.(a.Number)
	step := rr.First().(a.Number)
	return a.NewRange(low, high, step)
}

func (*forEachFunction) Apply(c a.Context, args a.Sequence) a.Value {
	a.AssertMinimumArity(args, 2)
	f, r, _ := args.Split()
	b := f.(a.Vector)
	bc := b.Count()
	if bc%2 != 0 {
		panic(a.ErrStr(ExpectedBindings))
	}

	var proc forProc
	depth := bc / 2
	for i := depth - 1; i >= 0; i-- {
		o := i * 2
		s, _ := b.ElementAt(o)
		e, _ := b.ElementAt(o + 1)
		n := s.(a.LocalSymbol).Name()
		if i == depth-1 {
			proc = makeTerminal(n, e, r)
		} else {
			proc = makeIntermediate(n, e, proc)
		}
	}

	proc(c)
	return a.Nil
}

func makeIntermediate(n a.Name, e a.Value, next forProc) forProc {
	return func(c a.Context) {
		s := a.Eval(c, e).(a.Sequence)
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			l := a.ChildVariables(c, a.Variables{n: f})
			next(l)
		}
	}
}

func makeTerminal(n a.Name, e a.Value, s a.Sequence) forProc {
	bl := a.MakeBlock(s)
	return func(c a.Context) {
		es := a.Eval(c, e).(a.Sequence)
		for f, r, ok := es.Split(); ok; f, r, ok = r.Split() {
			l := a.ChildVariables(c, a.Variables{n: f})
			bl.Eval(l)
		}
	}
}

func init() {
	var lazySequence *lazySequenceFunction
	var concat *concatFunction
	var filter *filterFunction
	var _map *mapFunction
	var reduce *reduceFunction
	var take *takeFunction
	var drop *dropFunction
	var _range *rangeFunction
	var forEach *forEachFunction

	RegisterBuiltIn(lazySequenceName, lazySequence)
	RegisterBuiltIn(concatName, concat)
	RegisterBuiltIn(filterName, filter)
	RegisterBuiltIn(mapName, _map)
	RegisterBuiltIn(reduceName, reduce)
	RegisterBuiltIn(takeName, take)
	RegisterBuiltIn(dropName, drop)
	RegisterBuiltIn(rangeName, _range)
	RegisterBuiltIn(forEachName, forEach)
}
