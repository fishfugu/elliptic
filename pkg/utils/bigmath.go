package utils

import (
	"fmt"
	"math/big"
)

// 2048 bit precision was chosen as approximately 4 times the MAXIMUM number of bits
// used for. keys in EC Cryptography in the highest security situations:
// https://en.wikipedia.org/wiki/Key_size?utm_source=chatgpt.com
// "521-bit keys: Deliver a security level of roughly 256 bits, used in scenarios requiring
// the highest security assurances."

var (
	OneInt, TwoInt, FourInt            = big.NewInt(1), big.NewInt(2), big.NewInt(4)
	ZeroRat, TwoRat, ThreeRat, FourRat = big.NewRat(0, 1), big.NewRat(2, 1), big.NewRat(3, 1), big.NewRat(4, 1)

	// precision_2048           uint     = 2048
	// tolerance_1024           uint64   = 1024
	precision_128 uint   = 128
	tolerance_256 uint64 = 256

	// toleranceInt                      = big.NewInt(int64(tolerance_1024))
	toleranceInt                      = big.NewInt(int64(tolerance_256))
	twoToThePowerOfTolerance          = new(big.Int).Exp(TwoInt, toleranceInt, nil)
	ToleranceFractionRat     *big.Rat = new(big.Rat).SetFrac(OneInt, twoToThePowerOfTolerance)

	HalfFloat *big.Float = NewFloat().SetFloat64(0.5)
)

// FindPrime finds a prime number within the range specified by the strings i and high.
// Returns the first prime found as a string, or an error if no prime is found or inputs are invalid.
func FindPrime(i, high *big.Int) *big.Int {
	// TODO: is there any reason to implement different way / direction of search?
	// Start searching for a prime at the low end of the range
	for p := i; new(big.Int).Set(p).Cmp(high) < 0; p = new(big.Int).Add(p, new(big.Int).SetInt64(1)) {
		if p.ProbablyPrime(1000) {
			return p
		}
	}
	return new(big.Int).SetInt64(0)
}

// newFloat creates a new big.Float with the default precision
// ue this whenever creating new Float (except in very specific circumstances...)
func NewFloat() *big.Float {
	return new(big.Float).SetPrec(0)
}

// SqrtRat computes the square root of a big.Rat with arbitrary precision.
// If the result is not an exact rational number, it computes an approximation with the specified precision.
func SqrtRat(input *big.Rat) (*big.Rat, error) {
	num, den := input.Num(), input.Denom()

	// Check if numerator and denominator are perfect squares
	if sqrtNum, sqrtDen := IntSqrt(num), IntSqrt(den); sqrtNum != nil && sqrtDen != nil {
		return new(big.Rat).SetFrac(sqrtNum, sqrtDen), nil
	}

	// Approximation for non-perfect square roots
	floatInput := new(big.Float).SetPrec(precision_128).SetRat(input)
	floatSqrt := SqrtFloat(floatInput, precision_128)

	result := new(big.Rat)
	floatSqrt.Rat(result)
	return result, nil
}

// SqrtFloat computes the square root of a big.Float using Newton's method with the specified precision.
func SqrtFloat(a *big.Float, prec uint) *big.Float {
	logger := InitialiseLogger("[sqrtFloat]")
	logger.Debug("starting function sqrtFloat")

	// Initial guess: x0 = a / 2
	guess := NewFloat().Quo(a, big.NewFloat(2))

	// Iteratively refine the guess
	for i := uint(0); i < prec; i++ {
		temp := NewFloat().Quo(a, guess)     // temp = a / guess
		temp2 := NewFloat().Add(guess, temp) // temp2 = (guess + a/guess)
		guess = NewFloat().Mul(temp2, HalfFloat)
	}

	return guess
}

// IntSqrt computes the integer square root of a big.Int if it is a perfect square.
// Otherwise, it returns nil.
func IntSqrt(x *big.Int) *big.Int {
	logger := InitialiseLogger("[intSqrt]")
	logger.Debug("starting function intSqrt")

	// Use binary search to find the integer square root
	low := big.NewInt(0)
	high := new(big.Int).Set(x)
	mid := new(big.Int)
	square := new(big.Int)

	for low.Cmp(high) <= 0 {
		logger.Debug(fmt.Sprintf("Low: %s, High: %s", low.String(), high.String()))
		mid = new(big.Int).Rsh(new(big.Int).Add(low, high), 1) // mid = (low + high) / 2
		logger.Debug(fmt.Sprintf("Mid: %s:", mid.String()))
		square = new(big.Int).Mul(mid, mid) // square = mid^2

		cmp := square.Cmp(x)
		logger.Debug(fmt.Sprintf("X: %s, Square: %s", x.String(), square.String()))

		if cmp == 0 {
			return mid // Perfect square
		} else if cmp < 0 {
			low.Add(mid, big.NewInt(1))
		} else {
			high.Sub(mid, big.NewInt(1))
		}
	}

	return nil // Not a perfect square
}

// Borwein's AGM-Based Method (Best for High Precision)
// The Arithmetic-Geometric Mean (AGM) method is the fastest way to compute high-precision trigonometric functions.
// Given:
// \[ a_0 = 1, \quad b_0 = \frac{1}{\sqrt{1+x^2}} \]
// Iterate:
// \[ a_{n+1} = \frac{a_n + b_n}{2}, \quad b_{n+1} = \sqrt{a_n b_n} \]
// until \( |a_n - b_n| \) is sufficiently small.
// Then:
// \[
// \arctan(x) = \frac{x}{a_n}
// \]
func ArctanAGM(x *big.Rat) (*big.Rat, error) {
	a := big.NewRat(1, 1)
	sqrtOnePlusXSquared, err := SqrtRat(new(big.Rat).Add(new(big.Rat).SetInt64(1), new(big.Rat).Mul(x, x)))
	if err != nil {
		return nil, err
	}
	b := new(big.Rat).Inv(sqrtOnePlusXSquared)

	for {
		aNext := new(big.Rat).Add(a, b)
		aNext.Quo(aNext, big.NewRat(2, 1))

		bNext := new(big.Rat).Mul(a, b)
		bNext, err = SqrtRat(bNext)
		if err != nil {
			return nil, err
		}

		if new(big.Rat).Abs(new(big.Rat).Sub(a, b)).Cmp(new(big.Rat).SetFrac64(1, 1000000)) < 0 {
			break
		}

		a.Set(aNext)
		b.Set(bNext)
	}
	return new(big.Rat).Quo(x, a), nil
}
