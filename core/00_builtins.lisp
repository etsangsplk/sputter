;;;; sputter core: builtins

(def-builtin def  :doc-asset "def" :special-form true)
(def-builtin if   :doc-asset "if"  :special-form true)
(def-builtin let  :doc-asset "let" :special-form true)
(def-builtin do   :doc-asset "do"  :special-form true)
(def-builtin read :doc-asset "read")
(def-builtin eval :doc-asset "eval")

;; metadata

(def-builtin with-meta :doc-asset "with-meta")
(def-builtin meta      :doc-asset "meta")
(def-builtin meta?     :doc-asset "has-meta")
(def-builtin !meta?    :doc-asset "has-meta")

;; macros

(def-builtin defmacro        :doc-asset "defmacro" :special-form true)
(def-builtin macroexpand1    :special-form true)
(def-builtin macroexpand     :special-form true)
(def-builtin macroexpand-all :special-form true)
(def-builtin quote           :doc-asset "quote" :macro true :special-form true)
(def-builtin syntax-quote    :macro true)
(def-builtin gensym)
(def-builtin macro?)
(def-builtin !macro?)

;; namespaces

(def-builtin with-ns :doc-asset "with-ns" :special-form true)
(def-builtin ns      :doc-asset "ns"      :special-form true)

;; basic predicates

(def-builtin eq    :doc-asset "eq")
(def-builtin !eq   :doc-asset "eq")
(def-builtin nil?  :doc-asset "is-nil")
(def-builtin !nil? :doc-asset "is-nil")
(def-builtin keyword?)
(def-builtin !keyword?)
(def-builtin symbol?)
(def-builtin !symbol?)
(def-builtin local?)
(def-builtin !local?)

;; strings

(def-builtin str   :doc-asset "str")
(def-builtin str!  :doc-asset "str")
(def-builtin str?  :doc-asset "is-str")
(def-builtin !str? :doc-asset "is-str")

;; sequences

(def-builtin first :doc-asset "first")
(def-builtin rest  :doc-asset "rest")
(def-builtin cons  :doc-asset "cons")
(def-builtin conj  :doc-asset "conj")
(def-builtin len   :doc-asset "len")
(def-builtin nth   :doc-asset "nth")
(def-builtin get   :doc-asset "get")

;; predicates

(def-builtin seq?      :doc-asset "is-seq")
(def-builtin !seq?     :doc-asset "is-seq")
(def-builtin len?      :doc-asset "is-len")
(def-builtin !len?     :doc-asset "is-len")
(def-builtin indexed?  :doc-asset "is-indexed")
(def-builtin !indexed? :doc-asset "is-indexed")

;; associatives

(def-builtin assoc    :doc-asset "assoc")
(def-builtin assoc?   :doc-asset "is-assoc")
(def-builtin !assoc?  :doc-asset "is-assoc")
(def-builtin mapped?  :doc-asset "is-mapped")
(def-builtin !mapped? :doc-asset "is-mapped")

;; lists

(def-builtin list   :doc-asset "list")
(def-builtin list?  :doc-asset "is-list")
(def-builtin !list? :doc-asset "is-list")

;; vectors

(def-builtin vector   :doc-asset "vector")
(def-builtin vector?  :doc-asset "is-vector")
(def-builtin !vector? :doc-asset "is-vector")

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

(def-builtin inf?   :doc "checks a number for positive infinity")
(def-builtin !inf?  :doc "checks a number for positive infinity")
(def-builtin -inf?  :doc "checks a number for negative infinity")
(def-builtin !-inf? :doc "checks a number for negative infinity")
(def-builtin nan?   :doc "checks that a value is not a number")
(def-builtin !nan?  :doc "checks that a value is not a number")

;; functions

(def-builtin closure      :doc-asset "closure" :macro true :special-form true)
(def-builtin lambda       :doc-asset "lambda"  :special-form true)
(def-builtin apply        :doc-asset "apply")
(def-builtin make-closure :macro true)
(def-builtin apply?)
(def-builtin !apply?)
(def-builtin special-form?)
(def-builtin !special-form?)

;; concurrency

(def-builtin make-go   :special-form true)
(def-builtin chan      :doc-asset "chan")
(def-builtin promise   :doc-asset "promise")
(def-builtin promise?  :doc-asset "is-promise")
(def-builtin !promise? :doc-asset "is-promise")

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
