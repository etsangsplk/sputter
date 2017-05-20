(defmacro defn
  {:doc-asset "defn"}
  [name & rest]
  (list 'sputter:def name
    (cons 'sputter:fn (cons name rest))))
