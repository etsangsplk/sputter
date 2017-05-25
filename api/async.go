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

type channelEmitter struct {
	ch chan channelResult
}

type channelSequence struct {
	once Do
	ch   chan channelResult

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
var closeChannel = channelResult{value: Nil, error: "close"}

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
func NewChannel() (Emitter, Sequence) {
	ch := make(chan channelResult, 0)
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending Values to a Go chan
func NewChannelEmitter(ch chan channelResult) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		if ch := e.ch; ch != nil {
			defer func() { recover() }()
			close(ch)
			e.ch = nil
		}
	})
	return r
}

// Emit will send a Value to the Go chan
func (e *channelEmitter) Emit(v Value) Emitter {
	if e.ch != nil {
		e.ch <- channelResult{v, nil}
		r, _ := <-e.ch
		if r.error != nil {
			e.Close()
		}
	}
	return e
}

// Error will send an Error to the Go chan
func (e *channelEmitter) Error(err interface{}) Emitter {
	if ch := e.ch; ch != nil {
		ch <- channelResult{nil, err}
		<-ch // consume the response
		e.Close()
	}
	return e
}

// Close will close the Go chan
func (e *channelEmitter) Close() Emitter {
	runtime.SetFinalizer(e, nil)
	if ch := e.ch; ch != nil {
		close(ch)
		e.ch = nil
	}
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
func NewChannelSequence(ch chan channelResult) Sequence {
	r := &channelSequence{
		once:   Once(),
		ch:     ch,
		result: emptyResult,
		rest:   EmptyList,
	}
	runtime.SetFinalizer(r, func(c *channelSequence) {
		if ch := c.ch; ch != nil {
			defer func() { recover() }()
			<-ch // consume whatever is there
			ch <- closeChannel
			c.ch = nil
		}
	})
	return r
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		runtime.SetFinalizer(c, nil)
		ch := c.ch
		if result, isSeq := <-ch; isSeq {
			ch <- emptyResult
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
