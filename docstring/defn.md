# (defn name meta? [name**] form+) assigns a namespace function
Will assign a function as a variable in the current namespace, which is 'user' by default. All variables are immutable and result in an error being raised if an attempt is made to re-assign them. This behavior is different than most Lisps, as they will generally fail silently in such cases.

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
