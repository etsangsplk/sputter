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

type (
	// Number can represent either floating point or rational numbers
	Number interface {
		Value
		Cmp(r Number) Comparison
		Add(r Number) Number
		Sub(r Number) Number
		Mul(r Number) Number
		Div(r Number) Number
		Mod(r Number) Number
		IsNaN() bool
		Float64() (float64, bool)
	}

	dec apd.Decimal
	rat big.Rat
)

var (
	ctx = apd.BaseContext.WithPrecision(53)

	// Zero is a convenience wrapper for the number 0
	Zero = NewFloat(0)

	// One is a convenience wrapper for the number 1
	One = NewFloat(1)

	// PosInfinity is a convenience wrapper for positive infinity
	PosInfinity = &dec{Form: apd.Infinite}

	// NegInfinity is a convenience wrapper for negative infinity
	NegInfinity = &dec{Form: apd.Infinite, Negative: true}
)

// NewFloat generates a new Number from a float64 value
func NewFloat(f float64) Number {
	if res, err := new(apd.Decimal).SetFloat64(f); err == nil {
		return (*dec)(res)
	}
	panic(ErrStr(ExpectedNumber, fmt.Sprintf("%f", f)))
}

// NewRatio generates a new Number from a ratio
func NewRatio(n int64, d int64) Number {
	return (*rat)(new(big.Rat).SetFrac64(n, d))
}

// ParseNumber attempts to parse a string into a Number Value
func ParseNumber(s Str) Number {
	ns := string(s)
	if f, _, err := ctx.SetString(new(apd.Decimal), ns); err == nil {
		return (*dec)(f)
	}
	if r, ok := new(big.Rat).SetString(ns); ok {
		return (*rat)(r)
	}
	panic(ErrStr(ExpectedNumber, s))
}

func (r *rat) toDecimal() *apd.Decimal {
	rr := (*big.Rat)(r)
	rf, _ := rr.Float64()
	if res, err := new(apd.Decimal).SetFloat64(rf); err == nil {
		return res
	}
	panic(ErrStr(ExpectedNumber, rr.String()))
}

func (d *dec) Cmp(n Number) Comparison {
	if dn, ok := n.(*dec); ok {
		return Comparison((*apd.Decimal)(d).Cmp((*apd.Decimal)(dn)))
	}
	rn, _ := n.(*rat)
	return Comparison((*apd.Decimal)(d).Cmp(rn.toDecimal()))
}

func (r *rat) Cmp(n Number) Comparison {
	if rn, ok := n.(*rat); ok {
		rr := (*big.Rat)(r)
		br := (*big.Rat)(rn)
		return Comparison(rr.Cmp(br))
	}
	dn := (*apd.Decimal)(n.(*dec))
	return Comparison(r.toDecimal().Cmp(dn))
}

func (d *dec) Add(n Number) Number {
	res := new(apd.Decimal)
	if dn, ok := n.(*dec); ok {
		ctx.Add(res, (*apd.Decimal)(d), (*apd.Decimal)(dn))
		return (*dec)(res)
	}
	rn, _ := n.(*rat)
	ctx.Add(res, (*apd.Decimal)(d), rn.toDecimal())
	return (*dec)(res)
}

func (r *rat) Add(n Number) Number {
	if rn, ok := n.(*rat); ok {
		return (*rat)(new(big.Rat).Add((*big.Rat)(r), (*big.Rat)(rn)))
	}
	res := new(apd.Decimal)
	ctx.Add(res, r.toDecimal(), (*apd.Decimal)(n.(*dec)))
	return (*dec)(res)
}

func (d *dec) Sub(n Number) Number {
	res := new(apd.Decimal)
	if dn, ok := n.(*dec); ok {
		ctx.Sub(res, (*apd.Decimal)(d), (*apd.Decimal)(dn))
		return (*dec)(res)
	}
	rn, _ := n.(*rat)
	ctx.Sub(res, (*apd.Decimal)(d), rn.toDecimal())
	return (*dec)(res)
}

