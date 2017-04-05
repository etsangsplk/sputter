# (concat seq+) lazily concatenates sequences
Creates a lazy sequence whose content is the result of concatenating the values of each provided sequence.

## An Example

  (concat [1 2 3] '(4 5 6))

This will return the sequence `(1 2 3 4 5 6)`
