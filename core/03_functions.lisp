;;;; sputter core: functions

(defmacro assert-args
  ([]       nil)
  ([clause] (panic :message "assert-args clauses must be paired"))
  ([& clauses]
    `(cond
      ~(clauses 0) (assert-args ~@(rest (rest clauses)))
      :else        (panic :message ~(clauses 1)))))

(defmacro defn
  {:doc-asset "defn"}
  [name & forms]
  `(def ~name (fn ~name ~@forms)))

(defmacro .
  [target method & args]
  `((get ~target ~method) ~@args))

(defn no-op
  {:doc "a function that does absolutely nothing"}
  [])

(defn identity
  {:doc "a function that returns its single argument"}
  [val]
  val)

(defn constantly
  {:doc "returns a function that always returns the provided value"}
  [val]
  (fn [] val))

(defn comp-outer
  {:private true}
  [func args rest-funcs]
  (if (is-seq rest-funcs)
    (comp-outer (first rest-funcs) (list func args) (rest rest-funcs))
    (list func args)))

(defmacro comp
  ([] identity)
  ([func] func)
  ([func & funcs]
    (let [args        (gensym "args")
          inner       (list 'apply func args)
          first-outer (first funcs)
          rest-outer  (rest funcs)]
      `(fn [& ~args]
        ~(comp-outer first-outer inner rest-outer)))))

(defmacro juxt
  [& funcs]
  (let [args (gensym "args")]
    `(fn [& ~args]
      [~@(map (fn [f] (list 'apply f args)) funcs)])))
