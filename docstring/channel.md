
(channel) returns a unidirectional channel

A channel is a data structure that is used to generate a lazy
sequence of values.  It consists of an `emit` function, a
`close` function, and a sequence.  Retrieving an element from
the sequence will block, waiting for the next value to be
emitted or for the channel to be closed.  A call to channel
returns a hash-map with the following elements:

    {
        :seq   <the sequence to be generated>
        :emit  an emit function of the form `(emit value)`
        :close a function to close the channel `(close)`
    }

For example:

    (let [c (channel) emit (:emit c) close (:close c) seq (:seq c)]
        (emit "foo")
        (emit "bar")
        (close)

        (to-vector seq))
