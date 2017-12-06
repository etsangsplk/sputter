;;;; sputter core: predicates

(defn lazy-predicate
  {:private true}
  [func coll]
  (lazy-seq
    (if (is-seq coll)
      (let [f (func (first coll))]
        (if f
          (cons f (lazy-predicate func (rest coll)))
          '(false)))
      '(true))))

(defmacro def-pos-predicate
  {:private true}
  [name func-name]
  `(def ~name
    (with-meta
      (fn [~'first & ~'rest]
        (last (lazy-predicate ~func-name (cons ~'first ~'rest))))
      (meta ~func-name)
      {:name ~(str name)})))

(defmacro def-neg-predicate
  {:private true}
  [name func-name]
  `(def ~name
    (with-meta
      (fn [~'first & ~'rest]
        (let [nf# (fn [x] (not (~func-name x)))]
          (last (lazy-predicate nf# (cons ~'first ~'rest)))))
      (meta ~func-name)
      {:name ~(str name)})))

(defmacro def-predicate
  [pos-name neg-name func-name]
  `(do
    (def-pos-predicate ~pos-name ~func-name)
    (def-neg-predicate ~neg-name ~func-name)))

(defn is-even
  {:doc "will return whether or not the number is even"}
  [value]
  (= (% value 2) 0))

(defn is-odd
  {:doc "will return whether or not the number is odd"}
  [value]
  (!= (% value 2) 0))

(def-predicate even? !even? is-even)
(def-predicate odd? !odd? is-odd)

(def-predicate nil? !nil? is-nil)
(def-predicate str? !str? is-str)
(def-predicate seq? !seq? is-seq)
(def-predicate len? !len? is-len)
(def-predicate indexed? !indexed? is-indexed)
(def-predicate assoc? !assoc? is-assoc)
(def-predicate mapped? !mapped? is-mapped)
(def-predicate list? !list? is-list)
(def-predicate vector? !vector? is-vector)
(def-predicate promise? !promise? is-promise)
(def-predicate meta? !meta? is-meta)
(def-predicate inf? !inf? is-pos-inf)
(def-predicate -inf? !-inf? is-neg-inf)
(def-predicate nan? !nan? is-nan)
(def-predicate macro? !macro? is-macro)
(def-predicate apply? !apply? is-apply)
(def-predicate special-form? !special-form? is-special-form)
(def-predicate keyword? !keyword? is-keyword)
(def-predicate symbol? !symbol? is-symbol)
(def-predicate local? !local? is-local)
