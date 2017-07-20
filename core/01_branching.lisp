;;;; sputter core: branching

(defmacro !
  {:doc "logically inverts the truthiness of the provided value"}
  [val]
  `(if ~val false true))

(defmacro when
  {:doc-asset "when"}
  [test & body]
  `(if ~test (do ~@body) nil))

(defmacro when-not
  {:doc-asset "when"}
  [test & body]
  `(if ~test nil (do ~@body)))

(defmacro cond
  {:doc-asset "cond"}
  [& clauses]
  (when (seq? clauses)
    (if (= (len clauses) 1)
      (clauses 0)
      `(if ~(clauses 0) ~(clauses 1)
        (sputter:cond ~@(rest (rest clauses)))))))

(defmacro and
  {:doc-asset "and"}
  [& clauses]
  (cond
    (!seq? clauses)     true
    (= (len clauses) 1) (clauses 0)
    :else
      `(let [and# ~(clauses 0)]
        (if and# (sputter:and ~@(rest clauses)) and#))))

(defmacro or
  {:doc-asset "or"}
  [& clauses]
  (cond
    (!seq? clauses)     nil
    (= (len clauses) 1) (clauses 0)
    :else
      `(let [or# ~(clauses 0)]
        (if or# or# (sputter:or ~@(rest clauses))))))
