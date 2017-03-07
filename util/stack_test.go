package util_test

import (
	"testing"

	u "github.com/kode4food/sputter/util"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	as := assert.New(t)

	st := u.NewStack()

	v, ok := st.Peek()
	as.Equal(nil, v, "1st peek is nil")
	as.False(ok, "1st peek is false")

	st.Push("you")
	v, ok = st.Peek()
	as.Equal("you", v, "2nd peek is correct")
	as.True(ok, "2nd peek is true")

	st.Push("are")
	st.Push("how")

	v, ok = st.Peek()
	as.Equal("how", v, "3rd peek is correct")
	as.True(ok, "3rd peek is true")

	v, ok = st.Pop()
	as.Equal("how", v, "1st pop is correct")
	as.True(ok, "1st pop is true")

	v, ok = st.Pop()
	as.Equal("are", v, "2nd pop is correct")
	as.True(ok, "2nd pop is true")

	v, ok = st.Pop()
	as.Equal("you", v, "3rd pop is correct")
	as.True(ok, "3rd pop is true")

	v, ok = st.Pop()
	as.Equal(nil, v, "4th pop is nil")
	as.False(ok, "4th pop is false")

	v, ok = st.Peek()
	as.Equal(nil, v, "last peek is nil")
	as.False(ok, "last peek is false")
}
