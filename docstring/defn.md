# (defn name meta? [args] form+) assigns a namespace function
Will declare a function in the current namespace, which is 'user' by default.

## An Example

  (defn fib    
    {:doc "returns an element of the fibonacci sequence"}
    [i]
    (cond
      (= i 0) 0
      (= i 1) 1
      (= i 2) 1
      :else   (+ (fib (- i 2)) (fib (- i 1)))))

This example performs recursion with no tail call optimization, and no memoization. For a more performant and stack-friendly fibonacci sequence generation example, see the documentation of `lazy-seq`.
