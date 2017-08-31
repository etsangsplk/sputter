;;;; sputter core: os

(defmacro time
  [& forms]
  `(let [start# (current-time)
         result# (do ~@forms)
         end# (current-time)]
    (println (/ (- end# start#) 1000000) "ms")
    result#))
