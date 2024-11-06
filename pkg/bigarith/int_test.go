package bigarith

import (
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Helper function to generate random integers as strings
func randomInt() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(1000000000) // Random up to a billion
	return strconv.FormatInt(n, 10)
}

// Helper function to generate random integers as strings
func randomSmallInt(i int64) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(i) // Random up to i
	return strconv.FormatInt(n, 10)
}

// Helper function to check if a number is prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Convert the bigarith Int to int64 for comparison
func intFromString(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func TestPlus(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		b := randomInt()

		// Parse them to int64 for verification
		intA, err := intFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}
		intB, err := intFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Expected result using standard math
		expectedResult := intA + intB

		// Perform the Plus operation
		intObj := NewInt(a)
		originalValue := intObj.Val() // Store original value for later comparison
		result := intObj.Plus(b)

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("Plus failed: %d + %d = %d, expected %d", intA, intB, resultInt, expectedResult)
		}

		// Verify that the internal value of the original object hasn't been altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestMinus(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		b := randomInt()

		// Parse them to int64 for verification
		intA, err := intFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}
		intB, err := intFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Expected result using standard math
		expectedResult := intA - intB

		// Perform the Minus operation
		intObj := NewInt(a)
		originalValue := intObj.Val() // Store original value for later comparison
		result := intObj.Minus(b)

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("Minus failed: %d - %d = %d, expected %d", intA, intB, resultInt, expectedResult)
		}

		// Verify that the internal value of the original object hasn't been altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestNeg(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integer as string
		a := randomInt()

		// Parse it to int64 for verification
		intA, err := intFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Expected result using standard math
		expectedResult := -intA

		// Perform the Neg operation
		intObj := NewInt(a)
		originalValue := intObj.Val()
		result := intObj.Neg()

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("Neg failed: -%d = %d, expected %d", intA, resultInt, expectedResult)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestTimes(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		b := randomInt()

		// Parse them to int64 for verification
		intA, err := intFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}
		intB, err := intFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Expected result using standard math
		expectedResult := intA * intB

		// Perform the Times operation
		intObj := NewInt(a)
		originalValue := intObj.Val()
		result := intObj.Times(b)

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("Times failed: %d * %d = %d, expected %d", intA, intB, resultInt, expectedResult)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestMod(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		b := randomInt()

		// Parse them to int64 for verification
		intA, err := intFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}
		intB, err := intFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Avoid divide by zero
		if intB == 0 {
			intB = 1
			b = "1"
		}

		// Expected result using standard math
		expectedResult := intA % intB

		// Perform the Mod operation
		intObj := NewInt(a)
		originalValue := intObj.Val()
		result := intObj.Mod(b)

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("Mod failed: %d %% %d = %d, expected %d", intA, intB, resultInt, expectedResult)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestToThePowerOf(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Limit base to values between 2 and 100
		base := strconv.Itoa(rand.Intn(99) + 2) // random base in range [2, 100]
		// Limit exponent to values between 0 and 10
		exponent := strconv.Itoa(rand.Intn(10)) // random exponent in range [0, 10]

		// Parse them to int64 for verification
		intBase, err := intFromString(base)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}
		intExponent, err := strconv.Atoi(exponent)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Calculate expected result using standard math
		expectedResult := int64(1)
		for j := 0; j < intExponent; j++ {
			expectedResult *= intBase
		}

		// Perform the ToThePowerOf operation
		intObj := NewInt(base)
		originalValue := intObj.Val()
		result := intObj.ToThePowerOf(exponent, "0") // use modulus "0" for basic exponentiation

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("ToThePowerOf failed: %d ^ %d = %d, expected %d", intBase, intExponent, resultInt, expectedResult)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestFactorial(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Limit the input range to prevent overflow for factorial calculations
		a := strconv.Itoa(rand.Intn(20) + 1)

		// Parse it to int64 for verification
		intA, err := strconv.Atoi(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Expected result using standard math
		expectedResult := int64(1)
		for j := 1; j <= intA; j++ {
			expectedResult *= int64(j)
		}

		// Perform the Factorial operation
		intObj := NewInt(a)
		originalValue := intObj.Val()
		result := intObj.Factorial()

		// Convert result to int64 for comparison
		resultInt, err := intFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to int: %v", err)
		}

		// Verify the result
		if resultInt != expectedResult {
			t.Errorf("Factorial failed: %d! = %d, expected %d", intA, resultInt, expectedResult)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestDividedBy(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		b := randomInt()

		// Parse them to int64 for verification
		intA, err := intFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}
		intB, err := intFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to int: %v", err)
		}

		// Avoid divide by zero
		if intB == 0 {
			intB = 1
			b = "1"
		}

		// Expected result using standard math
		expectedResult := float64(intA) / float64(intB)

		// Perform the DividedBy operation
		intObj := NewInt(a)
		originalValue := intObj.Val()
		result := intObj.DividedBy(b)

		resultNumFloat, _ := strconv.ParseFloat(strings.TrimSpace(result.Num().strVal), 64)
		resultDenomFloat, _ := strconv.ParseFloat(strings.TrimSpace(result.Denom().strVal), 64)
		floatResult := resultNumFloat / resultDenomFloat

		// Verify the result
		if floatResult != expectedResult {
			t.Errorf("DividedBy failed - result: %f, expected %f", floatResult, expectedResult)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}

func TestModularInverse_RelativelyPrime(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		p := randomInt()

		// Convert a and p to big.Int for GCD check
		aBigInt, _ := new(big.Int).SetString(a, 10)
		pBigInt, _ := new(big.Int).SetString(p, 10)
		gcd := new(big.Int).GCD(nil, nil, aBigInt, pBigInt)

		// Ensure they are relatively prime by adjusting a until GCD(a, p) == 1
		for gcd.Cmp(big.NewInt(1)) != 0 {
			aBigInt.Add(aBigInt, big.NewInt(1))
			a = aBigInt.String()
			gcd = new(big.Int).GCD(nil, nil, aBigInt, pBigInt)
		}

		// Perform ModularInverse
		intObj := NewInt(a)
		result := intObj.ModularInverse(p)

		// Verify the result: (result * a) % p should be 1
		resultBigInt := new(big.Int).Mul(bigInt(result.Val()), aBigInt)
		expected := new(big.Int).Mod(resultBigInt, pBigInt)

		if expected.Cmp(big.NewInt(1)) != 0 {
			t.Errorf("ModularInverse failed: (%s * %s) mod %s != 1", result.Val(), a, p)
		}
	}
}

