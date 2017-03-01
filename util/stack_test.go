package util_test

import (
	"testing"

	u "github.com/kode4food/sputter/util"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	a := assert.New(t)

	st := u.NewStack()

	v, ok := st.Peek()
	a.Equal(nil, v, "1st peek is nil")
	a.False(ok, "1st peek is false")

	st.Push("you")
	v, ok = st.Peek()
	a.Equal("you", v, "2nd peek is correct")
	a.True(ok, "2nd peek is true")

	st.Push("are")
	st.Push("how")

	v, ok = st.Peek()
	a.Equal("how", v, "3rd peek is correct")
	a.True(ok, "3rd peek is true")

	v, ok = st.Pop()
	a.Equal("how", v, "1st pop is correct")
	a.True(ok, "1st pop is true")

	v, ok = st.Pop()
	a.Equal("are", v, "2nd pop is correct")
	a.True(ok, "2nd pop is true")

	v, ok = st.Pop()
	a.Equal("you", v, "3rd pop is correct")
	a.True(ok, "3rd pop is true")

	v, ok = st.Pop()
	a.Equal(nil, v, "4th pop is nil")
	a.False(ok, "4th pop is false")

	v, ok = st.Peek()
	a.Equal(nil, v, "last peek is nil")
	a.False(ok, "last peek is false")
}
