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
  ([]             `(range* 0 inf 1))
  ([max]          `(range* 0 ~max 1))
  ([min max]      `(range* ~min ~max 1))
  ([min max step] `(range* ~min ~max ~step)))

(defmacro lazy-seq
  {:doc-asset "lazy-seq"}
  [& body]
  (cons 'sputter:lazy-seq* body))

(defn take-while
  [pred coll]
  (lazy-seq
    (when-let [s coll]
      (when (pred (first s))
        (cons (first s)
              (take-while pred (rest s)))))))

(defmacro for
  [seq-exprs & body]
  `(generate
    (for-each ~seq-exprs (emit (do ~@body)))))

(defn partition
  {:doc-asset "partition"}
  ([count coll] (partition count count coll))
  ([count step coll]
    (lazy-seq
      (when (is-seq coll)
        (cons (to-list (take count coll))
              (partition count step (drop step coll)))))))
