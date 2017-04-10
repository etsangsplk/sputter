# (def name form) assigns a namespace variable
Will assign a variable in the current namespace, which is 'user' by default. All variables are immutable and result in an error being raised if an attempt is made to re-assign them. This behavior is different than most Lisps, as they will generally fail silently in such cases.

## An Example

  (def x
    (map
      (fn [y] y ** 2)
      seq1 seq2 seq3))

This example will create a lazy map where each value of the three provided sequences is doubled upon request. It will then assign that lazy map to the namespace variable 'x'.
