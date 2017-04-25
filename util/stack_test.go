package util_test

import (
	"testing"

	"github.com/kode4food/sputter/assert"
	u "github.com/kode4food/sputter/util"
)

func TestStack(t *testing.T) {
	as := assert.New(t)

	st := u.NewStack()

	v, ok := st.Peek()
	as.Nil(v)
	as.False(ok)

	st.Push("you")
	v, ok = st.Peek()
	as.String("you", v)
	as.True(ok)

	st.Push("are")
	st.Push("how")

	v, ok = st.Peek()
	as.String("how", v)
	as.True(ok)

	v, ok = st.Pop()
	as.String("how", v)
	as.True(ok)

	v, ok = st.Pop()
	as.String("are", v)
	as.True(ok)

	v, ok = st.Pop()
	as.String("you", v)
	as.True(ok)

	v, ok = st.Pop()
	as.Nil(v)
	as.False(ok)

	v, ok = st.Peek()
	as.Nil(v)
	as.False(ok)
}
