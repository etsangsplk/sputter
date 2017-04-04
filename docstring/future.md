
(future form) evaluates a form asynchronously

Returns a future in the form of a function.  The provided form
will be evaluated in a separate thread of execution, and any calls
to the function will block until the form has been completely
evaluated.  For example:

    (def fut (future (do
        (to-vector (async
            (emit "red")
            (emit "orange")
            (emit "yellow")
        )))))

    (fut)
