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
  [clauses errSym]
  (make-lazy-seq
    (let [clause (first clauses)
          var    ((clause 1) 1)
          expr   (rest (rest clause))]
      (cons
        (list 'sputter:let
              [var errSym]
              [false (cons 'sputter:do expr)])
        (try-catch-clause (rest clauses) errSym)))))

(defn try-catch-clause
  {:private true}
  [clauses errSym]
  (make-lazy-seq
    (when (seq? clauses)
      (let [clause (first clauses)
            pred   ((clause 1) 0)]
        (cons
          (list pred errSym)
          (try-catch-branch clauses errSym))))))

(defn try-catch
  {:private true}
  [clauses]
  (let [err (gensym "err")]
  `(lambda [~err]
    (cond
      ~@(apply list (try-catch-clause clauses err))
      :else [true ~err]))))

(defmacro try
  [& clauses]
  (let [parsed# (try-parse clauses)]
    `(let [rec# (recover [false (do ~@(:block parsed#))]
                         ~(try-catch (:catch parsed#)))
           err# (rec# 0)
           res# (rec# 1)]
      (do ~@(rest (:finally parsed#)))
      (if err# (raise res#) res#))))
