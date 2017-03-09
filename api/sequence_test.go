package api_test

import (
	"math/big"
	"testing"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestAssertSequence(t *testing.T) {
	as := assert.New(t)
	a.AssertSequence(a.NewList("hello"))

	defer expectError(as, a.ExpectedSequence)
	a.AssertSequence(big.NewFloat(99))
}
