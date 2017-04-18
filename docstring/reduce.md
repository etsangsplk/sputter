# (reduce func seq+) reduces sequences
Iterates over a set of sequences, reducing their elements to a single resulting value. The function provided must take two arguments. The first and second sequence elements encountered are the initial values applied to that function. Thereafter, the result of the previous calculation is used as the first argument, while the next element is used as the second argument.

## An Example

  (def x [1 2 3 4 5])
  (def y '(6 7 8 9 10))
  (reduce + x y)

This will return the value _55_.

Unlike other implementations of the reduce function, this one only takes its initial left value from the sequences provided. But it is easy enough to reproduce this behavior using an initial single-element vector.
