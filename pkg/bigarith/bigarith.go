package bigarith

import (
	"fmt"
	"math/big"
	"strconv"
)

// Error refactor functions

// bigFloat converts the string representation of the bigarith.Float into a *big.Float
// An internal only function to simplify other code
// NOTE: Panics on bad string representation of big.Float
func bigFloat(a string) *big.Float {
	// in order to create a simplified application of string <-> big.*
	// no errors are handed back up the chain
	// instead simply Panic here with clear reasons why
	x, ok := new(big.Float).SetPrec(2048).SetString(a)
	if !ok {
		panic(fmt.Sprintf("Error creating big.Float using SetString - base = %d - string value = %s", 10, a))
	}
	return x
}

// bigInt converts the string representation of the bigarith.Int into a *big.Int
// An internal only function to simplify other code
// NOTE: Panics on bad string representation of big.Int
func bigInt(a string) *big.Int {
	// in order to create a simplified application of string <-> big.*
	// no errors are handed back up the chain
	// instead simply Panic here with clear reasons why
	x, ok := new(big.Int).SetString(a, 10)
	if !ok {
		panic(fmt.Sprintf("Error creating big.Int using SetString - base = %d - string value = %s", 10, a))
	}
	return x
}

// Logarithm approximation using Newton's method
// internal function only - just works in big.Float
func logNewton(x *big.Float) *big.Float {
	if x.Cmp(big.NewFloat(0)) <= 0 {
		panic("Logarithm only defined for positive numbers")
	}

	// Initial guess (logarithm approximation)
	approx := bigFloat("0.5")

	// Newton's method to refine the logarithm estimate
	two := bigFloat("2")
	tmp := bigFloat("0")
	for i := 0; i < 20; i++ { // Adjust iterations for better precision
		expApprox := expSeries(approx)
		tmp.Quo(x, expApprox)
		tmp.Sub(tmp, big.NewFloat(1))
		tmp.Quo(tmp, two)
		approx.Add(approx, tmp)
	}
	return approx
}

// Taylor series expansion for exp(x) = 1 + x + x^2/2! + x^3/3! + ...
// internal function only - just works in big.Float
func expSeries(x *big.Float) *big.Float {
	// Start with the first term in the series (1)
	result := bigFloat("1")
	term := bigFloat("1")
	factorial := bigFloat("1")

	tmp := bigFloat("0")      // Temporary variable for addition step
	for i := 1; i < 50; i++ { // Adjust the number of terms for better precision
		factorial = factorial.Mul(factorial, bigFloat(strconv.Itoa(i)))
		term = term.Mul(term, x)
		tmp.Quo(term, factorial)
		result.Add(result, tmp)
	}

	return result
}

// ===================
// Original functions

// // Add takes two string representations of integers, adds them, and returns the result as a string.
// // Returns an error if the input strings are not valid integers.
// func Add(a, b string) (string, error) {
// 	x, okX := new(big.Int).SetString(a, 10)
// 	y, okY := new(big.Int).SetString(b, 10)
// 	if !okX || !okY {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return new(big.Int).Add(x, y).String(), nil
// }

// // AddFloat takes two string representations of floats, adds them, and returns the result as a string.
// // Returns an error if the input strings are not valid integers.
// func AddFloat(a, b string) (string, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	y, okY := new(big.Float).SetPrec(2048).SetString(b)
// 	if !okX || !okY {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return new(big.Float).SetPrec(2048).Add(x, y).Text('f', int(2048)), nil
// }

// // Subtract takes two string representations of integers, subtracts the second from the first,
// // and returns the result as a string.
// // Returns an error if the input strings are not valid integers.
// func Subtract(a, b string) (string, error) {
// 	x, okX := new(big.Int).SetString(a, 10)
// 	y, okY := new(big.Int).SetString(b, 10)
// 	if !okX || !okY {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return new(big.Int).Sub(x, y).String(), nil
// }

