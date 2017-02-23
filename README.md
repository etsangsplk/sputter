# Sputter (a Lisp Interpreter)
Sputter is a simple Lisp Interpreter written in
[Go](https://golang.org/). Basically, it's just me having some fun
and trying to improve my Go skills. That means you're unlikely to
find something you'd want to use in production here.  On the other
hand, if you want to join in on the fun, you're more than welcome
to.

## How To Install
Make sure your `GOPATH` is set, then run `go get` to retrieve the
package.

```bash
go get github.com/kode4food/sputter
```

## How To Invoke The Interpreter
Once you've installed the package, you can run it from `GOPATH/bin`
like so.

```bash
sputter ./examples/test.lisp
```

## Current Status
I just started this thing and it's still pretty fragile, but
that will change rapidly.  The current built-in forms are:

  * Control: `quote`, `do`
  * Numeric: `+`, `-`, `*`, `/`, `=`, `<`, `<=`, `>`, `>=`
  * Variables and Functions: `defvar`, `let`, `defun`, `lambda`
  * Branching: `if`, `cond`
  * Predicates: `eq`, `nil?`, `list?`
  * List Stuff: `list`, `cons`, `car`, `cdr`, `first`, `second`, `third`
  * I/O: `print`, `println`

## License (MIT License)
Copyright (c) 2017 Thomas S. Bradford

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
