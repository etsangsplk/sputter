package native

import (
	"reflect"
	"regexp"
	"strings"

	a "github.com/kode4food/sputter/api"
)

// BadConversionType is raised if a type can't be converted
const BadConversionType = "can't convert between type: %s"

var camelCase = regexp.MustCompile("[a-z][A-Z]")

// New wraps a native wrapped using Go's reflection API
func New(i interface{}) a.Value {
	v := reflect.ValueOf(i)
	t := v.Type()
	c := getConvertOut(t)
	return c(v)
}

func typeName(t reflect.Type) a.Name {
	return a.Name(kebabCase(t.String()))
}

func kebabCase(n string) string {
	r := camelCase.ReplaceAllStringFunc(n, func(s string) string {
		return s[:1] + "-" + s[1:]
	})
	return strings.ToLower(r)
}
