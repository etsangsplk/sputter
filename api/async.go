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

const (
	ready uint32 = iota
	closeRequested
	closed
)

type (
	// Do is a callback interface for performing some action
	Do func(func())

	// Emitter is an interface that is used to emit values to a Channel
	Emitter interface {
		Writer
		Closer
		Error(interface{})
	}

	// Promise represents a Value that will eventually be resolved
	Promise interface {
		Value
		Deliver(Value) Value
		Resolve() Value
	}

	channelResult struct {
		value Value
		error interface{}
	}

	channelWrapper struct {
		seq    chan channelResult
		status uint32
	}

	channelEmitter struct {
		ch *channelWrapper
	}

	channelSequence struct {
		once Do
		ch   *channelWrapper

		isSeq  bool
		result channelResult
		rest   Sequence
	}

	promise struct {
		cond  *sync.Cond
		state uint32
		val   Value
	}
)

var emptyResult = channelResult{value: Nil, error: nil}

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

// Always returns a Do instance for always performing an action
func Always() Do {
	return func(f func()) {
		f()
	}
}

// Never returns a Do instance for never performing an action
func Never() Do {
	return func(_ func()) {
		// no-op
	}
}

func (ch *channelWrapper) Close() {
	if status := atomic.LoadUint32(&ch.status); status != closed {
		atomic.StoreUint32(&ch.status, closed)
		close(ch.seq)
	}
}

// NewChannel produces a Emitter and Sequence pair
func NewChannel() (Emitter, Sequence) {
	seq := make(chan channelResult, 0)
	ch := &channelWrapper{
		seq:    seq,
		status: ready,
	}
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending values to a Go chan
func NewChannelEmitter(ch *channelWrapper) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		defer func() { recover() }()
		if s := atomic.LoadUint32(&ch.status); s != closed {
			e.Close()
		}
	})
	return r
}

// Write will send a Value to the Go chan
func (e *channelEmitter) Write(v Value) {
	if s := atomic.LoadUint32(&e.ch.status); s == ready {
		e.ch.seq <- channelResult{v, nil}
	}
	if s := atomic.LoadUint32(&e.ch.status); s == closeRequested {
		e.Close()
	}
}

// Error will send an Error to the Go chan
func (e *channelEmitter) Error(err interface{}) {
	if s := atomic.LoadUint32(&e.ch.status); s == ready {
		e.ch.seq <- channelResult{nil, err}
	}
	e.Close()
}

// Close will Close the Go chan
func (e *channelEmitter) Close() {
	runtime.SetFinalizer(e, nil)
	e.ch.Close()
}

func (e *channelEmitter) Type() Name {
	return "channel-emitter"
}

func (e *channelEmitter) Str() Str {
	return MakeDumpStr(e)
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch *channelWrapper) Sequence {
	r := &channelSequence{
		once:   Once(),
		ch:     ch,
		result: emptyResult,
		rest:   EmptyList,
	}
	runtime.SetFinalizer(r, func(c *channelSequence) {
		defer func() { recover() }()
		if s := atomic.LoadUint32(&c.ch.status); s == ready {
			atomic.StoreUint32(&c.ch.status, closeRequested)
			<-ch.seq // consume whatever is there
		}
	})
	return r
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		runtime.SetFinalizer(c, nil)
		ch := c.ch
		if result, isSeq := <-ch.seq; isSeq {
			c.isSeq = isSeq
			c.result = result
			c.rest = NewChannelSequence(ch)
		}
	})
	if e := c.result.error; e != nil {
		panic(e)
	}
	return c
}

func (c *channelSequence) IsSequence() bool {
	return c.resolve().isSeq
}

func (c *channelSequence) First() Value {
	return c.resolve().result.value
}

func (c *channelSequence) Rest() Sequence {
	return c.resolve().rest
}

func (c *channelSequence) Split() (Value, Sequence, bool) {
	r := c.resolve()
	return r.result.value, r.rest, r.isSeq
}

func (c *channelSequence) Prepend(v Value) Sequence {
	return &channelSequence{
		once:   Never(),
		isSeq:  true,
		result: channelResult{value: v, error: nil},
		rest:   c,
	}
}

func (c *channelSequence) Type() Name {
	return "channel-sequence"
}

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

func (p *promise) Apply(_ Context, args Vector) Value {
	if AssertArityRange(args, 0, 1) == 1 {
		return p.Deliver(args[0])
	}
	return p.Resolve()
}

func (p *promise) Resolve() Value {
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
	panic(ErrStr(ExpectedUndelivered))
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

func (p *promise) Str() Str {
	return MakeDumpStr(p)
}
