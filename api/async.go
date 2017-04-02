package api

import (
	"runtime"
	"sync"
)

// ExpectedUndelivered is thrown on an attempt to deliver a Promise twice
const ExpectedUndelivered = "Can't deliver a Promise twice"

// Emitter is an interface that is used to emit Values to a Channel
type Emitter interface {
	Emit(Value) Emitter
	Close() Emitter
}

// Promise represents a Value that will eventually be resolved
type Promise interface {
	Deliver(Value)
	Value() Value
}

type emitter struct {
	ch chan Value
}

type channelSequence struct {
	ch    chan Value
	cond  *sync.Cond
	ready bool

	isSeq bool
	first Value
	rest  Sequence
}

type promise struct {
	cond  *sync.Cond
	ready bool
	val   Value
}

// NewChannel produces a Emitter and Sequence pair
func NewChannel(buf int) (Emitter, Sequence) {
	ch := make(chan Value, buf)
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending Values to a Go chan
func NewChannelEmitter(ch chan Value) Emitter {
	r := &emitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *emitter) {
		if e.ch != nil {
			close(e.ch)
			e.ch = nil
		}
	})
	return r
}

// Emit will send a Value to the Go chan
func (e *emitter) Emit(v Value) Emitter {
	if e.ch != nil {
		e.ch <- v
	}
	return e
}

// Close will close the Go chan
func (e *emitter) Close() Emitter {
	if e.ch != nil {
		close(e.ch)
		e.ch = nil
		runtime.SetFinalizer(e, nil)
	}
	return e
}

// NewChannelSequence produces a new Sequence whose Values come from a Go chan
func NewChannelSequence(ch chan Value) Sequence {
	return &channelSequence{
		ch:   ch,
		cond: &sync.Cond{L: &sync.Mutex{}},
		rest: EmptyList,
	}
}

func (c *channelSequence) resolve() *channelSequence {
	if c.ready {
		return c
	}

	cond := c.cond
	cond.L.Lock()

	if c.ch == nil {
		cond.Wait()
		cond.L.Unlock()
		return c
	}

	ch := c.ch
	c.ch = nil
	cond.L.Unlock()

	if first, isSeq := <-ch; isSeq {
		c.isSeq = isSeq
		c.first = first
		c.rest = NewChannelSequence(ch)
	}

	c.ready = true
	c.cond = nil
	cond.Broadcast()

	return c
}

func (c *channelSequence) IsSequence() bool {
	return c.resolve().isSeq
}

func (c *channelSequence) First() Value {
	return c.resolve().first
}

func (c *channelSequence) Rest() Sequence {
	return c.resolve().rest
}

func (c *channelSequence) Prepend(v Value) Sequence {
	return &channelSequence{
		ready: true,
		isSeq: true,
		first: v,
		rest:  c,
	}
}

// NewPromise instantiates a new Promise
func NewPromise() Promise {
	return &promise{
		cond: &sync.Cond{L: &sync.Mutex{}},
	}
}

func (p *promise) Value() Value {
	if p.ready {
		return p.val
	}

	cond := p.cond
	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()
	return p.val
}

func (p *promise) Deliver(v Value) {
	cond := p.cond
	cond.L.Lock()

	if p.ready {
		cond.L.Unlock()
		if v == p.val {
			return
		}
		panic(ExpectedUndelivered)
	}

	p.val = v
	p.ready = true
	cond.L.Unlock()
	cond.Broadcast()
}
