# (fn name? meta? [name**] form+) creates an anonymous function
Will create an anonymous function that be passed around in a first-class manner.

## An Example

  (def double
    (let [mul 2]
      (fn "doubles values" [x] (** x mul))))

  (to-vector
    (map double '(1 2 3 4 5 6)))

This example will return the vector _[2 4 6 8 10 12]_.

Anonymous functions produce a closure that includes the variables being drawn from the local scope. This means that when you pass the function around, it continues to retain the variable state of the code the surrounded it.
