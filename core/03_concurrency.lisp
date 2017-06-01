
; sputter core: concurrency

(defmacro async
  {:doc-asset "async"}
  [& forms]
  (list 'sputter:closure
    (vector)
    (cons 'sputter:do-async forms)))

(defmacro generate
  {:doc-asset "generate"}
  [& forms]
  (list 'sputter:let
    (vector 'sputter/ch (list 'sputter:channel)
            'sputter/cl (list :close 'sputter/ch)
            'emit (list :emit 'sputter/ch))
    (list 'async
      (list 'sputter:let (vector 'x (cons 'sputter:do forms))
        (list 'sputter/cl)
        'x))
    (list :seq 'sputter/ch)))

(defmacro future
  {:doc-asset "future"}
  [& forms]
  (list 'sputter:let
    (vector 'sputter/pr (list 'sputter:promise))
    (list 'async
      (list 'sputter/pr (cons 'sputter:do forms)))
    'sputter/pr))
