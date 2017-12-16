# (meta? form+) tests whether the provided forms have metadata
If all forms evaluate to annotated values, then this function will return _true_. The first non-annotated will result in the function returning _false_.

  (defn double [x] (** x 2))
  (meta? double)

This example will return _true_ because functions can be annotated.

Like most predicates, this function can also be negated by prepending the /!/ character. This means that all of the provided forms must not allow annotation.

  (!meta? "hello" 1)

This example will return _true_ because strings and numbers can't be annotated.
