package api

import "sync"

// NewChannel produces a new Sequence whose Values come from a Go chan
func NewChannel(ch chan Value) Sequence {
	return &channel{
		ch:   ch,
		cond: &sync.Cond{L: &sync.Mutex{}},
		rest: EmptyList,
	}
}

func (c *channel) resolve() *channel {
	c.cond.L.Lock()
	if c.ch == nil {
		c.cond.Wait()
		c.cond.L.Unlock()
		return c
	}

	ch := c.ch
	c.ch = nil
	c.cond.L.Unlock()

	if first, isSeq := <-ch; isSeq {
		c.isSeq = isSeq
		c.first = first
		c.rest = NewChannel(ch)
	}

	c.cond.Broadcast()
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
		cond:  &sync.Cond{L: &sync.Mutex{}},
		first: v,
		rest:  c,
	}
}
