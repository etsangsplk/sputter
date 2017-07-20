;;;; sputter core: concurrency

(defmacro go
  {:doc-asset "go"}
  [& forms]
  (list 'sputter:make-closure []
    (cons 'sputter:make-go forms)))

(defmacro generate
  {:doc-asset "generate"}
  [& forms]
  `(let [ch# (chan) cl# (:close ch#) emit (:emit ch#)]
    (go (let [x (do ~@forms)] (cl#) x))
    (:seq ch#)))

(defmacro future
  {:doc-asset "future"}
  [& forms]
  `(let [promise# (promise)]
    (go (promise# (do ~@forms)))
    promise#))
