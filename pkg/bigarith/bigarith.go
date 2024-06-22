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
	// From ---^ that site: "The probability of returning true for a randomly chosen non-prime is at most ¼ⁿ"
	// i.e. (1/4)^n - where n is the parameter handed in to the function
	// TODO: work out if there's a way to decide what that param should be set to
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

// Cmp compares two integers represented as strings and returns:
// -1 if a < b, 0 if a == b, 1 if a > b
func Cmp(a, b string) (int, error) {
	x, okX := new(big.Int).SetString(a, 10)
	y, okY := new(big.Int).SetString(b, 10)
	if !okX || !okY {
		return 0, fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}
	return x.Cmp(y), nil
}

// Exp computes the exponentiation of base a to exp e modulo m.
func Exp(b, e, m string) (string, error) {
	base, okB := new(big.Int).SetString(b, 10)
	exp, okE := new(big.Int).SetString(e, 10)
	mod, okM := new(big.Int).SetString(m, 10)
	if !okB || !okE || !okM {
		return "", fmt.Errorf("invalid input: b = %s, e = %s, m = %s - cannot create all the integers required, from this input", b, e, m)
	}
	result := new(big.Int).Exp(base, exp, mod)
	return result.String(), nil
}

// Mod performs the modulus operation of a by m.
func Mod(a, m string) (string, error) {
	x, okX := new(big.Int).SetString(a, 10)
	mod, okM := new(big.Int).SetString(m, 10)
	if !okX || !okM {
		return "", fmt.Errorf("invalid input: a = %s, m = %s - cannot create all the integers required, from this input", a, m)
	}
	result := new(big.Int).Mod(x, mod)
	return result.String(), nil
}

// FindPrime finds a prime number within the range specified by the strings low and high.
// Returns the first prime found as a string, or an error if no prime is found or inputs are invalid.
func FindPrime(low, high string) (string, error) {
	// TODO: is there any reason to implement different way / direction of search?
	// NOTE: currently from low bound, up
	lowInt, ok := new(big.Int).SetString(low, 10)
	if !ok {
		return "", fmt.Errorf("invalid input for lower bound")
	}
	highInt, ok := new(big.Int).SetString(high, 10)
	if !ok {
		return "", fmt.Errorf("invalid input for upper bound")
	}

	// Start searching for a prime at the low end of the range
	for p := lowInt; p.Cmp(highInt) <= 0; p.Add(p, big.NewInt(1)) {
		if p.ProbablyPrime(20) { // 20 iterations of Miller-Rabin, quite strong
			return p.String(), nil
		}
	}

	return "", fmt.Errorf("no prime found in the specified range")
}