// // SubtractFloat takes two string representations of floats, subtracts the second from the first,
// // and returns the result as a string.
// // Returns an error if the input strings are not valid integers.
// func SubtractFloat(a, b string) (string, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	y, okY := new(big.Float).SetPrec(2048).SetString(b)
// 	if !okX || !okY {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return new(big.Float).SetPrec(2048).Sub(x, y).Text('f', int(2048)), nil
// }

// // Multiply takes two string representations of integers, multiplies them,
// // and returns the result as a string.
// // Returns an error if the input strings are not valid integers.
// // big.Int Mul sets z to the product x*y and returns z.
// func Multiply(a, b string) (string, error) {
// 	x, okX := new(big.Int).SetString(a, 10)
// 	y, okY := new(big.Int).SetString(b, 10)
// 	if !okX || !okY {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return new(big.Int).Mul(x, y).String(), nil
// }

// // Multiply takes two string representations of floats, multiplies them,
// // and returns the result as a string.
// // Returns an error if the input strings are not valid floats.
// // big.Float Mul sets z to the rounded product x*y and returns z.
// // Precision, rounding, and accuracy reporting are as for [Float.Add].
// // Mul panics with [ErrNaN] if one operand is zero and the other operand an infinity.
// // The value of z is undefined in that case.
// func MultiplyFloat(a, b string) (string, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	y, okY := new(big.Float).SetPrec(2048).SetString(b)
// 	if !okX || !okY {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the floats required, from this input", a, b)
// 	}
// 	return new(big.Float).Mul(x, y).Text('f', int(2048)), nil
// }

// // DivideInField takes two string representations of integers, divides a by b,
// // and returns the result as a string in a finite field defined by a prime modulus p.
// // Returns an error if the input strings are not valid integers, or if b has no inverse modulo p,
// // (which includes if b is 0).
// func DivideInField(a, b, p string) (string, error) {
// 	// First, find the multiplicative inverse of b mod p
// 	bInv, err := ModularInverse(b, p)
// 	if err != nil {
// 		return "", fmt.Errorf("error finding inverse: %v", err)
// 	}

// 	// Calculate a * bInv mod p
// 	result, err := Multiply(a, bInv) // a * b^-1
// 	if err != nil {
// 		return "", fmt.Errorf("error nultiplying with modular inverse: %v", err)
// 	}
// 	// x, okX := new(big.Int).SetString(a, 10)
// 	resultModP, err := Mod(result, p) // mod p
// 	if err != nil {
// 		return "", fmt.Errorf("error finding mod p of result: %v", err)
// 	}
// 	return resultModP, nil
// }

// // ModularInverse calculates the multiplicative inverse of a modulo p using Fermat's Little Theorem
// // and returns it as a string. Returns an error if p is not prime or a and p are not coprime.
// func ModularInverse(a, p string) (string, error) {
// 	aInt, okA := new(big.Int).SetString(a, 10)
// 	pInt, okP := new(big.Int).SetString(p, 10)
// 	zero := big.NewInt(0)
// 	if !okA || !okP {
// 		return "", fmt.Errorf("invalid input: a = %s, p = %s - cannot create all the integers required, from this input", a, p)
// 	}
// 	if aInt.Cmp(zero) == 0 {
// 		return "", fmt.Errorf("invalid input: a is ZERO - no modular multiplicative inverse")
// 	}

// 	// Check if p is prime; if not, the multiplicative inverse might not exist
// 	// https://pkg.go.dev/math/big#Int.ProbablyPrime
// 	// From ---^ that site: "The probability of returning true for a randomly chosen non-prime is at most ¼ⁿ"
// 	// i.e. (1/4)^n - where n is the parameter handed in to the function
// 	// TODO: work out if there's a way to decide what that param should be set to
// 	if !pInt.ProbablyPrime(100) {
// 		return "", fmt.Errorf("invalid input: modulus %s is not prime", p)
// 	}

