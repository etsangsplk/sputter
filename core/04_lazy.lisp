;;;; sputter core: lazy sequences

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

(defmacro range
  {:doc-asset "range"}
  [& forms]
  (assert-args
    (<= (len forms) 3) "requires between 0 and 3 arguments")
  (cond
    (= (len forms) 0) (list 'sputter:make-range 0 'sputter:inf 1)
    (= (len forms) 1) (list 'sputter:make-range 0 (forms 0) 1)
    (= (len forms) 2) (list 'sputter:make-range (forms 0) (forms 1) 1)
    :else             (cons 'sputter:make-range forms)))

(defmacro lazy-seq
  [& body]
  (list 'sputter:make-closure []
    (cons 'sputter:make-lazy-seq body)))

(defmacro when-let
  [bindings & body]
  (assert-args
    (vector? bindings)   "a vector for its binding"
    (= 2 (len bindings)) "exactly 2 forms in binding vector")
  (let [form (bindings 0) test (bindings 1)]
    `(let [temp# ~test]
      (when temp# (let [~form temp#] ~@body)))))

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
