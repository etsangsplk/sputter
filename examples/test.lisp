(defvar blah "blah is stored")

(defun say-hello [name]
  (print "Hello there, " name "!"))

(let (a 10)
  (print "a = " a)
  (let (a 20)
    (print "a = " a))
  (print "a = " a))

(print
  "hello "
  (*
    (+ 9.5e+2 1.0)
    2.0)
  " "
  '(blah 99 "yep"))

(say-hello "Thom")

(print
  "howdy "
  "ho "
  blah)

(if (list? [1])
  (print "yep")
  (print "nope"))

(print [1 2 3])
