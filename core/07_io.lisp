;;;; sputter core: i/o

(defn pr-map-with-nil
  {:private true}
  [func seq]
  (map
    (fn [val] (if (is-nil val) val (func val)))
    seq))

(defn pr [& forms]
  (let [seq (pr-map-with-nil str! forms)]
    (if (is-seq seq)
      (. sputter:*stdout* :write (first seq)))
    (for-each [elem (rest seq)]
      (. sputter:*stdout* :write *space* elem))))

(defn prn [& forms]
  (apply pr forms)
  (. sputter:*stdout* :write *newline*))

(defn print [& forms]
  (let [seq (pr-map-with-nil str forms)]
    (if (is-seq seq)
      (. sputter:*stdout* :write (first seq)))
    (for-each [elem (rest seq)]
      (. sputter:*stdout* :write *space* elem))))

(defn println [& forms]
  (apply print forms)
  (. sputter:*stdout* :write *newline*))

(defn paired-vector?
  {:private true}
  [val]
  (and (is-vector val) (= (% (len val) 2) 0)))

(defmacro with-open [bindings & body]
  (assert-args
    (paired-vector? bindings) "with-open bindings must be a key-value vector")
  (cond
    (= (len bindings) 0)
      `(do ~@body)
    (>= (len bindings) 2)
      `(let [~(bindings 0) ~(bindings 1),
             close# (get ~(bindings 0) :close nil)]
        (try
          (with-open [~@(rest (rest bindings))] ~@body)
          (finally
            (when (is-apply close#) (close#)))))))
