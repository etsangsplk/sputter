;;;; sputter core: basics

(defmacro def
  {:doc-asset "def"}
  [name value]
  `(ns-put (ns) ~name ~value))

(defmacro eq
  [value & comps]
  `(is-eq ~value ~@comps))

(defmacro !eq
  [value & comps]
  `(not (is-eq ~value ~@comps)))
