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

;; namespaces

(def-builtin with-ns :doc-asset "with-ns" :special-form true)
(def-builtin ns      :doc-asset "ns" :special-form true)

;; basic predicates

(def-builtin eq    :doc-asset "eq")
(def-builtin !eq   :doc-asset "eq")
(def-builtin nil?  :doc-asset "is-nil")
(def-builtin !nil? :doc-asset "is-nil")
(def-builtin keyword?)
(def-builtin !keyword?)

;; strings

(def-builtin str   :doc-asset "str")
(def-builtin str?  :doc-asset "is-str")
(def-builtin !str? :doc-asset "is-str")

;; sequences

(def-builtin first :doc-asset "first")
(def-builtin rest  :doc-asset "rest")
(def-builtin cons  :doc-asset "cons")
(def-builtin conj  :doc-asset "conj")
(def-builtin len   :doc-asset "len")
(def-builtin nth   :doc-asset "nth")

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

(def-builtin +)
(def-builtin -)
(def-builtin *)
(def-builtin /)
(def-builtin =)
(def-builtin !=)
(def-builtin >)
(def-builtin >=)
(def-builtin <)
(def-builtin <=)

;; functions

(def-builtin closure      :doc-asset "closure" :macro true :special-form true)
(def-builtin lambda       :doc-asset "lambda" :special-form true)
(def-builtin apply        :doc-asset "apply")
(def-builtin make-closure :macro true)
(def-builtin apply?)
(def-builtin !apply?)

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

;; io

(def-builtin pr      :doc-asset "pr")
(def-builtin prn     :doc-asset "prn")
(def-builtin print   :doc-asset "print")
(def-builtin println :doc-asset "println")
