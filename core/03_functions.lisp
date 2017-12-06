;;;; sputter core: functions

(defmacro assert-args
  ([]       nil)
  ([clause] (panic :message "assert-args clauses must be paired"))
  ([& clauses]
    `(cond
      ~(clauses 0) (assert-args ~@(rest (rest clauses)))
      :else        (panic :message ~(clauses 1)))))

(defmacro fn
  {:doc-asset "fn"}
  [& forms]
  (list 'sputter:make-closure
    (cond
      (and (> (len forms) 0) (is-vector (forms 0))) (forms 0)
      (and (> (len forms) 1) (is-vector (forms 1))) (forms 1)
      (and (> (len forms) 2) (is-vector (forms 2))) (forms 2)
      (and (> (len forms) 3) (is-vector (forms 3))) (forms 3)
      :else [])
    (cons 'sputter:lambda forms)))

(defmacro defn
  {:doc-asset "defn"}
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

(defmacro .
  [target method & args]
  `((get ~target ~method) ~@args))
