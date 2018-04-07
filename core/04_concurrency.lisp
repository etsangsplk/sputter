;;;; sputter core: concurrency

(defmacro go
  {:doc-asset "go"}
  [& body]
  (cons 'sputter:go* body))

(defmacro generate
  {:doc-asset "generate"}
  [& body]
  `(let [chan# (chan),
         close# (:close chan#),
         emit (:emit chan#)]
    (go (let [result# (do ~@body)]
          (close#)
          result#))
    (:seq chan#)))

(defmacro future
  {:doc-asset "future"}
  [& body]
  `(let [promise# (promise)]
    (go (promise# (do ~@body)))
    promise#))
