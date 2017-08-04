;;;; sputter core: concurrency

(defmacro go
  {:doc-asset "go"}
  [& body]
  (list 'sputter:make-closure []
    (cons 'sputter:make-go body)))

(defmacro generate
  {:doc-asset "generate"}
  [& body]
  `(let [ch# (chan) cl# (:close ch#) emit (:emit ch#)]
    (go (let [x (do ~@body)] (cl#) x))
    (:seq ch#)))

(defmacro future
  {:doc-asset "future"}
  [& body]
  `(let [promise# (promise)]
    (go (promise# (do ~@body)))
    promise#))
