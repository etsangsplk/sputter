(defmacro when
  [test & body]
  (list 'sputter:if test (cons 'sputter:do body)))

(defmacro when-not
  [test & body]
  (list 'sputter:if test nil (cons 'sputter:do body)))

(defmacro cond
  {:doc-asset "cond"}
  [& clauses]
  (when (seq? clauses)
    (if (= 1 (len clauses))
      (clauses 0)
      (list 'sputter:if
        (clauses 0) (clauses 1)
        (cons 'cond (rest (rest clauses)))))))
