package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestMetadata(t *testing.T) {
	as := assert.New(t)

	v1 := a.Properties{
		s("foo"): a.True,
		s("bar"): a.False,
	}

	v2 := v1.Child(a.Properties{
		s("foo"):   a.False,
		s("hello"): s("there"),
	})

	as.True(a.IsTrue(v1, s("foo")))
	as.False(a.IsTrue(v2, s("foo")))
	as.False(a.IsTrue(v2, s("missing")))

	s1 := v1.Str()
	as.Contains(`"foo" true`, s1)
	as.Contains(`"bar" false`, s1)

	s2 := v2.Str()
	as.Contains(`"bar" false`, s2)
	as.Contains(`"hello" "there"`, s2)
	as.Contains(`"foo" false`, s2)
	as.NotContains(`"foo" true`, s2)

	v1 = a.Properties{}
	v2 = a.Properties{s("test"): a.True}
	v3 := v1.Child(v2.Flatten())
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

	defer as.ExpectError(a.ErrStr(a.ExpectedAnnotated, f(99)))
	a.AssertAnnotated(f(99))
}

func TestGetDocumentation(t *testing.T) {
	as := assert.New(t)

	f1 := a.NewFunction(nil).WithMetadata(a.Properties{
		a.DocKey: a.Str("hello there"),
	})

	f2 := a.NewFunction(nil).WithMetadata(a.Properties{
		a.DocAssetKey: a.Str("if"),
	})

	as.String("hello there", a.GetDocumentation(f1))
	as.Contains("## An Example", a.GetDocumentation(f2))
}