// 	// Calculate a^(p-2) mod p
// 	pMinusTwo := new(big.Int).Sub(pInt, big.NewInt(2)) // p-2
// 	return new(big.Int).Exp(aInt, pMinusTwo, pInt).String(), nil
// }

// // Divide takes two string representations of floats
// // and returns a / b as a string
// func Divide(a, b string) (string, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	y, okY := new(big.Float).SetPrec(2048).SetString(b)
// 	if !okX || !okY {
// 		return "0", fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return new(big.Float).SetPrec(2048).Quo(x, y).Text('f', int(2048)), nil
// }

// // Neg takes a float represented as a string
// // and returns its negative represented as a string
// func Neg(a string) (string, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	if !okX {
// 		return "0", fmt.Errorf("invalid input: a = %s - cannot create all the integers required, from this input", a)
// 	}
// 	return new(big.Float).SetPrec(2048).Neg(x).Text('f', int(2048)), nil
// }

// // Cmp compares two integers represented as strings and returns:
// // -1 if a < b, 0 if a == b, 1 if a > b
// func Cmp(a, b string) (int, error) {
// 	x, okX := new(big.Int).SetString(a, 10)
// 	y, okY := new(big.Int).SetString(b, 10)
// 	if !okX || !okY {
// 		return 0, fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return x.Cmp(y), nil
// }

// // CmpFloat compares two floats represented as strings and returns:
// // -1 if a < b, 0 if a == b, 1 if a > b
// func CmpFloat(a, b string) (int, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	y, okY := new(big.Float).SetPrec(2048).SetString(b)
// 	if !okX || !okY {
// 		return 0, fmt.Errorf("invalid input: a = %s, b = %s - cannot create all the integers required, from this input", a, b)
// 	}
// 	return x.Cmp(y), nil
// }

// // Exp computes the exponentiation of base b to exp e modulo m - represented as strings.
// // If m == "", returns b**e unless y <= 0 then returns "1".
// // big.Int Exp sets z = x**y mod |m| (i.e. the sign of m is ignored), and returns z.
// // If m == nil or m == 0, z = x**y unless y <= 0 then z = 1.
// // If m != 0, y < 0, and x and m are not relatively prime, z is unchanged and nil is returned.
// func Exp(b, e, m string) (string, error) {
// 	// Check if b is already an integer, and if not, scale it
// 	decimalIndex := strings.Index(b, ".")
// 	scaleFactor := 0

// 	if decimalIndex != -1 {
// 		if m != "" {
// 			return "", fmt.Errorf("cannot have non-integer base AND a modulus - either base must be integer OR there must be no modulus value: b = %s, e = %s, m = %s - cannot create all the integers required, from this input", b, e, m)
// 		}
// 		// b IS NOT an integer
// 		scaleFactor = len(b) - decimalIndex - 1
// 		b = strings.Replace(b, ".", "", 1) // Remove the decimal point
// 	}

// 	// Parse base and exponent as big.Int
// 	base, okB := new(big.Int).SetString(b, 10)
// 	exp, okE := new(big.Int).SetString(e, 10)

// 	var mod *big.Int
// 	var okM bool
// 	if m == "" {
// 		mod = nil
// 		okM = true
// 	} else {
// 		mod, okM = new(big.Int).SetString(m, 10)
// 	}

// 	if !okB || !okE || !okM {
// 		return "", fmt.Errorf("invalid input: b = %s, e = %s, m = %s - cannot create all the integers required, from this input", b, e, m)
// 	}

// 	// Perform the exponentiation: base^exp % mod
// 	result := new(big.Int).Exp(base, exp, mod)

// 	// Scale factor for the result: 10^{scaleFactor * exp}
// 	scaleFactorBig := big.NewInt(int64(scaleFactor))
// 	scaleFactorExp := new(big.Int).Mul(scaleFactorBig, exp)

