package sputter

import "math/big"
import "fmt"

// Builtins are the Context of built-in identifiers
var Builtins *Context

func addition(c *Context, l *List) Value {
	sum := big.NewFloat(0)
	for current := l; current != EmptyList; current = current.rest {
		sum.Add(sum, Evaluate(c, current.value).(*big.Float))
	}
	return sum
}

func subtraction(c *Context, l *List) Value {
	diff := Evaluate(c, l.value).(*big.Float)
	for current := l.rest; current != EmptyList; current = current.rest {
		diff.Sub(diff, Evaluate(c, current.value).(*big.Float))
	}
	return diff
}

func multiplication(c *Context, l *List) Value {
	prod := big.NewFloat(1)
	for current := l; current != EmptyList; current = current.rest {
		prod.Mul(prod, Evaluate(c, current.value).(*big.Float))
	}
	return prod
}

func division(c *Context, l *List) Value {
	quotient := Evaluate(c, l.value).(*big.Float)
	for current := l.rest; current != EmptyList; current = current.rest {
		quotient.Quo(quotient, Evaluate(c, current.value).(*big.Float))
	}
	return quotient
}

func puts(c *Context, l *List) Value {
	for current := l; current != EmptyList; current = current.rest {
		fmt.Println(EvaluateToString(c, current.value))
	}
	return nil
}

func init() {
	Builtins = NewContext()
	Builtins.Put("+", &Function{addition})
	Builtins.Put("-", &Function{subtraction})
	Builtins.Put("*", &Function{multiplication})
	Builtins.Put("/", &Function{division})
	Builtins.Put("puts", &Function{puts})
}
