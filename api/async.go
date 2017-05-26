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

// Do is a callback interface for performing some action
type Do func(func())

// Emitter is an interface that is used to emit Values to a Channel
type Emitter interface {
	Value
	Emit(Value) Emitter
	Error(interface{}) Emitter
	Close() Emitter
}

// Promise represents a Value that will eventually be resolved
type Promise interface {
	Value
	Deliver(Value) Value
	Value() Value
}

type channelResult struct {
	value Value
	error interface{}
}

type channelWrapper struct {
	seq    chan channelResult
	status uint32
}

type channelEmitter struct {
	ch *channelWrapper
}

type channelSequence struct {
	once Do
	ch   *channelWrapper

	isSeq  bool
	result channelResult
	rest   Sequence
}

type promise struct {
	cond  *sync.Cond
	state uint32
	val   Value
}

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

// Never returns a Do instance for performing never performing an action
func Never() Do {
	return func(_ func()) {
		// no-op
	}
}

func (ch *channelWrapper) close() {
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

// NewChannelEmitter produces an Emitter for sending Values to a Go chan
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

// Emit will send a Value to the Go chan
func (e *channelEmitter) Emit(v Value) Emitter {
	if s := atomic.LoadUint32(&e.ch.status); s == ready {
		e.ch.seq <- channelResult{v, nil}
	}
	if s := atomic.LoadUint32(&e.ch.status); s == closeRequested {
		e.Close()
	}
	return e
}

// Error will send an Error to the Go chan
func (e *channelEmitter) Error(err interface{}) Emitter {
	if s := atomic.LoadUint32(&e.ch.status); s == ready {
		e.ch.seq <- channelResult{nil, err}
	}
	e.Close()
	return e
}

// Close will close the Go chan
func (e *channelEmitter) Close() Emitter {
	runtime.SetFinalizer(e, nil)
	e.ch.close()
	return e
}

func (e *channelEmitter) Type() Name {
	return "channel-emitter"
}

func (e *channelEmitter) Eval(_ Context) Value {
	return e
}

func (e *channelEmitter) Str() Str {
	return MakeDumpStr(e)
}

// NewChannelSequence produces a new Sequence whose Values come from a Go chan
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
		c.ch = nil
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

func (c *channelSequence) Eval(_ Context) Value {
	return c
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

func (p *promise) Eval(_ Context) Value {
	return p
}

func (p *promise) Str() Str {
	return MakeDumpStr(p)
}
