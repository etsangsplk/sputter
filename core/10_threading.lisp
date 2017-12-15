;;;; sputter core: threading

(defn thread-to-list
  {:private true}
  [target]
  (unless (is-list target)
    (list target)
    target))

(defmacro ->
  {:doc-asset "thread-first"}
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
        (when-not (is-nil val#)
          (some-> (~f val# ~@r) ~@(rest forms)))))))

(defmacro some->>
  {:doc "like `->>`, but returns _nil_ if any form evaluates as such"}
  ([value] value)
  ([value & forms]
    (let [l (thread-to-list (first forms)),
          f (first l),
          r (rest l)]
      `(let [val# ~value]
        (when-not (is-nil val#)
          (some->> (~f ~@r val#) ~@(rest forms)))))))

(defmacro as->
  {:doc "threads using a bound name for positional flexibility"}
  ([value name] value)
  ([value name & forms]
    (let [l (thread-to-list (first forms))]
      `(let [~name ~value]
        (as-> ~l ~name ~@(rest forms))))))

(defn make-cond-clause
  {:private true :special-form true}
  [sym]
  (fn [clause]
    (let [pred (nth clause 0), form (nth clause 1)]
      `((sputter:fn [val] (if ~pred (~sym val ~form) val))))))

(defmacro cond->
  {:doc "conditionally threads through a set of forms as their first argument"}
  ([value] value)
  ([value & clauses]
    (assert-args
      (is-even (len clauses)) "clauses must be paired")
    `(-> ~value ~@(map (make-cond-clause ->) (partition 2 clauses)))))

(defmacro cond->>
  {:doc "conditionally threads through a set of forms as their last argument"}
  ([value] value)
  ([value & clauses]
    (assert-args
      (is-even (len clauses)) "clauses must be paired")
    `(-> ~value ~@(map (make-cond-clause ->>) (partition 2 clauses)))))