// 	// Calculate 10^{scaleFactor * exp}
// 	scaleDivisor := new(big.Int).Exp(big.NewInt(10), scaleFactorExp, nil)

// 	// Convert result and scaleDivisor to strings for division
// 	resultStr := result.String()
// 	scaleDivisorStr := scaleDivisor.String()

// 	// Perform the final division
// 	finalResult, err := Divide(resultStr, scaleDivisorStr)
// 	if err != nil {
// 		return "", fmt.Errorf("error during division: %v", err)
// 	}

// 	return finalResult, nil
// }

// // Sqrt calculates the square root of x with arbitrary precision using Newton's method.
// func Sqrt(a string, prec uint) (string, error) {
// 	x, okX := new(big.Float).SetPrec(prec).SetString(a)
// 	if !okX {
// 		return "0", fmt.Errorf("invalid input: a = %s - cannot create all the integers required, from this input", a)
// 	}

// 	// Initial guess: x / 2
// 	guess := new(big.Float).Quo(x, big.NewFloat(2).SetPrec(prec))

// 	// Variables for the iteration
// 	two := big.NewFloat(2).SetPrec(prec)
// 	tmp := new(big.Float).SetPrec(prec)

// 	// Iterate using Newton's method
// 	for i := 0; i < 10000; i++ { // Increase iterations for better precision
// 		// tmp = x / guess
// 		tmp.Quo(x, guess)

// 		// guess = (guess + tmp) / 2
// 		guess.Add(guess, tmp).Quo(guess, two)
// 	}

// 	return guess.Text('f', int(prec)), nil
// }

// // CubeRoot calculates the cube root of x with arbitrary precision using Newton's method.
// func CubeRoot(a string, prec uint) (string, error) {
// 	x, okX := new(big.Float).SetPrec(2048).SetString(a)
// 	if !okX {
// 		return "0", fmt.Errorf("invalid input: a = %s - cannot create all the integers required, from this input", a)
// 	}

// 	// Set precision
// 	x.SetPrec(prec)

// 	// Initial guess: start with x/3 (a rough approximation)
// 	guess := new(big.Float).Quo(x, big.NewFloat(3)).SetPrec(prec)

// 	// Variables used in the iteration
// 	three := big.NewFloat(3).SetPrec(prec)
// 	two := big.NewFloat(2).SetPrec(prec)
// 	tmp := new(big.Float).SetPrec(prec)

// 	// Iterate using Newton's method
// 	for i := 0; i < 10000; i++ { // Increase iterations for higher precision
// 		// y^2
// 		tmp.Mul(guess, guess)

// 		// x / y^2
// 		tmp.Quo(x, tmp)

// 		// (2 * y + x / y^2) / 3
// 		guess.Mul(guess, two).Add(guess, tmp).Quo(guess, three)
// 	}

// 	return guess.Text('f', int(prec)), nil
// }

// // ExpFloat computes a^b = exp(b * ln(a)) using arbitrary precision with big.Float
// func ExpFloat(a, b string) (string, error) {
// 	// Parse inputs as big.Float
// 	base, okA := new(big.Float).SetPrec(2048).SetString(a)
// 	exp, okB := new(big.Float).SetPrec(2048).SetString(b)

// 	if !okA || !okB {
// 		return "", fmt.Errorf("invalid input: a = %s, b = %s", a, b)
// 	}

// 	// Ensure base > 0 since ln(a) for non-positive a is undefined
// 	// Convert to positive - convert back at end
// 	baseIsNonNegative := true
// 	if base.Cmp(big.NewFloat(0)) <= 0 {
// 		baseIsNonNegative = false
// 		base = new(big.Float).Abs(base)
// 	}

// 	// Compute ln(a) using Newton's method for big.Float
// 	logBase := logNewton(base)

// 	// Compute b * ln(a)
// 	expLogBase := new(big.Float).SetPrec(2048).Mul(exp, logBase)

