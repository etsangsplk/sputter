package api_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

type ncSeq struct{}

func (n *ncSeq) First() a.Value                     { return nil }
func (n *ncSeq) Rest() a.Sequence                   { return nil }
func (n *ncSeq) Split() (a.Value, a.Sequence, bool) { return nil, nil, false }
func (n *ncSeq) Prepend(v a.Value) a.Sequence       { return nil }
func (n *ncSeq) IsSequence() bool                   { return false }
func (n *ncSeq) Str() a.Str                         { return s("()") }

func TestNonCountableSequence(t *testing.T) {
	as := assert.New(t)
	nc := new(ncSeq)

	e := cvtErr("*api_test.ncSeq", "api.CountedSequence", "Count")
	defer as.ExpectError(e)
	a.Count(nc)
}
