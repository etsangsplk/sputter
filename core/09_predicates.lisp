;;;; sputter core: predicates

(defn predicate-lazy-seq
  {:private true}
  [func coll]
  (lazy-seq
    (if (is-seq coll)
      (let [f (func (first coll))]
        (if f
          (cons f (predicate-lazy-seq func (rest coll)))
          '(false)))
      '(true))))

(defmacro def-predicate-pos
  {:private true}
  [func name]
  (let [func-name (sym (str name "?"))]
    `(def ~func-name
      (with-meta
        (fn [~'first & ~'rest]
          (last (predicate-lazy-seq ~func (cons ~'first ~'rest))))
        (meta ~func)
        {:name ~(str func-name)}))))

(defmacro def-predicate-neg
  {:private true}
  [func name]
  (let [func-name (sym (str "!" name "?"))]
    `(def ~func-name
      (with-meta
        (fn [~'first & ~'rest]
          (let [nf# (fn [x] (not (~func x)))]
            (last (predicate-lazy-seq nf# (cons ~'first ~'rest)))))
        (meta ~func)
        {:name ~(str func-name)}))))

(defmacro def-predicate
  [func name]
  `(do
    (def-predicate-pos ~func ~name)
    (def-predicate-neg ~func ~name)))

(def-predicate is-nil "nil")
(def-predicate is-str "str")
(def-predicate is-seq "seq")
(def-predicate is-len "len")
(def-predicate is-indexed "indexed")
(def-predicate is-assoc "assoc")
(def-predicate is-mapped "mapped")
(def-predicate is-list "list")
(def-predicate is-vector "vector")
(def-predicate is-promise "promise")
(def-predicate is-meta "meta")
(def-predicate is-pos-inf "inf")
(def-predicate is-neg-inf "-inf")
(def-predicate is-nan "nan")
(def-predicate is-macro "macro")
(def-predicate is-apply "apply")
(def-predicate is-special-form "special-form")
(def-predicate is-keyword "keyword")
(def-predicate is-symbol "symbol")
(def-predicate is-local "local")

(defn is-even
  {:doc "will return whether or not the number is even"}
  [value]
  (= (% value 2) 0))

(defn is-odd
  {:doc "will return whether or not the number is odd"}
  [value]
  (= (% value 2) 1))

(def-predicate is-even "even")
(def-predicate is-odd "odd")
