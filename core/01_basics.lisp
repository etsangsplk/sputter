;;;; sputter core: basics

(defmacro def
  {:doc-asset "def"}
  [name value]
  `(ns-put (sputter:ns) ~name ~value))
