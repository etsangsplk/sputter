;;;; sputter core: threading

(defn thread-to-list
  {:private true}
  [target]
  (unless (list? target) (list target) target))

(defn make-threader
  {:private true}
  [func]
  (fn [value forms]
    (if (seq? forms)
      (let [v (func value (first forms)),
            f (rest forms)]
        (self v f))
      value)))

(def thread-first
  (make-threader
    (fn [value target]
      (let [l (thread-to-list target),
            f (first l),
            r (rest l)]
        (to-list [f value] r)))))

(def thread-last
  (make-threader
    (fn [value target]
      (let [l (thread-to-list target),
            f (first l),
            r (rest l)]
        (to-list [f] r [value])))))

(defmacro ->
  {:doc "threads value through a series of forms, as their first argument"}
  [value & forms]
  (thread-first value forms))

(defmacro ->>
  {:doc "threads value through a series of forms, as their last argument"}
  [value & forms]
  (thread-last value forms))
