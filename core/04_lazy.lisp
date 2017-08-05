;;;; sputter core: lazy sequences

(defmacro to-assoc
  {:doc-asset "to-assoc"}
  [& seqs]
  `(apply assoc (concat ~@seqs)))

(defmacro to-list
  {:doc-asset "to-list"}
  [& seqs]
  `(apply list (concat ~@seqs)))

(defmacro to-vector
  {:doc-asset "to-vector"}
  [& seqs]
  `(apply vector (concat ~@seqs)))

(defmacro range
  {:doc-asset "range"}
  ([]             `(sputter:make-range 0 sputter:inf 1))
  ([max]          `(sputter:make-range 0 ~max 1))
  ([min max]      `(sputter:make-range ~min ~max 1))
  ([min max step] `(sputter:make-range ~min ~max ~step)))

(defmacro lazy-seq
  {:doc-asset "lazy-seq"}
  [& body]
  (list 'sputter:make-closure []
    (cons 'sputter:make-lazy-seq body)))

(defmacro when-let
  [bindings & body]
  (assert-args
    (vector? bindings)   "binding vector must be supplied"
    (= 2 (len bindings)) "binding vector must contain 2 elements")
  (let [form (bindings 0),
        test (bindings 1)]
    `(let [temp# ~test]
      (when temp#
        (let [~form temp#] ~@body)))))

(defn take-while
  [pred coll]
  (assert-args
    (apply? pred) "predicate must be supplied"
    (seq? coll)   "collection must be supplied")
  (lazy-seq
    (when-let [s coll]
      (when (pred (first s))
        (cons (first s) (sputter:take-while pred (rest s)))))))

(defmacro for
  [seq-exprs & body]
  `(generate
     (for-each ~seq-exprs (emit (do ~@body)))))
