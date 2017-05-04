package api

import "fmt"

// Err generates a standard interpreter error
func Err(s string, args ...interface{}) string {
	sargs := make([]interface{}, len(args))
	for i, a := range args {
		if v, ok := a.(Value); ok {
			sargs[i] = string(v.Str())
		} else {
			sargs[i] = a
		}
	}
	return fmt.Sprintf(s, sargs...)
}
