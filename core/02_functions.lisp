;;;; sputter core: functions

(defmacro assert-args
  [& clauses]
  (cond
    (!seq? clauses)
      nil
    (= (len clauses) 1)
      (panic :message "assert-args clauses must be paired")
    (>= (len clauses) 2)
      `(cond
        ~(clauses 0) (sputter:assert-args ~@(rest (rest clauses)))
        :else        (panic :message ~(clauses 1)))))

(defmacro fn
  {:doc-asset "fn"}
  [& forms]
  (list 'sputter:make-closure
    (conj (cond
      (and (> (len forms) 0) (vector? (forms 0))) (forms 0)
      (and (> (len forms) 1) (vector? (forms 1))) (forms 1)
      (and (> (len forms) 2) (vector? (forms 2))) (forms 2)
      (and (> (len forms) 3) (vector? (forms 3))) (forms 3)
      :else []) 'recur)
    (cons 'sputter:lambda forms)))

(defmacro defn
  {:doc-asset "defn"}
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

(defmacro .
  [target method & args]
  `((get ~target ~method) ~@args))
