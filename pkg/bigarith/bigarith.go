package bigarith

import (
	"fmt"
	"math/big"
)

// Add takes two string representations of integers, adds them, and returns the result as a string.
// Returns an error if the input strings are not valid integers.
func Add(a, b string) (string, error) {
	x := new(big.Int)
	y := new(big.Int)
	_, okX := x.SetString(a, 10)
	_, okY := y.SetString(b, 10)
	if !okX || !okY {
		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}

	result := new(big.Int).Add(x, y)
	return result.String(), nil
}

// Subtract takes two string representations of integers, subtracts the second from the first,
// and returns the result as a string.
// Returns an error if the input strings are not valid integers.
func Subtract(a, b string) (string, error) {
	x := new(big.Int)
	y := new(big.Int)
	_, okX := x.SetString(a, 10)
	_, okY := y.SetString(b, 10)
	if !okX || !okY {
		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}

	result := new(big.Int).Sub(x, y)
	return result.String(), nil
}

// Multiply takes two string representations of integers, multiplies them,
// and returns the result as a string.
// Returns an error if the input strings are not valid integers.
func Multiply(a, b string) (string, error) {
	x := new(big.Int)
	y := new(big.Int)
	_, okX := x.SetString(a, 10)
	_, okY := y.SetString(b, 10)
	if !okX || !okY {
		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}

	result := new(big.Int).Mul(x, y)
	return result.String(), nil
}

// DivideInField takes two string representations of integers, divides a by b,
// and returns the result as a string in a finite field defined by a prime modulus p.
// Returns an error if the input strings are not valid integers, or if b has no inverse modulo p,
// (which includes if b is 0).
func DivideInField(a, b, p string) (string, error) {
	x := new(big.Int)
	y := new(big.Int)
	prime := new(big.Int)
	_, okX := x.SetString(a, 10)
	_, okY := y.SetString(b, 10)
	_, okP := prime.SetString(p, 10)
	if !okX || !okY || !okP {
		return "", fmt.Errorf("invalid input: a = %s, b = %s, p = %s - cannot create all the integers required, from this input", a, b, p)
	}

	// First, find the multiplicative inverse of b mod p
	bInv, err := ModularInverse(b, p)
	if err != nil {
		return "", fmt.Errorf("error finding inverse: %v", err)
	}

	// Calculate a * bInv mod p
	inverse := new(big.Int)
	inverse.SetString(bInv, 10)
	result := new(big.Int).Mul(x, inverse) // a * b^-1
	result.Mod(result, prime)              // mod p

	return result.String(), nil
}

// ModularInverse calculates the multiplicative inverse of a modulo p using Fermat's Little Theorem
// and returns it as a string. Returns an error if p is not prime or a and p are not coprime.
func ModularInverse(a, p string) (string, error) {
	aInt := new(big.Int)
	pInt := new(big.Int)
	_, okA := aInt.SetString(a, 10)
	_, okP := pInt.SetString(p, 10)
	zero := big.NewInt(0)
	if !okA || !okP {
		return "", fmt.Errorf("invalid input: a = %s, p = %s - cannot create all the integers required, from this input", a, p)
	}
	if aInt.Cmp(zero) == 0 {
		return "", fmt.Errorf("invalid input: a is ZERO - no modular multiplicative inverse")
	}

	// Check if p is prime; if not, the multiplicative inverse might not exist
	// https://pkg.go.dev/math/big#Int.ProbablyPrime
	// The probability of returning true for a randomly chosen non-prime is at most ¼ⁿ.
	// i.e. (1/4)^n - where n is the parameter handed in to the function
	if !pInt.ProbablyPrime(100) {
		return "", fmt.Errorf("invalid input: modulus %s is not prime", p)
	}

	// Calculate a^(p-2) mod p
	pMinusTwo := new(big.Int).Sub(pInt, big.NewInt(2)) // p-2
	inverse := new(big.Int).Exp(aInt, pMinusTwo, pInt)
	return inverse.String(), nil
}

// Divide takes two string representations of integers
// Returns an error explaining why it isn't implemented.
func Divide(a, b string) (string, error) {
	return "", fmt.Errorf("division not implemented - due to ambiguous integer results")

}
