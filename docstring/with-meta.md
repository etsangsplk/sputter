# (with-meta form mapped+) attaches metadata to a form
If the specified form is capable of being annotated, this function will return a new instance of that value with the specified metadata attached. The metadata is provided in the form of mapped sequences, such as associative structures.

## An Example

  (with-meta
    (fn [x] (** x 2))
    {:doc  "a function that doubles things"})
