;;;; sputter core: i/o

(defn pr-map-with-nil
  {:private true}
  [func seq]
  (map
    (fn [val] (if (nil? val) val (func val)))
    seq))

(defn pr [& forms]
  (let [s (pr-map-with-nil str! forms)]
    (if (seq? s) (. sputter:*stdout* :write (first s)))
    (for-each [e (rest s)]
      (. sputter:*stdout* :write *space* e))))

(defn prn [& forms]
  (apply pr forms)
  (. sputter:*stdout* :write *newline*))

(defn print [& forms]
  (let [s (pr-map-with-nil str forms)]
    (if (seq? s) (. sputter:*stdout* :write (first s)))
    (for-each [e (rest s)]
      (. sputter:*stdout* :write *space* e))))

(defn println [& forms]
  (apply print forms)
  (. sputter:*stdout* :write *newline*))

(defn paired-vector?
  {:private true}
  [val]
  (and (vector? val) (= (% (len val) 2) 0)))

(defmacro with-open [bindings & body]
  (assert-args
    (paired-vector? bindings) "with-open bindings must be a paired vector")
  (cond
    (= (len bindings) 0)
      `(sputter:do ~@body)
    (>= (len bindings) 2)
      `(let [~(bindings 0) ~(bindings 1),
             cl# (get ~(bindings 0) :close nil),
             res# (sputter:with-open [~@(rest (rest bindings))] ~@body)]
        (when (apply? cl#) (cl#))
        res#)))
