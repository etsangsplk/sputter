package api_test

import (
	"sync"
	"testing"
	"time"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestChannel(t *testing.T) {
	as := assert.New(t)

	e, seq := a.NewChannel(1)
	seq = seq.Prepend(f(1))
	as.Contains(":type channel-emitter", e)
	as.Contains(":type channel-sequence", seq)

	var wg sync.WaitGroup

	gen := func() {
		e.Emit(f(2))
		time.Sleep(time.Millisecond * 50)
		e.Emit(f(3))
		time.Sleep(time.Millisecond * 30)
		e.Emit(s("foo"))
		time.Sleep(time.Millisecond * 10)
		e.Emit(s("bar"))
		e.Close()
		wg.Done()
	}

	check := func() {
		as.Number(1, seq.First())
		as.Number(2, seq.Rest().First())
		as.Number(3, seq.Rest().Rest().First())
		as.True(seq.Rest().Rest().Rest().IsSequence())
		as.String("foo", seq.Rest().Rest().Rest().First())
		as.True(seq.Rest().Rest().Rest().Rest().IsSequence())
		as.String("bar", seq.Rest().Rest().Rest().Rest().First())
		as.False(seq.Rest().Rest().Rest().Rest().Rest().IsSequence())
		wg.Done()
	}

	wg.Add(4)
	go check()
	go check()
	go gen()
	go check()
	wg.Wait()
}

func TestPromise(t *testing.T) {
	as := assert.New(t)
	p1 := a.NewPromise()

	go func() {
		time.Sleep(time.Millisecond * 50)
		p1.Deliver(s("hello"))
	}()

	as.Contains(":type promise", p1)
	as.String("hello", p1.Value())
	p1.Deliver(s("hello"))
	as.String("hello", p1.Value())

	defer expectError(as, a.ExpectedUndelivered)
	p1.Deliver(s("goodbye"))
}
