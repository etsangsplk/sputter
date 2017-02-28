;;;; this is a comment at the top of the file

(defvar blah "blah is stored")

(defun say-hello [name]
  (println "Hello there, " name "!"))

(let (a 10)   ; this is an EOL comment
  (user:println "a = " a)
  (let (a 20)
    (println "a = " a))
  (println "a = " a))

(println
  "hello "
  (*
    (+ 9.5e+2 1.0)
    2.0)
  " "
  '(blah 99 "yep"))

(say-hello "Thom")

(println
  "howdy "
  "ho "
  blah)

(if (list? [1])
  (println "yep")
  (println "nope"))

(println [1 2 3])

(sputter:let
  [a (cons 1 (list 2 3 4))]
  (println a))

(println (car '(1 2 3)))
(println (cdr '(1 2 3)))

(println (list 'a 'b 'c))
(println (cons 50 (cons 99 100)))
(println (cons 100 200))

(defvar r (cons 3 (cons 4 (cons 9 10))))
(println (car r))
(println (cdr r))
(println (cdr (cdr r)))
(println (cdr (cdr (cdr r))))

(println (first r))
(println (second r))
(println (third r))

(defvar s (cons 4 (cons 9 10)))
(println (first s))
(println (second s))

(println -10)

(println (if (> 10 9 8 7 6) "hello" "not"))
(println (if (nil? false nil ()) "nil" "no"))

(defmacro foo [x] (println "hello"))
(println "macro: " (foo blah))
