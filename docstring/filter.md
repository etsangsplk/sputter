# (filter func seq+) lazily filters sequences
Creates a lazy sequence whose content is the result of applying the provided function to the values of the provided sequences. If the result of the application is truthy (not false, not nil) then the value will be included in the resulting sequence.

## An Example

  (filter (fn [x] (< x 3)) [1 2 3 4])

This will return the lazy sequence *(1 2)*
