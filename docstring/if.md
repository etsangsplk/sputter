# (if pred then else?) performs simple branching
If the evaluated predicate is truthy (not false, not nil), the 'then' form is evaluated and returned, otherwise the 'else' form, if any, will be evaluated and returned.

## An Example

  (def x '(1 2 3 4 5 6 7 8))

  (if (> (len x) 3)
    "x is big"
    "x is small")
