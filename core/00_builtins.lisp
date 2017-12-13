;;;; sputter core: builtins

(def-builtin if    :doc-asset "if"  :special-form true)
(def-builtin let   :doc-asset "let" :special-form true)
(def-builtin do    :doc-asset "do"  :special-form true)
(def-builtin read  :doc-asset "read")
(def-builtin eval  :doc-asset "eval")
(def-builtin is-eq :doc-asset "is-eq")

;; basic predicates

(def-builtin is-nil)
(def-builtin is-keyword)

;; metadata

(def-builtin with-meta :doc-asset "with-meta")
(def-builtin meta      :doc-asset "meta")
(def-builtin is-meta)

;; macros

(def-builtin defmacro     :doc-asset "defmacro" :special-form true)
(def-builtin quote        :doc-asset "quote" :macro true :special-form true)
(def-builtin syntax-quote :macro true)
(def-builtin macroexpand1)
(def-builtin macroexpand)
(def-builtin macroexpand-all)
(def-builtin is-macro)

;; symbols

(def-builtin sym       :doc-asset "sym")
(def-builtin gensym    :doc-asset "gensym")
(def-builtin is-symbol)
(def-builtin is-local)

;; namespaces

(def-builtin with-ns :doc-asset "with-ns" :special-form true)
(def-builtin ns      :doc-asset "ns"      :special-form true)
(def-builtin ns-put  :doc-asset "ns-put"  :special-form true)

;; strings

(def-builtin str    :doc-asset "str")
(def-builtin str!   :doc-asset "str")
(def-builtin is-str)

;; sequences

(def-builtin seq    :doc "attempts to convert value to a seq, nil otherwise")
(def-builtin first  :doc-asset "first")
(def-builtin rest   :doc-asset "rest")
(def-builtin last   :doc-asset "last")
(def-builtin cons   :doc-asset "cons")
(def-builtin conj   :doc-asset "conj")
(def-builtin len    :doc-asset "len")
(def-builtin nth    :doc-asset "nth")
(def-builtin get    :doc-asset "get")
(def-builtin assoc  :doc-asset "assoc")
(def-builtin list   :doc-asset "list")
(def-builtin vector :doc-asset "vector")

(def-builtin is-seq)
(def-builtin is-len)
(def-builtin is-indexed)
(def-builtin is-assoc)
(def-builtin is-mapped)
(def-builtin is-list)
(def-builtin is-vector)

;; numeric

(def-builtin inc :doc "increases a number by one")
(def-builtin dec :doc "decreases a number by one")
(def-builtin +   :doc "adds a set of numbers")
(def-builtin -   :doc "subtracts a set of numbers")
(def-builtin *   :doc "multiplies a set of numbers")
(def-builtin /   :doc "divides a set of numbers")
(def-builtin %   :doc "produces the remainder for a divided set of numbers")

(def-builtin =  :doc "checks a set of numbers for equality")
(def-builtin != :doc "checks a set of numbers for inequality")
(def-builtin >  :doc "checks that a set of numbers increases in value")
(def-builtin >= :doc "checks that a set of numbers doesn't decrease in value")
(def-builtin <  :doc "checks that a set of numbers decreases in value")
(def-builtin <= :doc "checks that a set of numbers doesn't increase in value")

(def-builtin is-pos-inf)
(def-builtin is-neg-inf)
(def-builtin is-nan)

;; functions

(def-builtin closure      :doc-asset "closure" :macro true :special-form true)
(def-builtin lambda       :doc-asset "lambda"  :special-form true)
(def-builtin apply        :doc-asset "apply")
(def-builtin partial      :doc "generates a function based on a partial apply")
(def-builtin make-closure :macro true)
(def-builtin is-apply)
(def-builtin is-special-form)

;; concurrency

(def-builtin make-go    :special-form true)
(def-builtin chan       :doc-asset "chan")
(def-builtin promise    :doc-asset "promise")
(def-builtin is-promise)

;; lazy sequences

(def-builtin make-lazy-seq :special-form true)
(def-builtin for-each      :special-form true)
(def-builtin concat        :doc-asset "concat")
(def-builtin filter        :doc-asset "filter")
(def-builtin map           :doc-asset "map")
(def-builtin reduce        :doc-asset "reduce")
(def-builtin take          :doc-asset "take")
(def-builtin drop          :doc-asset "drop")
(def-builtin make-range)

;; raise and recover

(def-builtin make-error)
(def-builtin raise)
(def-builtin recover :special-form true)

;; current time

(def-builtin current-time)
