package api

import "sync"

type channel struct {
	ch    chan Value
	cond  *sync.Cond
	ready bool

	isSeq bool
	first Value
	rest  Sequence
}

// NewChannel produces a new Sequence whose Values come from a Go chan
func NewChannel(ch chan Value) Sequence {
	return &channel{
		ch:   ch,
		cond: &sync.Cond{L: &sync.Mutex{}},
		rest: EmptyList,
	}
}

func (c *channel) resolve() *channel {
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
		c.rest = NewChannel(ch)
	}

	c.ready = true
	c.cond = nil
	cond.Broadcast()

	return c
}

func (c *channel) IsSequence() bool {
	return c.resolve().isSeq
}

func (c *channel) First() Value {
	return c.resolve().first
}

func (c *channel) Rest() Sequence {
	return c.resolve().rest
}

func (c *channel) Prepend(v Value) Sequence {
	return &channel{
		ready: true,
		isSeq: true,
		first: v,
		rest:  c,
	}
}
