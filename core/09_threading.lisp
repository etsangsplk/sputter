;;;; sputter core: threading

(defn thread-to-list
  {:private true}
  [target]
  (unless (list? target) (list target) target))

(defn thread-first
  {:private true}
  [value target]
  (let [l (thread-to-list target),
        f (first l),
        r (rest l)]
    (apply list (concat [f value] r))))

(defn thread-last
  {:private true}
  [value target]
  (let [l (thread-to-list target),
        f (first l),
        r (rest l)]
    (apply list (concat [f] r [value]))))

(defn thread
  {:private true}
  [func value forms]
  (if (seq? forms)
    (let [v (func value (first forms)),
          f (rest forms)]
      (sputter:thread func v f))
    value))

(defmacro ->
  {:doc "threads value through a series of forms, as their first argument"}
  [value & forms]
  (thread thread-first value forms))

(defmacro ->>
  {:doc "threads value through a series of forms, as their last argument"}
  [value & forms]
  (thread thread-last value forms))
