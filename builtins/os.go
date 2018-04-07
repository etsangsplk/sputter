package builtins

import (
	"os"
	"regexp"
	"time"

	a "github.com/kode4food/sputter/api"
)

const (
	envName         = "*env*"
	argsName        = "*args*"
	currentTimeName = "current-time"
)

type currentTimeFunction struct{ BaseBuiltIn }

var envPairRegex = regexp.MustCompile("^(?P<Key>[^=]+)=(?P<Value>.*)$")

func env() a.Value {
	var r []a.Vector
	for _, v := range os.Environ() {
		if e := envPairRegex.FindStringSubmatch(v); len(e) == 3 {
			r = append(r, a.Vector{
				a.NewKeyword(a.Name(e[1])),
				a.Str(e[2]),
			})
		}
	}
	return a.NewAssociative(r...)
}

func args() a.Value {
	r := a.Vector{}
	for _, v := range os.Args {
		r = append(r, a.Str(v))
	}
	return r
}

func (*currentTimeFunction) Apply(_ a.Context, args a.Vector) a.Value {
	return a.NewFloat(float64(time.Now().UnixNano()))
}

func init() {
	var currentTime *currentTimeFunction

	Namespace.Put(envName, env())
	Namespace.Put(argsName, args())

	RegisterBuiltIn(currentTimeName, currentTime)
}
