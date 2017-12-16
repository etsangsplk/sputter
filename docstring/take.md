# (take count seq+) takes the first elements of sequences
Will return a lazy sequence of either /count/ or fewer elements from the beginning of the provided sequences (concatenated). If the source sequences are shorter than the requested count, the resulting sequence will be truncated.

## An Example

  (def x '(1 2 3 4))
  (def y [5 6 7 8])
  (take 6 x y)

This example will return the lazy sequence _(1 2 3 4 5 6)_.
