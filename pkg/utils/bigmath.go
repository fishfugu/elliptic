package utils

import "math/big"

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

func ratIsWholeNumber(r *big.Rat) bool {
	// The denominator must be 1 for it to be a whole number
	return r.Denom().Cmp(big.NewInt(1)) == 0
}
