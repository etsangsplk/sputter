; sputter core: sequences

(defmacro lazy-seq [& forms]
  (cons 'sputter:async forms))

(defn concat
  {:doc-asset "concat"}
  [& forms]
  (lazy-seq
    (for-each [f forms, e f] (emit e))))

(defn filter
  {:doc-asset "filter"}
  [fn & forms]
  (lazy-seq
    (for-each [f forms, e f]
      (when (fn e) (emit e)))))

(defn map
  {:doc-asset "map"}
  [fn & forms]
  (lazy-seq
    (for-each [f forms, e f] (emit (fn e)))))

(defmacro to-assoc
  {:doc-asset "to-assoc"}
  [& forms]
  (list 'sputter:apply 'sputter:assoc
    (cons 'concat forms)))

(defmacro to-list
  {:doc-asset "to-list"}
  [& forms]
  (list 'sputter:apply 'sputter:list
    (cons 'concat forms)))

(defmacro to-vector
  {:doc-asset "to-vector"}
  [& forms]
  (list 'sputter:apply 'sputter:vector
    (cons 'concat forms)))
