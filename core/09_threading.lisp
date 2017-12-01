;;;; sputter core: threading

(defn thread-to-list
  {:private true}
  [target]
  (unless (list? target) (list target) target))

(defmacro ->
  {:doc "threads value through a series of forms, as their first argument"}
  ([value] value)
  ([value & forms]
    (let [l (thread-to-list (first forms)),
          f (first l),
          r (rest l)]
      `(-> (~f ~value ~@r) ~@(rest forms)))))

(defmacro ->>
  {:doc "threads value through a series of forms, as their last argument"}
  ([value] value)
  ([value & forms]
    (let [l (thread-to-list (first forms)),
          f (first l),
          r (rest l)]
      `(->> (~f ~@r ~value) ~@(rest forms)))))
