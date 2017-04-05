# (channel size?) creates a unidirectional channel
A channel is a data structure that is used to generate a lazy sequence of values. If a size is specified, the channel will be buffered. The result is a hash-map consisting of an `emit` function, a `close` function, and a sequence. Retrieving an element from the sequence will *block*, waiting for the next value to be emitted or for the channel to be closed. Emitting a value to a channel will also block if the buffer hasn't been flushed as a result of iterating over the sequence.

## Channel Keys

*:seq*     the sequence to be generated
*:emit*    an emitter function of the form `(emit value)`
*:close*   a function to close the channel `(close)`

## An Example

  (let [c (channel 2) emit (:emit c) close (:close c)]
    (emit "foo")
    (emit "bar")
    (close)

    (to-vector (:seq c)))
