package api

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// ExpectedUndelivered is thrown on an attempt to deliver a Promise twice
const ExpectedUndelivered = "can't deliver a promise twice"

const (
	undeliveredState uint32 = iota
	deliveredState
)

// Do is a callback interface for performing some action
type Do func(func())

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
	once Do
	ch   chan Value

	isSeq bool
	first Value
	rest  Sequence
}

type promise struct {
	cond  *sync.Cond
	state uint32
	val   Value
}

// Once creates a Do instance for performing an action only once
func Once() Do {
	var state = undeliveredState
	var mutex sync.Mutex

	return func(f func()) {
		if atomic.LoadUint32(&state) == deliveredState {
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		if state == undeliveredState {
			defer atomic.StoreUint32(&state, deliveredState)
			f()
		}
	}
}

// Never returns a Do instance for performing never performing an action
func Never() Do {
	return func(_ func()) {
		// no-op
	}
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
		once: Once(),
		ch:   ch,
		rest: EmptyList,
	}
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		ch := c.ch
		if first, isSeq := <-ch; isSeq {
			c.isSeq = isSeq
			c.first = first
			c.rest = NewChannelSequence(ch)
		}
		c.ch = nil
	})
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
		once:  Never(),
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
		cond:  sync.NewCond(new(sync.Mutex)),
		state: undeliveredState,
	}
}

func (p *promise) Value() Value {
	if atomic.LoadUint32(&p.state) == deliveredState {
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
	if atomic.LoadUint32(&p.state) == deliveredState {
		return p.checkNewValue(v)
	}

	cond := p.cond
	cond.L.Lock()
	defer cond.L.Unlock()

	if p.state == undeliveredState {
		p.val = v
		atomic.StoreUint32(&p.state, deliveredState)
		cond.Broadcast()
		return v
	}

	cond.Wait()
	return p.checkNewValue(v)
}

func (p *promise) Type() Name {
	return "promise"
}

// Str converts this Value into a Str
func (p *promise) Str() Str {
	return MakeDumpStr(p)
}
