;;;; sputter core: basics

(defmacro def
  {:doc-asset "def"}
  [name value]
  `(ns-put (sputter:ns) ~name ~value))

(defmacro eq
  [value & comps]
  `(is-eq ~value ~@comps))

(defmacro !eq
  [value & comps]
  `(not (is-eq ~value ~@comps)))

(defmacro cons-many
  {:private true}
  ([]     ())
  ([args] args)
  ([arg & args]
    `(cons ~arg (cons-many ~@args))))

(defmacro apply
  {:doc-asset "apply"}
  [func & args]
  `(apply* ~func (cons-many ~@args)))

(defmacro let
  {:doc-asset "let"}
  [bindings & forms]
  `(let* ~bindings ~@forms))
