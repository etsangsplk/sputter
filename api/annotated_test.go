package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestMetadata(t *testing.T) {
	as := assert.New(t)

	v1 := a.Metadata{
		s("foo"): a.True,
		s("bar"): a.False,
	}
	v2 := v1.Merge(a.Metadata{
		s("foo"):   a.False,
		s("hello"): s("there"),
	})

	s1 := v1.Str()
	as.Contains(`"foo" true`, s1)
	as.Contains(`"bar" false`, s1)

	s2 := v2.Str()
	as.Contains(`"bar" false`, s2)
	as.Contains(`"hello" "there"`, s2)
	as.Contains(`"foo" false`, s2)
	as.NotContains(`"foo" true`, s2)

	v1 = a.Metadata{}
	v2 = a.Metadata{s("test"): a.True}
	v3 := v1.Merge(v2)
	s3 := v3.Str()
	as.Contains(`"test" true`, s3)
	as.Equal(v2, v3)

	r1, ok := v2.Get(s("test"))
	as.True(ok)
	as.True(r1)

	r2, ok := v2.Get(s("missing"))
	as.False(ok)
	as.Nil(r2)
}

func TestAnnotated(t *testing.T) {
	as := assert.New(t)

	fn := a.NewFunction(nil)
	as.Identical(fn, a.AssertAnnotated(fn))

	defer expectError(as, a.Err(a.ExpectedAnnotated, f(99)))
	a.AssertAnnotated(f(99))
}
