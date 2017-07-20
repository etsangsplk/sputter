package api_test

import (
	"sync"
	"testing"
	"time"

	a "github.com/kode4food/sputter/api"
	"github.com/kode4food/sputter/assert"
)

func TestConditionals(t *testing.T) {
	as := assert.New(t)

	i := 0
	inc := func() {
		i++
	}

	once := a.Once()
	never := a.Never()
	always := a.Always()

	as.Number(0, i)
	once(inc)
	as.Number(1, i)
	once(inc)
	as.Number(1, i)

	never(inc)
	as.Number(1, i)
	never(inc)
	as.Number(1, i)

	always(inc)
	as.Number(2, i)
	always(inc)
	as.Number(3, i)
	always(inc)
	as.Number(4, i)
}

func TestChannel(t *testing.T) {
	as := assert.New(t)

	e, seq := a.NewChannel()
	seq = seq.Prepend(f(1))
	as.Contains(":type channel-emitter", e)
	as.Contains(":type channel-sequence", seq)

	var wg sync.WaitGroup

	gen := func() {
		e.Write(f(2))
		time.Sleep(time.Millisecond * 50)
		e.Write(f(3))
		time.Sleep(time.Millisecond * 30)
		e.Write(s("foo"))
		time.Sleep(time.Millisecond * 10)
		e.Write(s("bar"))
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
	as.String("hello", p1.Resolve())
	p1.Deliver(s("hello"))
	as.String("hello", p1.Resolve())

	defer as.ExpectError(a.ErrStr(a.ExpectedUndelivered))
	p1.Deliver(s("goodbye"))
}
