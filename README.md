# Sputter (a Lisp Experiment)
[![Go Report Card](https://goreportcard.com/badge/github.com/kode4food/sputter)](https://goreportcard.com/report/github.com/kode4food/sputter) [![Build Status](https://travis-ci.org/kode4food/sputter.svg?branch=master)](https://travis-ci.org/kode4food/sputter) [![Coverage Status](https://coveralls.io/repos/github/kode4food/sputter/badge.svg?branch=master)](https://coveralls.io/github/kode4food/sputter?branch=master)

Sputter is a Lisp Environment written in [Go](https://golang.org/).
Basically, it's just me having some fun and trying to improve my Go skills.
That means you're unlikely to find something you'd want to use in production
here. On the other hand, if you want to join in on the fun, you're more than
welcome to.

## How To Install
Make sure your `GOPATH` is set, then run `go get` to retrieve the package.

```bash
go get github.com/kode4food/sputter
```

## How To Run A Source File
Once you've installed the package, you can run it from `GOPATH/bin` like so:

```bash
sputter somefile.lisp

# or

cat somefile.lisp | sputter
```

## How To Start The REPL
Sputter has a very crude Read-Eval-Print Loop that will be more than happy
to start if you call it with no arguments from the terminal:

<img src="docs/img/repl.jpeg" />

## Current Status
I just started this thing and it's still pretty fragile, but that will
change rapidly. The current built-in forms are:

  * Basics: `read` `eval` `do`
  * Branching: `if` `not` `unless` `when` `when-not` `cond` `and` `or`
  * Numeric: `+` `-` `*` `/` `%` `=` `!=` `<` `>` `<=` `>=` `inf` `-inf` `inc` `dec` `range`
  * Numeric Predicates: `inf?` `-inf?` `nan?` `even?` `odd?`
  * Values: `def` `let` `when-let` `put-ns` `ns` `with-ns`
  * Symbols and Keywords: `sym` `gensym` `sym?` `local?` `keyword?`
  * Functions: `defn` `fn` `lambda` `apply` `partial` `apply?` `closure` `special-form?`
  * Threading: `->` `->>` `some->` `some->>` `as->` `cond->` `cond->>`
  * Macros: `defmacro` `macroexpand1` `macroexpand` `macroexpand-all` `macro?`
  * Errors: `error` `raise` `panic` `try` `assert-args`
  * Quoting: `quote` `syntax-quote`
  * Predicates: `eq` `nil?`
  * Sequences: `cons` `conj` `first` `rest` `last` `for-each` `seq?`
  * Lists: `list` `to-list` `list?`
  * Vectors: `vector` `to-vector` `vector?`
  * Associative Arrays: `assoc` `to-assoc` `assoc?`
  * Counted Sequences: `len` `len?`
  * Indexed Sequences: `nth` `indexed?`
  * Mapped Sequences: `get` `.` `mapped?`
  * Comprehensions: `concat` `map` `filter` `reduce` `take` `take-while` `drop`
  * Lazy Sequences: `lazy-seq` `partition`
  * Metadata: `meta` `with-meta` `meta?`
  * Concurrency: `go` `chan` `generate` `future` `promise` `promise?`
  * Strings: `str` `str!` `str?`
  * I/O: `print` `println` `pr` `prn` `with-open`
  * Operating System: `time` `*env*` `*args*` `*stdout*` `*stderr*` `*stdin*`

Documentation for most of these forms may be viewed in the REPL using the
`doc` function.
