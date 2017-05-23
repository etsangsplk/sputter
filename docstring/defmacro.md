# (defmacro name meta? [name**] form+) assigns a reader macro
Will declare a macro in the current namespace, which is 'user' by default. A macro is expanded by the reader in order to alter the source code's data representation before it is evaluated.

## An Example

  (defmacro cond
    "implements the 'cond' form"
    [& clauses]
    (when (seq? clauses)
      (if (= 1 (len clauses))
        (clauses 0)
        (list 'sputter:if
          (clauses 0) (clauses 1)
          (cons 'cond (rest (rest clauses)))))))
