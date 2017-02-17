package api_test

import (
	"testing"

	s "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestVector(t *testing.T) {
	a := assert.New(t)

	v := &s.Vector{"hello", "how", "are", "you?"}
	a.Equal(4, v.Count(), "vector count is correct")
	a.Equal(4, s.Count(v), "vector general count is correct")
	a.Equal("are", v.Get(2), "get by index is correct")
}
