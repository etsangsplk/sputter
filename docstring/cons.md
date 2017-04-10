# (cons value seq) combines a value with a sequence
With an ordered sequence, such as a list or vector, the value is prepended to the sequence. With an unordered sequence, such as an associative array, there is no guarantee regarding position.

The name `cons` is a vestige of when Lisp implementations *cons*tructed new lists or cells by pairing a `car` (*c*ontents of the *a*ddress part of *r*egister) with a `cdr` (*c*ontents of the *d*ecrement part of *r*egister).

## An Example

  (def x '(3 4 5 6))
  (def y (cons 2 x))
  (def z (cons 1 y))
  
  (prn z)
