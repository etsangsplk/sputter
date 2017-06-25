; sputter core: branching

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
  (sputter:when (seq? clauses)
    (if (= 1 (len clauses))
      (clauses 0)
      (list 'sputter:if
        (clauses 0) (clauses 1)
        (cons 'cond (rest (rest clauses)))))))

(defmacro and
  {:doc-asset "and"}
  [& clauses]
  (sputter:cond
    (= 0 (len clauses)) true
    (= 1 (len clauses)) (clauses 0)
    :else
      (list 'sputter:let
        (vector 'sputter/and (clauses 0))
        (list 'sputter:if
          'sputter/and
          (cons 'and (rest clauses))
          'sputter/and))))

(defmacro or
  {:doc-asset "or"}
  [& clauses]
  (sputter:cond
    (= 0 (len clauses)) true
    (= 1 (len clauses)) (clauses 0)
    :else
      (list 'sputter:let
        (vector 'sputter/or (clauses 0))
        (list 'sputter:if
          'sputter/or 'sputter/or
          (cons 'or (rest clauses))))))
