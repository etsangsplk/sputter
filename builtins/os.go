package builtins

import (
	"os"
	"regexp"

	a "github.com/kode4food/sputter/api"
)

var envPairRegex = regexp.MustCompile("^(?P<Key>[^=]+)=(?P<Value>.*)$")

func env() a.Value {
	r := []a.Vector{}
	for _, v := range os.Environ() {
		e := envPairRegex.FindStringSubmatch(v)
		r = append(r, a.NewVector(
			a.NewKeyword(a.Name(e[1])),
			a.Str(e[2])))
	}
	return a.NewAssociative(r...)
}

func args() a.Value {
	r := a.Values{}
	for _, v := range os.Args {
		r = append(r, a.Str(v))
	}
	return a.NewVector(r...)
}

func init() {
	Namespace.Put("*env*", env())
	Namespace.Put("*args*", args())
}
