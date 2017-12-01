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

(defmacro some->
  {:doc "like `->`, but returns _nil_ if any form evaluates as such"}
  ([value] value)
  ([value & forms]
    (let [l (thread-to-list (first forms)),
          f (first l),
          r (rest l)]
      `(let [val# ~value]
        (when-not (nil? val#)
          (some-> (~f val# ~@r) ~@(rest forms)))))))

(defmacro some->>
  {:doc "like `->>`, but returns _nil_ if any form evaluates as such"}
  ([value] value)
  ([value & forms]
    (let [l (thread-to-list (first forms)),
          f (first l),
          r (rest l)]
      `(let [val# ~value]
        (when-not (nil? val#)
          (some->> (~f ~@r val#) ~@(rest forms)))))))

(defmacro as->
  {:doc "threads using a bound name for positional flexibility"}
  ([value name] value)
  ([value name & forms]
    (let [l (thread-to-list (first forms))]
      `(let [~name ~value]
        (as-> ~l ~name ~@(rest forms))))))
