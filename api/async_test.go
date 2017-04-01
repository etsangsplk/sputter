package api_test

import (
	"sync"
	"testing"
	"time"

	a "github.com/kode4food/sputter/api"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	as := assert.New(t)

	ch := make(chan a.Value)
	s := a.NewChannel(ch).Prepend(1)

	var wg sync.WaitGroup

	gen := func() {
		ch <- 2
		time.Sleep(time.Millisecond * 50)
		ch <- 3
		time.Sleep(time.Millisecond * 30)
		ch <- "foo"
		time.Sleep(time.Millisecond * 10)
		ch <- "bar"
		close(ch)
		wg.Done()
	}

	check := func() {
		as.Equal(1, s.First(), "first is right")
		as.Equal(2, s.Rest().First(), "second is right")
		as.Equal(3, s.Rest().Rest().First(), "third is right")
		as.True(s.Rest().Rest().Rest().IsSequence(), "more!")
		as.Equal("foo", s.Rest().Rest().Rest().First(), "foo is right")
		as.True(s.Rest().Rest().Rest().Rest().IsSequence(), "more!")
		as.Equal("bar", s.Rest().Rest().Rest().Rest().First(), "bar is right")
		as.False(s.Rest().Rest().Rest().Rest().Rest().IsSequence(), "eof")
		wg.Done()
	}

	wg.Add(3)
	go check()
	go check()
	go gen()
	wg.Wait()
}