func TestModularInverse_NotRelativelyPrime(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate random integers as strings
		a := randomInt()
		p := randomInt()

		// Force a and p to not be relatively prime by making p a multiple of a
		aBigInt, _ := new(big.Int).SetString(a, 10)
		pBigInt, _ := new(big.Int).SetString(p, 10)
		pBigInt.Mul(pBigInt, aBigInt)
		p = pBigInt.String()

		intObj := NewInt(a)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for a = %s, p = %s (not relatively prime), but no panic occurred", a, p)
			}
		}()

		// This should panic
		intObj.ModularInverse(p)
	}
}

func TestDivideInField_RelativelyPrime(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate 3 random integers
		a := randomInt()
		d := randomInt()
		p := randomInt()

		// Convert them to big.Int for comparison and ordering
		aBigInt, _ := new(big.Int).SetString(a, 10)
		dBigInt, _ := new(big.Int).SetString(d, 10)
		pBigInt, _ := new(big.Int).SetString(p, 10)

		// Ensure the modulus is the largest value
		if pBigInt.Cmp(aBigInt) < 0 || pBigInt.Cmp(dBigInt) < 0 {
			// Set p to be the largest by swapping
			if aBigInt.Cmp(dBigInt) > 0 {
				aTempInt := new(big.Int).Set(pBigInt)
				pBigInt.Set(aBigInt)
				aBigInt.Set(aTempInt)
			} else {
				aTempInt := new(big.Int).Set(pBigInt)
				pBigInt.Set(dBigInt)
				dBigInt.Set(aTempInt)
			}
			a = aBigInt.String()
			d = dBigInt.String()
			p = pBigInt.String()

		}

		// Now make sure d and p are relatively prime
		gcd := new(big.Int).GCD(nil, nil, dBigInt, pBigInt)

		for gcd.Cmp(big.NewInt(1)) != 0 {
			dBigInt.Add(dBigInt, big.NewInt(1))
			d = dBigInt.String()
			gcd = new(big.Int).GCD(nil, nil, dBigInt, pBigInt)
		}

		// Perform DivideInField operation (a / d mod p)
		intObj := NewInt(a)
		result := intObj.DivideInField(d, p)

		// Verify the result: (result * d) mod p should be a
		resultMod := new(big.Int).Mul(bigInt(result.Val()), bigInt(d))
		expectedResult := new(big.Int).Mod(resultMod, bigInt(p))

		if expectedResult.Cmp(bigInt(a)) != 0 {
			t.Errorf("DivideInField failed: (%s * %s) mod %s != %s", result.Val(), d, p, a)
		}
	}
}

func TestDivideInField_PanicWhenPGreaterThan(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Generate 3 random integers
		a := randomInt()
		d := randomInt()
		p := randomInt()

		// Convert them to big.Int for comparison
		aBigInt, _ := new(big.Int).SetString(a, 10)
		dBigInt, _ := new(big.Int).SetString(d, 10)
		pBigInt, _ := new(big.Int).SetString(p, 10)

		// Ensure that p is smaller than a or d to cause panic
		if pBigInt.Cmp(aBigInt) > 0 || pBigInt.Cmp(dBigInt) > 0 {
			// Make p smaller than both a and d
			if aBigInt.Cmp(dBigInt) > 0 {
				pBigInt.Set(aBigInt.Sub(aBigInt, big.NewInt(1)))
			} else {
				pBigInt.Set(dBigInt.Sub(dBigInt, big.NewInt(1)))
			}
			p = pBigInt.String()
		}

		// Perform DivideInField operation and expect a panic
		intObj := NewInt(a)

		// Using defer and recover to capture the panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when p is less than i or d, but no panic occurred: i = %s, d = %s, p = %s", a, d, p)
			}
		}()

		// This should panic since p is smaller than a or d
		intObj.DivideInField(d, p)
	}
}

func TestProbablyPrime(t *testing.T) {
	// Known prime numbers
	primes := []string{"2", "3", "5", "7", "11", "13", "17", "19", "23", "29"}

	for _, prime := range primes {
		intObj := NewInt(prime)
		originalValue := intObj.Val()

		// Perform ProbablyPrime operation
		isPrime := intObj.ProbablyPrime()

		// Verify result
		if !isPrime {
			t.Errorf("ProbablyPrime failed: %s should be prime", prime)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}

	// Known non-prime numbers
	nonPrimes := []string{"4", "6", "8", "9", "10", "12", "14", "15", "16", "18"}

	for _, nonPrime := range nonPrimes {
		intObj := NewInt(nonPrime)
		originalValue := intObj.Val()

		// Perform ProbablyPrime operation
		isPrime := intObj.ProbablyPrime()

		// Verify result
		if isPrime {
			t.Errorf("ProbablyPrime failed: %s should not be prime", nonPrime)
		}

		// Verify the original object is not altered
		if intObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, intObj.Val())
		}
	}
}
