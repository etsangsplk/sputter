# (channel) creates a unidirectional channel
A channel is a data structure that is used to generate a lazy sequence of values. The result is a hash-map consisting of an `emit` function, a `close` function, and a sequence. Retrieving an element from the sequence may *block*, waiting for the next value to be emitted or for the channel to be closed. Emitting a value to a channel will also block until the buffer is flushed as a result of iterating over the sequence.

## Channel Keys

*:seq*     the sequence to be generated
*:emit*    an emitter function of the form `(emit value)`
*:close*   a function to close the channel `(close)`

## An Example

  (let [c     (channel)
        c-emit  (:emit c)
        c-close (:close c)]
    (async
      (c-emit "foo")
      (c-emit "bar")
      (c-close))

    (to-vector (:seq c)))
