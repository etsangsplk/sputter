;;;; sputter core: i/o

(defn pr-map-with-nil [f s]
    (map (lambda [v]
        (if (nil? v) v (f v))) s))

(defn pr [& forms]
    (let [s (pr-map-with-nil str! forms)]
        (if (seq? s) (. *stdout* :write (first s)))
        (for-each [e (rest s)]
            (. *stdout* :write " " e))))

(defn prn [& forms]
    (apply pr forms)
    (. *stdout* :write *newline*))

(defn print [& forms]
    (let [s (pr-map-with-nil str forms)]
        (if (seq? s) (. *stdout* :write (first s)))
        (for-each [e (rest s)]
            (. *stdout* :write " " e))))

(defn println [& forms]
    (apply print forms)
    (. *stdout* :write *newline*))
