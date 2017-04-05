# (async form+) generates an asynchronous sequence
Evaluates the specified forms in a separate thread of execution. Returns a sequence that will iterate over any of the values that are emitted using a locally scoped function of the form `(emit value)`.

## An Example

  (def aseq (async
    (emit "red")
    (emit "orange")
    (emit "yellow")))

  (to-vector aseq)
