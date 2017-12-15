# (-> value forms**) threads value through calls as their first argument
Threads 'value' through the supplied forms. Any form that is not already an application will be converted into one. The value will then be inserted as its first argument and applied, and so on.

## An Example

  (-> 0 (+ 10) (* 2) (/ 5))

Will expand to `(/ (* (+ 0 10) 2) 5)` and return _4_. In order to better visualize what's going on, one might choose to insert a `,` as a placeholder for the threaded value.

  (-> 0 (+ , 10) (* , 2) (/ , 5))
