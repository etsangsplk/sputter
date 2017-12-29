;;;; sputter core: predicates

(defn is-even [value] (= (mod value 2) 0))
(defn is-odd [value] (= (mod value 2) 1))

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
  [func name metadata]
  (let [func-name (sym (str name "?"))]
    `(def ~func-name
      (with-meta
        (fn ~func-name [~'first & ~'rest]
          (last (predicate-lazy-seq ~func (cons ~'first ~'rest))))
        ~metadata))))

(defmacro def-predicate-neg
  {:private true}
  [func name metadata]
  (let [func-name (sym (str "!" name "?"))]
    `(def ~func-name
      (with-meta
        (fn ~func-name [~'first & ~'rest]
          (let [nf# (fn [x] (not (~func x)))]
            (last (predicate-lazy-seq nf# (cons ~'first ~'rest)))))
        ~metadata))))

(defmacro def-predicate
  [func name & metadata]
  (assert-args
    (is-even (len metadata)) "metadata must be paired")
  (let [md (to-assoc metadata)]
    `(do
      (def-predicate-pos ~func ~name ~md)
      (def-predicate-neg ~func ~name ~md))))

(def-predicate is-atom "atom"       :doc-asset "is-atom")
(def-predicate is-nil "nil"         :doc-asset "is-nil")
(def-predicate is-str "str"         :doc-asset "is-str")
(def-predicate is-seq "seq"         :doc-asset "is-seq")
(def-predicate is-len "len"         :doc-asset "is-len")
(def-predicate is-indexed "indexed" :doc-asset "is-indexed")
(def-predicate is-assoc "assoc"     :doc-asset "is-assoc")
(def-predicate is-mapped "mapped"   :doc-asset "is-mapped")
(def-predicate is-list "list"       :doc-asset "is-list")
(def-predicate is-vector "vector"   :doc-asset "is-vector")

(def-predicate is-promise "promise" :doc-asset "is-promise")
(def-predicate is-meta "meta"       :doc-asset "has-meta")

(def-predicate is-pos-inf "inf"  :doc "checks a number for positive infinity")
(def-predicate is-neg-inf "-inf" :doc "checks a number for negative infinity")
(def-predicate is-nan "nan"      :doc "checks that a value is not a number")

(def-predicate is-macro "macro" :doc "tests if a value is a macro")
(def-predicate is-apply "apply" :doc "tests if a value can be applied")
(def-predicate is-special-form "special-form" :doc "tests for a special form")

(def-predicate is-keyword "keyword" :doc "tests if a value is a keyword")
(def-predicate is-symbol "symbol"   :doc "tests if a value is a symbol")
(def-predicate is-local "local"     :doc "tests if a value is an unqualified symbol")

(def-predicate is-even "even" :doc "tests if a number is even")
(def-predicate is-odd "odd"   :doc "tests if a number is odd")
