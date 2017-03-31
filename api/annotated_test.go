package api_test

import (
	"fmt"
	"strings"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	as := assert.New(t)

	v1 := a.Metadata{"foo": true, "bar": false}
	v2 := v1.Merge(a.Metadata{"foo": false, "hello": "there"})

	s1 := v1.String()
	as.True(strings.Contains(s1, `"foo" true`))
	as.True(strings.Contains(s1, `"bar" false`))

	s2 := v2.String()
	as.True(strings.Contains(s2, `"bar" false`))
	as.True(strings.Contains(s2, `"hello" "there"`))
	as.True(strings.Contains(s2, `"foo" false`))
	as.False(strings.Contains(s2, `"foo" true`))
}

func TestAnnotated(t *testing.T) {
	as := assert.New(t)

	f := a.NewFunction(nil)
	as.Equal(f, a.AssertAnnotated(f), "function asserts as annotated")

	defer expectError(as, fmt.Sprintf(a.ExpectedAnnotated, "99"))
	a.AssertAnnotated(a.NewFloat(99))
}
