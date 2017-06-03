; sputter core: definitions

(defmacro fn
  {:doc-asset "fn"}
  [& forms]
  (list 'sputter:closure
    (sputter:cond
      (and (> (len forms) 0) (vector? (forms 0))) (forms 0)
      (and (> (len forms) 1) (vector? (forms 1))) (forms 1)
      (and (> (len forms) 2) (vector? (forms 2))) (forms 2)
      (and (> (len forms) 3) (vector? (forms 3))) (forms 3)
      :else (vector))
    (cons 'sputter:make-fn forms)))

(defmacro defn
  {:doc-asset "defn"}
  [name & forms]
  (list 'sputter:def name
    (cons 'fn (cons name forms))))
