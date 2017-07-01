;;;; sputter core: functions

(defmacro fn
  {:doc-asset "fn"}
  [& forms]
  (list 'sputter:make-closure
    (cond
      (and (> (len forms) 0) (vector? (forms 0))) (forms 0)
      (and (> (len forms) 1) (vector? (forms 1))) (forms 1)
      (and (> (len forms) 2) (vector? (forms 2))) (forms 2)
      (and (> (len forms) 3) (vector? (forms 3))) (forms 3)
      :else (vector))
    (cons 'sputter:lambda forms)))

(defmacro defn
  {:doc-asset "defn"}
  [name & forms]
  `(def ~name (fn ~name ~@forms)))
