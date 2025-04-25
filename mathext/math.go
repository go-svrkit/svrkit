// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mathext

import (
	"math"
	"math/bits"
)

func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Dim[T Number](x, y T) T {
	return max(x-y, 0)
}

func SafeDiv[T Number](x, y T) T {
	var zero T
	if y == zero {
		return zero
	}
	return x / y
}

func SafeMulUint64(a, b uint64) (product uint64, overflow bool) {
	hi, lo := bits.Mul64(a, b)
	if hi != 0 {
		return 0, true
	}
	return lo, false
}

func SafeMulInt64(a, b int64) (product int64, overflow bool) {
	var sign int64 = 1
	if a < 0 {
		a = -a
		sign = -sign
	}
	if b < 0 {
		b = -b
		sign = -sign
	}
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	if hi != 0 || lo > math.MaxInt64 {
		return 0, true
	}
	if sign < 0 {
		return -int64(lo), false
	}
	return int64(lo), false
}

// Truncate 截断浮点数的`n`位后，n不应过大
func Truncate(f float64, n int) float64 {
	var x = math.Pow10(n)
	return float64(int64(f*x)) / x
}

// RoundFloat round a float to a specific decimal place or precision
// see https://github.com/montanaflynn/stats/blob/master/round.go
func RoundFloat(x float64, places int) float64 {
	// If the float is not a number
	if math.IsNaN(x) {
		return math.NaN()
	}

	// Find out the actual sign and correct the input for later
	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	// Use the places arg to get the amount of precision wanted
	precision := math.Pow(10, float64(places))

	// Find the decimal place we are looking to round
	digit := x * precision

	// Get the actual decimal number as a fraction to be compared
	_, decimal := math.Modf(digit)

	// If the decimal is less than .5 we round down otherwise up
	var rounded float64
	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	// Finally we do the math to actually create a rounded number
	return rounded / precision * sign
}

// RoundHalf 四舍五入
func RoundHalf(v float64) int {
	return int(RoundFloat(v, 0))
}
