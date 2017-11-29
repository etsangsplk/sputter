# (ns-put ns name form) binds a namespace entry
Will bind a value to a name in the specified namespace. All bindings are immutable and result in an error being raised if an attempt is made to re-bind them.

## An Example

  (ns-put (ns maps) x
    (map
      (fn [y] y ** 2)
      seq1 seq2 seq3))

This example will create a lazy map where each value of the three provided sequences is doubled upon request. It will then bind that lazy map to the namespace entry `x` in a namespace called `maps`.