func (r *rat) Sub(n Number) Number {
	if rn, ok := n.(*rat); ok {
		return (*rat)(new(big.Rat).Sub((*big.Rat)(r), (*big.Rat)(rn)))
	}
	res := new(apd.Decimal)
	ctx.Sub(res, r.toDecimal(), (*apd.Decimal)(n.(*dec)))
	return (*dec)(res)
}

func (d *dec) Mul(n Number) Number {
	res := new(apd.Decimal)
	if dn, ok := n.(*dec); ok {
		ctx.Mul(res, (*apd.Decimal)(d), (*apd.Decimal)(dn))
		return (*dec)(res)
	}
	rn, _ := n.(*rat)
	ctx.Mul(res, (*apd.Decimal)(d), rn.toDecimal())
	return (*dec)(res)
}

func (r *rat) Mul(n Number) Number {
	if rn, ok := n.(*rat); ok {
		return (*rat)(new(big.Rat).Mul((*big.Rat)(r), (*big.Rat)(rn)))
	}
	res := new(apd.Decimal)
	ctx.Mul(res, r.toDecimal(), (*apd.Decimal)(n.(*dec)))
	return (*dec)(res)
}

func (d *dec) Div(n Number) Number {
	res := new(apd.Decimal)
	if dn, ok := n.(*dec); ok {
		ctx.Quo(res, (*apd.Decimal)(d), (*apd.Decimal)(dn))
		return (*dec)(res)
	}
	rn, _ := n.(*rat)
	ctx.Quo(res, (*apd.Decimal)(d), rn.toDecimal())
	return (*dec)(res)
}

func (r *rat) Div(n Number) Number {
	if rn, ok := n.(*rat); ok {
		return (*rat)(new(big.Rat).Quo((*big.Rat)(r), (*big.Rat)(rn)))
	}
	res := new(apd.Decimal)
	ctx.Quo(res, r.toDecimal(), (*apd.Decimal)(n.(*dec)))
	return (*dec)(res)
}

func (d *dec) Mod(n Number) Number {
	res := new(apd.Decimal)
	if dn, ok := n.(*dec); ok {
		ctx.Rem(res, (*apd.Decimal)(d), (*apd.Decimal)(dn))
		return (*dec)(res)
	}
	rn, _ := n.(*rat)
	ctx.Rem(res, (*apd.Decimal)(d), rn.toDecimal())
	return (*dec)(res)
}

func (r *rat) Mod(n Number) Number {
	res := new(apd.Decimal)
	ctx.Rem(res, r.toDecimal(), (*apd.Decimal)(n.(*dec)))
	return (*dec)(res)
}

func (d *dec) IsNaN() bool {
	return d.Form == apd.NaN || d.Form == apd.NaNSignaling
}

func (r *rat) IsNaN() bool {
	return false
}

func (d *dec) Float64() (float64, bool) {
	v, err := (*apd.Decimal)(d).Float64()
	return v, err == nil
}

func (r *rat) Float64() (float64, bool) {
	return (*big.Rat)(r).Float64()
}

func (d *dec) Str() Str {
	return Str((*apd.Decimal)(d).ToStandard())
}

func (r *rat) Str() Str {
	return Str((*big.Rat)(r).String())
}

// AssertNumber will cast a Value into a Number or explode violently
func AssertNumber(v Value) Number {
	if r, ok := v.(Number); ok {
		return r
	}
	panic(ErrStr(ExpectedNumber, v))
}

// AssertInteger will cast a Value into an Integer or explode violently
func AssertInteger(v Value) int {
	n := AssertNumber(v)
	f, _ := n.Float64()
	i := int(f)
	if f-float64(i) == 0 {
		return i
	}
	panic(ErrStr(ExpectedInteger, n))
}
