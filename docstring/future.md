# (future form) evaluates a form asynchronously
Returns a future in the form of a function. The provided form will be evaluated in a separate thread of execution, and any calls to the function will *block* until the form has been completely evaluated.

## An Example

  (def fut (future
    (to-vector (async
      (emit "red")
      (emit "orange")
      (emit "yellow")))))

  (fut)

This example produces a future called `fut` that converts the results of an asynchronous block into a vector. The `(fut)` call will block until the future returns a value.
