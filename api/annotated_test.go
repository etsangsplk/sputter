package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	as := assert.New(t)

	v1 := a.Metadata{"foo": true, "bar": false}
	v2 := v1.Merge(a.Metadata{"foo": false, "hello": "there"})

	s1 := v1.String()
	as.Equal(`{"foo" true, "bar" false}`, s1)

	s2 := v2.String()
	as.Equal(`{"bar" false, "hello" "there", "foo" false}`, s2)
}

func TestAnnotated(t *testing.T) {
	as := assert.New(t)

	f := a.NewFunction(nil)
	as.Equal(f, a.AssertAnnotated(f), "function asserts as annotated")

	defer expectError(as, a.ExpectedAnnotated)
	a.AssertAnnotated(&a.Symbol{})
}