// 	// Compute exp(b * ln(a)) using Taylor series for big.Float
// 	result := expSeries(expLogBase)
// 	if !baseIsNonNegative {
// 		result = new(big.Float).SetPrec(2048).Neg(result)
// 	}

// 	return result.Text('f', 2048), nil
// }

// // Mod performs the modulus operation of a by m.
// // big.Int Mod sets z to the modulus x%y for y != 0 and returns z.
// // If y == 0, a division-by-zero run-time panic occurs.
// // Mod implements Euclidean modulus (unlike Go); see [Int.DivMod] for more details.
// func Mod(a, m string) (string, error) {
// 	x, okX := new(big.Int).SetString(a, 10)
// 	mod, okM := new(big.Int).SetString(m, 10)
// 	if !okX || !okM {
// 		return "", fmt.Errorf("invalid input: a = %s, m = %s - cannot create all the integers required, from this input", a, m)
// 	}
// 	result := new(big.Int).Mod(x, mod)
// 	return result.String(), nil
// }

// // Factorial calculates the factorial of an integer a, represented as a string,
// // and returns the result as a string. Returns an error if the input string is not a valid integer.
// func Factorial(a string) (string, error) {
// 	x, ok := new(big.Int).SetString(a, 10)
// 	if !ok {
// 		return "", fmt.Errorf("invalid input: a = %s - cannot create the integer required, from this input", a)
// 	}

// 	// Factorial calculation
// 	result := big.NewInt(1)
// 	one := big.NewInt(1)
// 	for x.Cmp(one) > 0 {
// 		result.Mul(result, x)
// 		x.Sub(x, one)
// 	}

// 	return result.String(), nil
// }

// // Sin calculates the sine of x (in radians), represented as a string, using the Taylor series expansion.
// // It returns the result as a string with the specified precision.
// func Sin(a string, prec uint) (string, error) {
// 	x, ok := new(big.Float).SetPrec(prec).SetString(a)
// 	if !ok {
// 		return "", fmt.Errorf("invalid input: a = %s - cannot create all the floats required, from this input", a)
// 	}

// 	// Variables for the Taylor series
// 	result := new(big.Float).SetPrec(prec)
// 	term := new(big.Float).SetPrec(prec).Set(x)        // First term is x
// 	xSquared := new(big.Float).Mul(x, x).SetPrec(prec) // x^2

// 	// one := big.NewFloat(1).SetPrec(prec)
// 	factorial := big.NewFloat(1).SetPrec(prec) // 1!

// 	// Iterate through Taylor series terms
// 	sign := 1.0
// 	for i := int64(1); i < 50; i++ {
// 		if i > 1 {
// 			// Update factorial: (2*i - 1)(2*i)
// 			factorial.Mul(factorial, big.NewFloat(float64(2*i-2)))
// 			factorial.Mul(factorial, big.NewFloat(float64(2*i-1)))

// 			// Update term: x^(2i+1) / factorial
// 			term.Mul(term, xSquared)
// 			term.Quo(term, factorial)
// 		}

// 		// Add or subtract the term
// 		if sign == 1 {
// 			result.Add(result, term)
// 		} else {
// 			result.Sub(result, term)
// 		}

// 		// Toggle the sign for the next term
// 		sign = -sign
// 	}

// 	return result.Text('f', int(prec)), nil
// }

// // ModFloat calculates a mod b for big.Float, analogous to a % b for floating-point numbers.
// func ModFloat(a, b *big.Float, prec uint) *big.Float {
// 	// Create a result big.Float to store the modulo result
// 	result := new(big.Float).SetPrec(prec)

// 	// Compute quotient q = a / b
// 	q := new(big.Float).SetPrec(prec).Quo(a, b)

// 	// Truncate q to an integer (floor it)
// 	qInt, _ := q.Int(nil)

