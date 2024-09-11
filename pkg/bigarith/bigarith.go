package bigarith

import (
	"fmt"
	"math/big"
	"strings"
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
	return new(big.Int).Add(x, y).String(), nil
}

// AddFloat takes two string representations of floats, adds them, and returns the result as a string.
// Returns an error if the input strings are not valid integers.
func AddFloat(a, b string) (string, error) {
	x, okX := new(big.Float).SetString(a)
	y, okY := new(big.Float).SetString(b)
	if !okX || !okY {
		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}
	return new(big.Float).Add(x, y).String(), nil
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
	return new(big.Int).Sub(x, y).String(), nil
}

// SubtractFloat takes two string representations of floats, subtracts the second from the first,
// and returns the result as a string.
// Returns an error if the input strings are not valid integers.
func SubtractFloat(a, b string) (string, error) {
	x, okX := new(big.Float).SetString(a)
	y, okY := new(big.Float).SetString(b)
	if !okX || !okY {
		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}
	return new(big.Float).Sub(x, y).String(), nil
}

// Multiply takes two string representations of integers, multiplies them,
// and returns the result as a string.
// Returns an error if the input strings are not valid integers.
// big.Int Mul sets z to the product x*y and returns z.
func Multiply(a, b string) (string, error) {
	x := new(big.Int)
	y := new(big.Int)
	_, okX := x.SetString(a, 10)
	_, okY := y.SetString(b, 10)
	if !okX || !okY {
		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}
	return new(big.Int).Mul(x, y).String(), nil
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
	return new(big.Int).Exp(aInt, pMinusTwo, pInt).String(), nil
}

// Divide takes two string representations of integers
// Returns an error explaining why it isn't implemented.
// See `DivideInField` for division in modular case.
func Divide(a, b string) (string, error) {
	x, okX := new(big.Float).SetString(a)
	y, okY := new(big.Float).SetString(b)
	if !okX || !okY {
		return "0", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}
	return new(big.Float).Quo(x, y).String(), nil
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

// CmpFloat compares two floats represented as strings and returns:
// -1 if a < b, 0 if a == b, 1 if a > b
func CmpFloat(a, b string) (int, error) {
	x, okX := new(big.Float).SetString(a)
	y, okY := new(big.Float).SetString(b)
	if !okX || !okY {
		return 0, fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
	}
	return x.Cmp(y), nil
}

// Exp computes the exponentiation of base b to exp e modulo m - represented as strings.
// If m == "", returns b**e unless y <= 0 then returns "1".
// big.Int Exp sets z = x**y mod |m| (i.e. the sign of m is ignored), and returns z.
// If m == nil or m == 0, z = x**y unless y <= 0 then z = 1.
// If m != 0, y < 0, and x and m are not relatively prime, z is unchanged and nil is returned.
func Exp(b, e, m string) (string, error) {
	// Check if b is already an integer, and if not, scale it
	decimalIndex := strings.Index(b, ".")
	scaleFactor := 0

	if decimalIndex != -1 {
		if m != "" {
			return "", fmt.Errorf("cannot have non-integer base AND a modulus - either base must be integer OR there must be no modulus value: b = %s, e = %s, m = %s - cannot create all the integers required, from this input", b, e, m)
		}
		// b IS NOT an integer
		scaleFactor = len(b) - decimalIndex - 1
		b = strings.Replace(b, ".", "", 1) // Remove the decimal point
	}

	// Parse base and exponent as big.Int
	base, okB := new(big.Int).SetString(b, 10)
	exp, okE := new(big.Int).SetString(e, 10)

	var mod *big.Int
	var okM bool
	if m == "" {
		mod = nil
		okM = true
	} else {
		mod, okM = new(big.Int).SetString(m, 10)
	}

	if !okB || !okE || !okM {
		return "", fmt.Errorf("invalid input: b = %s, e = %s, m = %s - cannot create all the integers required, from this input", b, e, m)
	}

	// Perform the exponentiation: base^exp % mod
	result := new(big.Int).Exp(base, exp, mod)

	// Scale factor for the result: 10^{scaleFactor * exp}
	scaleFactorBig := big.NewInt(int64(scaleFactor))
	scaleFactorExp := new(big.Int).Mul(scaleFactorBig, exp)

	// Calculate 10^{scaleFactor * exp}
	scaleDivisor := new(big.Int).Exp(big.NewInt(10), scaleFactorExp, nil)

	// Convert result and scaleDivisor to strings for division
	resultStr := result.String()
	scaleDivisorStr := scaleDivisor.String()

	// Perform the final division
	finalResult, err := Divide(resultStr, scaleDivisorStr)
	if err != nil {
		return "", fmt.Errorf("error during division: %v", err)
	}

	return finalResult, nil
}

// Taylor series expansion for exp(x) = 1 + x + x^2/2! + x^3/3! + ...
func expSeries(x *big.Float, prec uint) *big.Float {
	// Set precision
	x.SetPrec(prec)
	// Start with the first term in the series (1)
	result := big.NewFloat(1).SetPrec(prec)
	term := big.NewFloat(1).SetPrec(prec)
	factorial := big.NewFloat(1).SetPrec(prec)

	// Temporary variables
	tmp := new(big.Float).SetPrec(prec)
	for i := 1; i < 50; i++ { // Adjust the number of terms for better precision
		factorial = factorial.Mul(factorial, big.NewFloat(float64(i)))
		term = term.Mul(term, x)
		tmp.Quo(term, factorial)
		result.Add(result, tmp)
	}

	return result
}

// Logarithm approximation using Newton's method
func logNewton(x *big.Float, prec uint) *big.Float {
	if x.Cmp(big.NewFloat(0)) <= 0 {
		panic("Logarithm only defined for positive numbers")
	}

	// Initial guess (logarithm approximation)
	approx := big.NewFloat(0.5).SetPrec(prec)

	// Newton's method to refine the logarithm estimate
	two := big.NewFloat(2).SetPrec(prec)
	tmp := new(big.Float).SetPrec(prec)
	for i := 0; i < 20; i++ { // Adjust iterations for better precision
		expApprox := expSeries(approx, prec)
		tmp.Quo(x, expApprox)
		tmp.Sub(tmp, big.NewFloat(1))
		tmp.Quo(tmp, two)
		approx.Add(approx, tmp)
	}
	return approx
}

// ExpFloat computes a^b = exp(b * ln(a)) using arbitrary precision with big.Float
func ExpFloat(a, b string, prec uint) (string, error) {
	// Parse inputs as big.Float
	base, okA := new(big.Float).SetString(a)
	exp, okB := new(big.Float).SetString(b)

	if !okA || !okB {
		return "", fmt.Errorf("invalid input: a = %s, b = %s", a, b)
	}

	// Ensure base > 0 since ln(a) for non-positive a is undefined
	if base.Cmp(big.NewFloat(0)) <= 0 {
		return "", fmt.Errorf("base must be greater than 0")
	}

	// Compute ln(a) using Newton's method for big.Float
	logBase := logNewton(base, prec)

	// Compute b * ln(a)
	expLogBase := new(big.Float).Mul(exp, logBase)

	// Compute exp(b * ln(a)) using Taylor series for big.Float
	result := expSeries(expLogBase, prec)

	return result.Text('f', int(prec)), nil
}

// Mod performs the modulus operation of a by m.
// big.Int Mod sets z to the modulus x%y for y != 0 and returns z.
// If y == 0, a division-by-zero run-time panic occurs.
// Mod implements Euclidean modulus (unlike Go); see [Int.DivMod] for more details.
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
