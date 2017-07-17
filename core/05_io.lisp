;;;; sputter core: i/o

(defn pr-map-with-nil [f s]
  (map (lambda [v]
    (if (nil? v) v (f v))) s))

(defn pr [& forms]
  (let [s (pr-map-with-nil str! forms)]
    (if (seq? s) (. sputter:*stdout* :write (first s)))
    (for-each [e (rest s)]
      (. sputter:*stdout* :write sputter:*space* e))))

(defn prn [& forms]
  (apply pr forms)
  (. sputter:*stdout* :write sputter:*newline*))

(defn print [& forms]
  (let [s (pr-map-with-nil str forms)]
    (if (seq? s) (. sputter:*stdout* :write (first s)))
    (for-each [e (rest s)]
      (. sputter:*stdout* :write sputter:*space* e))))

(defn println [& forms]
  (apply print forms)
  (. sputter:*stdout* :write sputter:*newline*))
