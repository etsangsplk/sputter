(defn is-call
  {:private true}
  [sym clause]
  (and
    (local? sym)
    (list? clause)
    (eq sym (first clause))))

(defn is-catch
  {:private true}
  [clause parsed]
  (and
    (sputter:is-call 'catch clause)
    (!seq? (:block parsed))))

(defn is-finally
  {:private true}
  [clause parsed]
  (and
    (sputter:is-call 'finally clause)
    (!seq? (:catch parsed)
    (!seq? (:block parsed)))))

(defn is-valid-expr
  {:private true}
  [clause parsed]
  (!or (sputter:is-call 'catch clause)
       (sputter:is-call 'finally clause)))

(defn try-prepend
  {:private true}
  [parsed keyword clause]
  (conj parsed [keyword (cons clause (get parsed keyword))]))

(defn try-parse
  {:private true}
  [clauses]
  (unless (seq? clauses)
    {:block () :catch () :finally ()}
    (let [f (first clauses) r (rest clauses) p (try-parse r)]
      (cond
        (is-catch f p)      (try-prepend p :catch f)
        (is-finally f p)    (try-prepend p :finally f)
        (is-valid-expr f p) (try-prepend p :block f)
        :else               (panic :message "malformed try macro")))))

(defmacro try
  [& clauses]
  (try-recover (try-parse ~clauses)))
