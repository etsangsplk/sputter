package builtins_test

import (
	"testing"

	a "github.com/kode4food/sputter/api"
)

func TestAssoc(t *testing.T) {
	testCode(t, `(len {:name "Sputter", :age 45})`, f(2))
	testCode(t, `(len (assoc :name "Sputter", :age 45))`, f(2))
	testCode(t, `(assoc? {:name "Sputter" :age 45})`, a.True)
	testCode(t, `(assoc? (assoc :name "Sputter" :age 45))`, a.True)
	testCode(t, `(assoc? '(:name "Sputter" :age 45))`, a.False)
	testCode(t, `(assoc? [:name "Sputter" :age 45])`, a.False)
	testCode(t, `(!assoc? '(:name "Sputter" :age 45))`, a.True)
	testCode(t, `(!assoc? [:name "Sputter" :age 45])`, a.True)
	testCode(t, `(:name {:name "Sputter" :age 45})`, s("Sputter"))

	testCode(t, `
		(:name (apply assoc (append '(:name "Sputter") '(:age 45))))
	`, s("Sputter"))

	a.GetNamespace(a.UserDomain).Delete("x")
	testCode(t, `
		(def x {:name "bob" :age 45})
		(x :name)
	`, s("bob"))

	testBadCode(t, `(assoc :too "few" :args)`, a.ExpectedPair)

	testBadCode(t, `
		(apply assoc (append '(:name "Sputter") '(:age)))
	`, a.ExpectedPair)
}

func TestMapped(t *testing.T) {
	testCode(t, `(mapped? {:name "Sputter" :age 45})`, a.True)
	testCode(t, `(mapped? (assoc :name "Sputter" :age 45))`, a.True)
	testCode(t, `(mapped? '(:name "Sputter" :age 45))`, a.False)
	testCode(t, `(mapped? [:name "Sputter" :age 45])`, a.False)
	testCode(t, `(!mapped? '(:name "Sputter" :age 45))`, a.True)
	testCode(t, `(!mapped? '(:name "Sputter" :age 45))`, a.True)
	testCode(t, `(!mapped? [:name "Sputter" :age 45])`, a.True)
}

func TestKeywordApply(t *testing.T) {
	ns := a.GetNamespace(a.UserDomain)
	ns.Delete("test")

	testCode(t, `
		(def test {
			:double (lambda [x] (* x 2))
		})
		(:double test 16)
	`, f(32))
}
