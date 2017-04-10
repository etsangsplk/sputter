package api

import (
	"fmt"
	"math/big"

	"github.com/cockroachdb/apd"
)

const (
	// ExpectedNumber is thrown when a Value is not a Number
	ExpectedNumber = "value is not a number: %s"

	// ExpectedInteger is thrown when a Value is not an Integer
	ExpectedInteger = "value is not an integer: %s"
)

// Number can represent either floating point or rational numbers
type Number struct {
	decimal *apd.Decimal
	ratio   *big.Rat
}

var ctx = apd.BaseContext.WithPrecision(53)

// NewFloat generates a new Number from a float64 value
func NewFloat(f float64) *Number {
	if res, err := nativeDecimal().SetFloat64(f); err == nil {
		return &Number{decimal: res}
	}
	panic(Err(ExpectedNumber, fmt.Sprintf("%f", f)))
}

// NewRatio generates a new Number from a ratio
func NewRatio(n int64, d int64) *Number {
	res := new(big.Rat).SetFrac64(n, d)
	return &Number{ratio: res}
}

func nativeDecimal() *apd.Decimal {
	return &apd.Decimal{}
}

func nativeRatio() *big.Rat {
	return big.NewRat(1, 1)
}

// ParseNumber attempts to parse a string into a Number Value
func ParseNumber(s string) *Number {
	if f, _, err := ctx.SetString(nativeDecimal(), s); err == nil {
		return &Number{decimal: f}
	}
	if r, ok := nativeRatio().SetString(s); ok {
		return &Number{ratio: r}
	}
	panic(Err(ExpectedNumber, String(s)))
}

func (l *Number) toDecimal() *apd.Decimal {
	if lf := l.decimal; lf != nil {
		return lf
	}

	rf, _ := l.ratio.Float64()
	if res, err := nativeDecimal().SetFloat64(rf); err == nil {
		return res
	}
	panic(Err(ExpectedNumber, l.ratio.String()))
}

// Cmp compares this Number Value to another Value
func (l *Number) Cmp(r *Number) Comparison {
	if l.decimal != nil || r.decimal != nil {
		return Comparison(l.toDecimal().Cmp(r.toDecimal()))
	}
	return Comparison(l.ratio.Cmp(r.ratio))
}

// Add will add two Numbers
func (l *Number) Add(r *Number) *Number {
	if l.decimal != nil || r.decimal != nil {
		res := nativeDecimal()
		ctx.Add(res, l.toDecimal(), r.toDecimal())
		return &Number{decimal: res}
	}
	return &Number{ratio: nativeRatio().Add(l.ratio, r.ratio)}
}

// Sub will subtract two Numbers
func (l *Number) Sub(r *Number) *Number {
	if l.decimal != nil || r.decimal != nil {
		res := nativeDecimal()
		ctx.Sub(res, l.toDecimal(), r.toDecimal())
		return &Number{decimal: res}
	}
	return &Number{ratio: nativeRatio().Sub(l.ratio, r.ratio)}
}

// Mul will multiply two Numbers
func (l *Number) Mul(r *Number) *Number {
	if l.decimal != nil || r.decimal != nil {
		res := nativeDecimal()
		ctx.Mul(res, l.toDecimal(), r.toDecimal())
		return &Number{decimal: res}
	}
	return &Number{ratio: nativeRatio().Mul(l.ratio, r.ratio)}
}

// Div will divide two Numbers
func (l *Number) Div(r *Number) *Number {
	if l.decimal != nil || r.decimal != nil {
		res := nativeDecimal()
		ctx.Quo(res, l.toDecimal(), r.toDecimal())
		return &Number{decimal: res}
	}
	return &Number{ratio: nativeRatio().Quo(l.ratio, r.ratio)}
}

// Float64 converts the value to a native float64
func (l *Number) Float64() (float64, bool) {
	if nf := l.decimal; nf != nil {
		v, err := nf.Float64()
		return v, err == nil
	}
	return l.ratio.Float64()
}

func (l *Number) String() string {
	if nf := l.decimal; nf != nil {
		return nf.ToStandard()
	}
	return l.ratio.String()
}

// AssertNumber will cast a Value into a Number or explode violently
func AssertNumber(v Value) *Number {
	if r, ok := v.(*Number); ok {
		return r
	}
	panic(Err(ExpectedNumber, String(v)))
}

// AssertInteger will cast a Value into an Integer or explode violently
func AssertInteger(v Value) int {
	n := AssertNumber(v)
	f, _ := n.Float64()
	i := int(f)
	if f-float64(i) == 0 {
		return i
	}
	panic(Err(ExpectedInteger, n.String()))
}
