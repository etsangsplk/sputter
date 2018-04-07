package builtins

import a "github.com/kode4food/sputter/api"

const (
	lazySequenceName = "lazy-seq*"
	concatName       = "concat"
	filterName       = "filter"
	mapName          = "map"
	reduceName       = "reduce"
	takeName         = "take"
	dropName         = "drop"
	rangeName        = "range*"
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
		r := f.Apply(c, a.EmptyVector)
		if r != a.Nil {
			s := r.(a.Sequence)
			if sf, sr, ok := s.Split(); ok {
				return sf, sr, true
			}
		}
		return a.Nil, a.EmptyList, false
	}
}

func (*lazySequenceFunction) Apply(c a.Context, args a.Vector) a.Value {
	fn := a.NewBlockFunction(args)
	return a.NewLazySequence(makeLazyResolver(c, fn))
}

func (*concatFunction) Apply(_ a.Context, args a.Vector) a.Value {
	if a.AssertMinimumArity(args, 1) == 1 {
		return args[0].(a.Sequence)
	}
	return a.Concat(args)
}

func (*filterFunction) Apply(c a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 2)
	fn := args[0].(a.Applicable)
	s := args[1].(a.Sequence)
	return a.Filter(c, s, fn)
}

func (*mapFunction) Apply(c a.Context, args a.Vector) a.Value {
	cnt := a.AssertMinimumArity(args, 2)
	fn := args[0].(a.Applicable)
	if cnt == 2 {
		s := args[1].(a.Sequence)
		return a.Map(c, s, fn)
	}
	return a.MapParallel(c, args[1:], fn)
}

func (*reduceFunction) Apply(c a.Context, args a.Vector) a.Value {
	cnt := a.AssertArityRange(args, 2, 3)
	fn := args[0].(a.Applicable)
	if cnt == 2 {
		s := args[1].(a.Sequence)
		return a.Reduce(c, s, fn)
	}
	s := args[2].(a.Sequence).Prepend(args[1])
	return a.Reduce(c, s, fn)
}

func (*takeFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 2)
	n := a.AssertInteger(args[0])
	s := args[1].(a.Sequence)
	return a.Take(s, n)
}

func (*dropFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 2)
	n := a.AssertInteger(args[0])
	s := args[1].(a.Sequence)
	return a.Drop(s, n)
}

func (*rangeFunction) Apply(_ a.Context, args a.Vector) a.Value {
	a.AssertArity(args, 3)
	low := args[0].(a.Number)
	high := args[1].(a.Number)
	step := args[2].(a.Number)
	return a.NewRange(low, high, step)
}

func (*forEachFunction) Apply(c a.Context, args a.Vector) a.Value {
	a.AssertMinimumArity(args, 2)
	b := args[0].(a.Vector)
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
			proc = makeTerminal(n, e, args[1:])
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
			l := a.ChildContext(c, a.Variables{n: f})
			next(l)
		}
	}
}

func makeTerminal(n a.Name, e a.Value, s a.Sequence) forProc {
	bl := a.MakeBlock(s)
	return func(c a.Context) {
		es := a.Eval(c, e).(a.Sequence)
		for f, r, ok := es.Split(); ok; f, r, ok = r.Split() {
			l := a.ChildContext(c, a.Variables{n: f})
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
