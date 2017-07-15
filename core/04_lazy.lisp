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
  (cond
    (= (len forms) 0) (list 'sputter:make-range 0 'sputter:inf 1)
    (= (len forms) 1) (list 'sputter:make-range 0 (forms 0) 1)
    (= (len forms) 2) (list 'sputter:make-range (forms 0) (forms 1) 1)
    :else             (cons 'sputter:make-range forms)))

(defmacro lazy-seq
  [& body]
  (list 'sputter:make-closure
    (vector)
    (cons 'sputter:make-lazy-seq body)))

(defmacro for
  [seq-exprs & body]
  `(generate
     (for-each ~seq-exprs (emit (do ~@body)))))
