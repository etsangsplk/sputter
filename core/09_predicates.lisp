;;;; sputter core: predicates

(defn even?
  {:doc "will return whether or not the number is even"}
  [value]
  (= (% value 2) 0))

(defn odd?
  {:doc "will return whether or not the number is odd"}
  [value]
  (!= (% value 2) 0))
