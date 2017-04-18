# (drop count seq+) drops the first elements of sequences
Will return a lazy sequence that excludes the first 'count' elements of the provided sequences (concatenated). If the source sequences are shorter than the requested count, an empty list will be returned.

## An Example

  (def x '(1 2 3 4))
  (def y [5 6 7 8])
  (drop 3 x y)

This example will return the lazy sequence _(4 5 6)_.
