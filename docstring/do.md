# (do form**) evaluates multiple forms
Will evaluate each form in turn, returning the final evaluation as its result.

## An Example

  (when (> x 9)
    (do (prn x)
        (** x 10)))