// 	// Compute result = a - qInt * b
// 	bq := new(big.Float).SetPrec(prec).Mul(new(big.Float).SetInt(qInt), b)
// 	result.Sub(a, bq)

// 	return result
// }

// func Mod2Pi(a string, prec uint) (string, error) {
// 	// Use a more precise Pi value for accurate modulo calculations
// 	pi := Pi(2048)
// 	twoPi := new(big.Float).Mul(pi, big.NewFloat(2))

// 	x, ok := new(big.Float).SetPrec(prec).SetString(a)
// 	if !ok {
// 		return "", fmt.Errorf("invalid input: a = %s - cannot create the float", a)
// 	}

// 	// Perform x mod 2*Pi
// 	modResult := ModFloat(x, twoPi, 2048)

// 	// Ensure it's in the range [0, 2*Pi]
// 	if modResult.Cmp(big.NewFloat(0)) < 0 {
// 		modResult.Add(modResult, twoPi)
// 	}

// 	return modResult.Text('f', int(prec)), nil
// }

// func Pi(prec uint) *big.Float {
// 	// Constants in the Chudnovsky algorithm
// 	// C = 426880 * sqrt(10005)
// 	C := new(big.Float).SetPrec(prec).SetFloat64(426880)
// 	K := new(big.Float).SetPrec(prec).SetFloat64(13591409)
// 	M := new(big.Float).SetPrec(prec).SetFloat64(1)
// 	L := new(big.Float).SetPrec(prec).SetFloat64(13591409)
// 	X := new(big.Float).SetPrec(prec).SetFloat64(1)
// 	S := new(big.Float).SetPrec(prec).Set(L)

// 	// Multiply C by sqrt(10005)
// 	sqrt10005 := new(big.Float).SetPrec(prec).SetFloat64(10005)
// 	sqrt10005.Sqrt(sqrt10005)
// 	C.Mul(C, sqrt10005)

// 	// Set factorials
// 	factK := new(big.Int).SetInt64(1)

// 	// Iterations to improve precision (increase this for more digits)
// 	iterations := 5

// 	// Use Chudnovsky's series to compute Pi
// 	for i := 1; i < iterations; i++ {
// 		// Update K, M, L, X
// 		k := float64(i)
// 		K.Mul(K, new(big.Float).SetFloat64(545140134))
// 		M.Mul(M, new(big.Float).SetFloat64(k*k*k))

// 		L.Add(L, new(big.Float).SetFloat64(545140134))
// 		X.Mul(X, new(big.Float).SetFloat64(-262537412640768000))

// 		// Calculate term and add to the sum S
// 		term := new(big.Float).SetPrec(prec).Quo(new(big.Float).Mul(M, L), X)
// 		term.Quo(term, new(big.Float).SetPrec(prec).SetInt(factK))
// 		S.Add(S, term)

// 		// Update the factorial
// 		factK.Mul(factK, big.NewInt(int64(6*i)))
// 	}

// 	// Calculate pi = C / S
// 	pi := new(big.Float).SetPrec(prec).Quo(C, S)

// 	return pi
// }

// // Cos calculates the cosine of x (in radians), represented as a string, using the Taylor series expansion.
// // It returns the result as a string with the specified precision.
// func Cos(a string, prec uint) (string, error) {
// 	// First, reduce the input to [0, 2*Pi]
// 	modA, err := Mod2Pi(a, prec)
// 	if err != nil {
// 		return "", fmt.Errorf("error reducing a to [0, 2*Pi]: %v", err)
// 	}

// 	// Now continue with the Taylor series expansion as before
// 	x, ok := new(big.Float).SetPrec(prec).SetString(modA)
// 	if !ok {
// 		return "", fmt.Errorf("invalid input: a = %s - cannot create the float", a)
// 	}

// 	// Variables for the Taylor series
// 	result := new(big.Float).SetPrec(prec).Set(big.NewFloat(1)) // First term is 1
// 	term := big.NewFloat(1).SetPrec(prec)
// 	xSquared := new(big.Float).Mul(x, x).SetPrec(prec) // x^2

