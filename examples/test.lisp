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

(say-hello "Thom")

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

(prn (car '(1 2 3)))
(prn (cdr '(1 2 3)))

(prn (list 'a 'b 'c))
(prn (cons 50 (cons 99 100)))
(prn (cons 100 200))

(def r (cons 3 (cons 4 (cons 9 10))))
(prn (car r))
(prn (cdr r))
(prn (cdr (cdr r)))
(prn (cdr (cdr (cdr r))))

(prn (first r))
(prn (rest r))

(def s (cons 4 (cons 9 10)))
(prn (first s))
(prn (rest s))

(prn -10)

(prn (if (> 10 9 8 7 6) "hello" "not"))
(prn (if (nil? false nil ()) "nil" "no"))
(prn (if (!nil? false nil ()) "nil" "no"))

(sputter:defmacro foo [x] (prn "hello"))
(prn "macro:" (foo blah))

(with-ns foo
  (defn bar [x y] (+ x y)))

(prn (ns foo))
