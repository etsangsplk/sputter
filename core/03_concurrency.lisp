;;;; sputter core: concurrency

(defmacro go
  {:doc-asset "go"}
  [& forms]
  (list 'sputter:make-closure
    (vector)
    (cons 'sputter:make-go forms)))

(defmacro generate
  {:doc-asset "generate"}
  [& forms]
  (list 'sputter:let
    (vector 'sputter/ch (list 'sputter:chan)
            'sputter/cl (list :close 'sputter/ch)
            'emit (list :emit 'sputter/ch))
    (list 'sputter:go
      (list 'sputter:let (vector 'x (cons 'sputter:do forms))
        (list 'sputter/cl)
        'x))
    (list :seq 'sputter/ch)))

(defmacro future
  {:doc-asset "future"}
  [& forms]
  `(let [promise# (promise)]
    (go (promise# (do ~@forms)))
    promise#))
