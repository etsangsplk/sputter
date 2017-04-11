# (map func seq+) lazily maps sequences
Creates a lazy sequence whose values are the result of applying the provided function to the values of the provided sequences.

## An Example

  (map (fn [x] (* x 2)) [1 2 3 4])

This will return the lazy sequence *(2 4 6 8)*
