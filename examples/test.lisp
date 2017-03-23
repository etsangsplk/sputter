;;;; this is a comment at the top of the file
(def blah "blah is stored")

(defn say-hello [name]
  (println "Hello there," name "!"))

(let (a 10)   ; this is an EOL comment
  (println "a =" a)
  (let (a 20)
    (println "a =" a))
  (println "a =" a))

(prn
  "hello"
  (*
    (+ 9.5e+2 1.0)
    2.0)
  '(blah 99 "yep"))

(say-hello "Sputter")

(prn
  "howdy"
  "ho"
  blah)

(if (list? [1])
  (prn "yep")
  (prn "nope"))

(prn [1 2 3])

(sputter:let
  [a (cons 1 (list 2 3 4))]
  (prn a))

(prn (first '(1 2 3)))
(prn (rest '(1 2 3)))

(prn (list 'a 'b 'c))
(prn (cons 50 (list 99 100)))
(prn (list 100 200))

(def r (cons 3 (cons 4 (list 9 10))))
(prn (first r))
(prn (rest r))
(prn (rest (rest r)))
(prn (rest (rest (rest r))))

(prn (first r))
(prn (rest r))

(def s (cons 4 (list 9 10)))
(prn (first s))
(prn (rest s))

(prn -10)

(prn (if (> 10 9 8 7 6) "hello" "not"))
(prn (if (nil? false nil ()) "nil" "no"))
(prn (if (!nil? false nil ()) "nil" "no"))

(sputter:defmacro foo [x] (prn "hello macro"))
(prn "macro:" (foo blah))

(with-ns foo
  (defn bar [x y] (+ x y)))

(prn (ns foo))
