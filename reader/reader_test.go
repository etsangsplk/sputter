package reader_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
	_ "github.com/kode4food/sputter/builtins"
	r "github.com/kode4food/sputter/reader"
)

func s(s string) a.Str {
	return a.Str(s)
}

func f(f float64) a.Number {
	return a.NewFloat(f)
}

func TestCreateReader(t *testing.T) {
	as := assert.New(t)
	l := r.Scan("99")
	tr := r.Read(l)
	as.NotNil(tr)
}

func TestReadList(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(`(99 "hello" 55.12)`)
	tr := r.Read(l)
	v := tr.First()
	list, ok := v.(*a.List)
	as.True(ok)

	i := a.Iterate(list)
	val, ok := i.Next()
	as.True(ok)
	as.Number(99, val)

	val, ok = i.Next()
	as.True(ok)
	as.String("hello", val)

	val, ok = i.Next()
	as.True(ok)
	as.Number(55.12, val)

	_, ok = i.Next()
	as.False(ok)
}

func TestReadVector(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(`[99 "hello" 55.12]`)
	tr := r.Read(l)
	v := tr.First()
	vector, ok := v.(a.Vector)
	as.True(ok)

	res, ok := vector.ElementAt(0)
	as.True(ok)
	as.Number(99, res)

	res, ok = vector.ElementAt(1)
	as.True(ok)
	as.String("hello", res)

	res, ok = vector.ElementAt(2)
	as.True(ok)
	as.Number(55.120, res)
}

func TestReadMap(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(`{:name "blah" :age 99}`)
	tr := r.Read(l)
	v := tr.First()
	m, ok := v.(a.Associative)
	as.True(ok)
	as.Number(2, m.Count())
}

func TestReadNestedList(t *testing.T) {
	as := assert.New(t)
	l := r.Scan(`(99 ("hello" "there") 55.12)`)
	tr := r.Read(l)
	v := tr.First()
	list, ok := v.(*a.List)
	as.True(ok)

	i1 := a.Iterate(list)
	val, ok := i1.Next()
	as.True(ok)
	as.Number(99, val)

	// get nested list
	val, ok = i1.Next()
	as.True(ok)
	list2, ok := val.(*a.List)
	as.True(ok)

	// iterate over the rest of top-level list
	val, ok = i1.Next()
	as.True(ok)
	as.Number(55.12, val)

	_, ok = i1.Next()
	as.False(ok)

	// iterate over the nested list
	i2 := a.Iterate(list2)
	val, ok = i2.Next()
	as.True(ok)
	as.String("hello", val)

	val, ok = i2.Next()
	as.True(ok)
	as.String("there", val)

	_, ok = i2.Next()
	as.False(ok)
}

func testReaderError(t *testing.T, src string, err error) {
	as := assert.New(t)

	defer as.ExpectError(err)

	l := r.Scan(s(src))
	tr := r.Read(l)
	a.Last(tr)
}

func TestReaderErrors(t *testing.T) {
	testReaderError(t, "(99 100 ", a.ErrStr(r.ListNotClosed))
	testReaderError(t, "[99 100 ", a.ErrStr(r.VectorNotClosed))
	testReaderError(t, "{:key 99", a.ErrStr(r.MapNotClosed))

	testReaderError(t, "99 100)", a.ErrStr(r.UnmatchedListEnd))
	testReaderError(t, "99 100]", a.ErrStr(r.UnmatchedVectorEnd))
	testReaderError(t, "99}", a.ErrStr(r.UnmatchedMapEnd))
	testReaderError(t, "{99}", a.ErrStr(r.MapNotPaired))

	testReaderError(t, "(", a.ErrStr(r.ListNotClosed))
	testReaderError(t, "'", a.ErrStr(r.PrefixedNotPaired, "sputter:quote"))
	testReaderError(t, "~@", a.ErrStr(r.PrefixedNotPaired, "sputter:unquote-splicing"))
	testReaderError(t, "~", a.ErrStr(r.PrefixedNotPaired, "sputter:unquote"))
}