// 	factorial := big.NewFloat(1).SetPrec(prec) // 0! = 1

// 	// Iterate through Taylor series terms
// 	sign := -1.0
// 	for i := int64(1); i < 50; i++ {
// 		// Update factorial: (2*i - 1)(2*i)
// 		factorial.Mul(factorial, big.NewFloat(float64(2*i-1)))
// 		factorial.Mul(factorial, big.NewFloat(float64(2*i)))

// 		// Update term: x^(2i) / factorial
// 		term.Mul(term, xSquared)
// 		term.Quo(term, factorial)

// 		// Add or subtract the term
// 		if sign == 1 {
// 			result.Add(result, term)
// 		} else {
// 			result.Sub(result, term)
// 		}

// 		// Toggle the sign for the next term
// 		sign = -sign
// 	}

// 	// Ensure the result stays between -1 and 1 (due to precision errors in float handling)
// 	one := big.NewFloat(1).SetPrec(prec)
// 	negOne := new(big.Float).Neg(one)

// 	if result.Cmp(one) > 0 {
// 		result.Set(one)
// 	} else if result.Cmp(negOne) < 0 {
// 		result.Set(negOne)
// 	}

// 	return result.Text('f', int(prec)), nil
// }

// // Asin calculates the arcsine (inverse sine) of x, represented as a string, using a Taylor series expansion.
// // It returns the result as a string with the specified precision.
// func Asin(a string, prec uint) (string, error) {
// 	x, ok := new(big.Float).SetPrec(prec).SetString(a)
// 	if !ok {
// 		return "", fmt.Errorf("invalid input: a = %s - cannot create all the floats required, from this input", a)
// 	}

// 	// Variables for the Taylor series
// 	result := new(big.Float).SetPrec(prec).Set(x)
// 	term := new(big.Float).SetPrec(prec).Set(x)
// 	one := big.NewFloat(1).SetPrec(prec)
// 	factorial := new(big.Float).SetPrec(prec).Set(one)

// 	xSquared := new(big.Float).Mul(x, x).SetPrec(prec)

// 	// Iterate through the Taylor series terms
// 	for i := 1; i < 50; i++ {
// 		// Multiply the term by x^2
// 		term.Mul(term, xSquared)

// 		// Update factorial: (2*i - 1) * (2*i)
// 		factorial.Mul(factorial, big.NewFloat(float64(2*i-1)))
// 		factorial.Mul(factorial, big.NewFloat(float64(2*i)))

// 		// Add the next term in the series
// 		nextTerm := new(big.Float).Quo(term, factorial)
// 		if i%2 == 0 {
// 			result.Add(result, nextTerm)
// 		} else {
// 			result.Sub(result, nextTerm)
// 		}
// 	}

// 	return result.Text('f', int(prec)), nil
// }

// // Acos calculates the arccosine (inverse cosine) of x, represented as a string, using the identity:
// // acos(x) = Pi/2 - asin(x). It returns the result as a string with the specified precision.
// func Acos(a string, prec uint) (string, error) {
// 	// First, calculate asin(x)
// 	asinResult, err := Asin(a, prec)
// 	if err != nil {
// 		return "", fmt.Errorf("error calculating asin: %v", err)
// 	}

// 	asinFloat, ok := new(big.Float).SetPrec(prec).SetString(asinResult)
// 	if !ok {
// 		return "", fmt.Errorf("invalid input for asin result")
// 	}

// 	// Calculate Pi/2
// 	pi := Pi(prec)
// 	piOverTwo := new(big.Float).Quo(pi, big.NewFloat(2))

// 	// Acos(x) = Pi/2 - Asin(x)
// 	result := new(big.Float).Sub(piOverTwo, asinFloat)

// 	return result.Text('f', int(prec)), nil
// }
