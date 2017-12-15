# (map func seq+) lazily maps sequences
Creates a lazy sequence whose elements are the result of applying the provided function to the elements of the provided sequences.

## An Example

  (map (fn [x] (** x 2)) [1 2 3 4])

This will return the lazy sequence _(2 4 6 8)_.
