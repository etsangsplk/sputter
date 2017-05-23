# (defn name meta? [name**] form+) assigns a namespace function
Will declare a function in the current namespace, which is 'user' by default.

## An Example

  (defn fib
    "returns an element of the fibonacci sequence"
    {:tail true}
    [i]
    (cond
      (= i 0) 0
      (= i 1) 1
      (= i 2) 1
      :else   (+ (fib (- i 2)) (fib (- i 1)))))
