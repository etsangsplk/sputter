; sputter core: sequences

(defn concat
  {:doc-asset "concat"}
  [& forms]
  (sputter:generate
    (for-each [form forms, elem form] (emit elem))))

(defn filter
  {:doc-asset "filter"}
  [func & forms]
  (sputter:generate
    (sputter:for-each [form forms, elem form]
      (sputter:when (func elem) (emit elem)))))

(defn map
  {:doc-asset "map"}
  [func & forms]
  (sputter:generate
    (for-each [form forms, elem form] (emit (func elem)))))

(defmacro to-assoc
  {:doc-asset "to-assoc"}
  [& forms]
  `(apply assoc (concat ~@forms)))

(defmacro to-list
  {:doc-asset "to-list"}
  [& forms]
  `(apply list (concat ~@forms)))

(defmacro to-vector
  {:doc-asset "to-vector"}
  [& forms]
  `(apply vector (concat ~@forms)))
