;;;; sputter core: exceptions

(defmacro error
  [& clauses]
  `(make-error (assoc ~@clauses)))

(defmacro panic
  [& clauses]
  `(raise (error ~@clauses)))

(defn is-call
  {:private true}
  [sym clause]
  (and (local? sym)
       (list? clause)
       (eq sym (first clause))))

(defn is-catch-binding
  {:private true}
  [form]
  (and (vector? form)
       (= 2 (len form))
       (local? (form 1))))

(defn is-catch
  {:private true}
  [clause parsed]
  (and (is-call 'catch clause)
       (is-catch-binding (nth clause 1))
       (!seq? (:block parsed))))

(defn is-finally
  {:private true}
  [clause parsed]
  (and (is-call 'finally clause)
       (!seq? (:finally parsed))
       (!seq? (:catch parsed))
       (!seq? (:block parsed))))

(defn is-expr
  {:private true}
  [clause parsed]
  (!or (is-call 'catch clause)
       (is-call 'finally clause)))

(defn try-prepend
  {:private true}
  [parsed keyword clause]
  (conj parsed [keyword (cons clause (get parsed keyword))]))

(defn try-parse
  {:private true}
  [clauses]
  (unless (seq? clauses)
    {:block () :catch () :finally ()}
    (let [f (first clauses),
          r (rest clauses),
          p (try-parse r)]
      (cond
        (is-catch f p)   (try-prepend p :catch f)
        (is-finally f p) (try-prepend p :finally f)
        (is-expr f p)    (try-prepend p :block f)
        :else            (panic :message "malformed try-catch-finally")))))

(defn try-catch-branch
  {:private true}
  [clauses err-sym]
  (assert-args
    (seq? clauses) "catch branch not paired")
  (make-lazy-seq
    (let [clause (first clauses),
          var    ((clause 1) 1),
          expr   (rest (rest clause))]
      (cons
        (list 'sputter:let
              [var err-sym]
              [false (cons 'sputter:do expr)])
        (try-catch-clause (rest clauses) err-sym)))))

(defn try-catch-clause
  {:private true}
  [clauses err-sym]
  (make-lazy-seq
    (when (seq? clauses)
      (let [clause (first clauses),
            pred   ((clause 1) 0)]
        (cons
          (list pred err-sym)
          (try-catch-branch clauses err-sym))))))

(defn try-catch
  {:private true}
  [clauses]
  (let [err (gensym "err")]
    `(lambda [~err]
      (cond
        ~@(apply list (try-catch-clause clauses err))
        :else [true ~err]))))

(defn try-finally
  {:private true}
  [clauses]
  (map
    (lambda [clause] (first (rest clause)))
    clauses))

(defmacro try
  [& clauses]
  (assert-args
    (seq? clauses) "try-catch-finally requires at least one clause")
  (let [parsed#  (try-parse clauses),
        block#   (:block parsed#),
        catches# (:catch parsed#),
        finally# (:finally parsed#)]
    `(let [rec# (recover [false (do ~@block#)] ~(try-catch catches#)),
           err# (rec# 0),
           res# (rec# 1)]
      ~@(try-finally finally#)
      (if err# (raise res#) res#))))
