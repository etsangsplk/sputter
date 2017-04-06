# (cond <pred then>**) performs conditional branching
For each `<pred then>` branch, the predicate will be evaluated, and if it is truthy (not false, not nil) the 'then' form is evaluated and returned, otherwise the next branch is processed.

## An Example

  (def x 99)

  (cond
    (< x 50)  "was less than 50"
    (> x 100) "was greater than 100"
    :else     "was in between")
