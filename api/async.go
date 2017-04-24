package api

import (
	"runtime"
	"sync"
)

// ExpectedUndelivered is thrown on an attempt to deliver a Promise twice
const ExpectedUndelivered = "Can't deliver a promise twice"

type asyncState int

const (
	undeliveredState asyncState = iota
	deliveringState
	deliveredState
)

// Emitter is an interface that is used to emit Values to a Channel
type Emitter interface {
	Value
	Emit(Value) Emitter
	Close() Emitter
}

// Promise represents a Value that will eventually be resolved
type Promise interface {
	Value
	Deliver(Value) Value
	Value() Value
}

type channelEmitter struct {
	ch chan Value
}

type channelSequence struct {
	cond  *sync.Cond
	state asyncState
	ch    chan Value

	isSeq bool
	first Value
	rest  Sequence
}

type promise struct {
	cond  *sync.Cond
	state asyncState
	val   Value
}

// NewChannel produces a Emitter and Sequence pair
func NewChannel(buf int) (Emitter, Sequence) {
	ch := make(chan Value, buf)
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending Values to a Go chan
func NewChannelEmitter(ch chan Value) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		if e.ch != nil {
			close(e.ch)
			e.ch = nil
		}
	})
	return r
}

// Emit will send a Value to the Go chan
func (e *channelEmitter) Emit(v Value) Emitter {
	if e.ch != nil {
		e.ch <- v
	}
	return e
}

// Close will close the Go chan
func (e *channelEmitter) Close() Emitter {
	if e.ch != nil {
		close(e.ch)
		e.ch = nil
		runtime.SetFinalizer(e, nil)
	}
	return e
}

func (e *channelEmitter) Type() Name {
	return "channel-emitter"
}

// Str converts this Value into a Str
func (e *channelEmitter) Str() Str {
	return MakeDumpStr(e)
}

// NewChannelSequence produces a new Sequence whose Values come from a Go chan
func NewChannelSequence(ch chan Value) Sequence {
	return &channelSequence{
		cond:  &sync.Cond{L: &sync.Mutex{}},
		state: undeliveredState,
		ch:    ch,
		rest:  EmptyList,
	}
}

func (c *channelSequence) resolve() *channelSequence {
	cond := c.cond
	if c.state == deliveredState {
		return c
	}

	cond.L.Lock()
	if c.state == deliveredState {
		cond.L.Unlock()
		return c
	}

	if c.state == deliveringState {
		cond.Wait()
		cond.L.Unlock()
		return c
	}

	c.state = deliveringState
	ch := c.ch
	c.ch = nil
	cond.L.Unlock()

	if first, isSeq := <-ch; isSeq {
		c.isSeq = isSeq
		c.first = first
		c.rest = NewChannelSequence(ch)
	}

	c.state = deliveredState
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
		state: deliveredState,
		isSeq: true,
		first: v,
		rest:  c,
	}
}

func (c *channelSequence) Type() Name {
	return "channel-sequence"
}

// Str converts this Value into a Str
func (c *channelSequence) Str() Str {
	return MakeDumpStr(c)
}

// NewPromise instantiates a new Promise
func NewPromise() Promise {
	return &promise{
		cond:  &sync.Cond{L: &sync.Mutex{}},
		state: undeliveredState,
	}
}

func (p *promise) Value() Value {
	if p.state == deliveredState {
		return p.val
	}

	cond := p.cond
	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()
	return p.val
}

func (p *promise) checkNewValue(v Value) Value {
	if v == p.val {
		return p.val
	}
	panic(ExpectedUndelivered)
}

func (p *promise) Deliver(v Value) Value {
	cond := p.cond
	if p.state == deliveredState {
		return p.checkNewValue(v)
	}

	cond.L.Lock()
	if p.state == deliveredState {
		cond.L.Unlock()
		return p.checkNewValue(v)
	}

	if p.state == deliveringState {
		cond.Wait()
		cond.L.Unlock()
		return p.checkNewValue(v)
	}

	p.state = deliveringState
	cond.L.Unlock()

	p.val = v
	p.state = deliveredState
	p.cond = nil
	cond.Broadcast()
	return v
}

func (p *promise) Type() Name {
	return "promise"
}

// Str converts this Value into a Str
func (p *promise) Str() Str {
	return MakeDumpStr(p)
}
