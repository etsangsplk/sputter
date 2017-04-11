# (meta form) returns the metadata for a form
If the form is annotated, then the metadata for that form will be returned as an associative structure.

## An Example

  (defn double [x] (** x 2))
  (meta double)

This will return the following associative structure:

_  {:type "function", :name double, :doc ""}_
