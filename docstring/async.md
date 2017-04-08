# (async form+) generates an asynchronous sequence
Evaluates the specified forms in a separate thread of execution. Returns a sequence that will iterate over any of the values that are emitted using a locally scoped function of the form `(emit value)`.

## An Example

  (def colors (async
    (emit "red")
    (emit "orange")
    (emit "yellow")))

  (to-vector colors)
