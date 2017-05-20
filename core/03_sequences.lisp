
(defmacro to-assoc
  {:doc-asset "to-assoc"}
  [& forms]
  (list 'sputter:apply 'sputter:assoc
    (cons 'sputter:concat forms)))

(defmacro to-list
  {:doc-asset "to-list"}
  [& forms]
  (list 'sputter:apply 'sputter:list
    (cons 'sputter:concat forms)))

(defmacro to-vector
  {:doc-asset "to-vector"}
  [& forms]
  (list 'sputter:apply 'sputter:vector
    (cons 'sputter:concat forms)))
