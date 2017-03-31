package api_test

import (
	"testing"
	"time"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	ch := make(chan a.Value)

	go func() {
		ch <- 2
		time.Sleep(time.Millisecond * 50)
		ch <- 3
		time.Sleep(time.Millisecond * 30)
		ch <- "foo"
		time.Sleep(time.Millisecond * 10)
		ch <- "bar"
		close(ch)
	}()

	check := func() {
		as := assert.New(t)
		s := a.NewChannel(ch).Prepend(1)
		as.Equal(1, s.First(), "first is right")
		as.Equal(2, s.Rest().First(), "second is right")
		as.Equal(3, s.Rest().Rest().First(), "third is right")
		as.True(s.Rest().Rest().Rest().IsSequence(), "more!")
		as.Equal("foo", s.Rest().Rest().Rest().First(), "foo is right")
		as.True(s.Rest().Rest().Rest().Rest().IsSequence(), "more!")
		as.Equal("bar", s.Rest().Rest().Rest().Rest().First(), "bar is right")
		as.False(s.Rest().Rest().Rest().Rest().Rest().IsSequence(), "eof")
	}

	go check()
	go check()
}
