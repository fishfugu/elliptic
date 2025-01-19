package utils

import (
	"fmt"
	"math/big"
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
