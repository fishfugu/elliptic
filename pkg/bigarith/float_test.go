package bigarith

import (
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// Helper function to generate random float numbers as strings
func randomFloat() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Float64() * 1000000 // Random up to a million
	return strconv.FormatFloat(n, 'f', 50, 64)
}

// Helper function to convert a float string to a float64 for comparison
func floatFromString(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func TestFloatPlus(t *testing.T) {
	for i := 0; i < 500; i++ {
		// Generate random floats as strings
		a := randomFloat()
		b := randomFloat()

		// Parse them to float64 for verification
		floatA, err := floatFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}
		floatB, err := floatFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}

		// Expected result using standard math
		expectedResult := floatA + floatB

		// Perform the Plus operation
		floatObj := NewFloat(a)
		originalValue := floatObj.Val() // Store original value for later comparison
		result := floatObj.Plus(b)

		// Convert result to float64 for comparison
		resultFloat, err := floatFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to float: %v", err)
		}

		// Verify the result
		if resultFloat != expectedResult {
			t.Errorf("Plus failed: %s + %s = %.50f + %.50f = %.50f, expected %.50f", a, b, floatA, floatB, resultFloat, expectedResult)
		}

		// Verify that the internal value of the original object hasn't been altered
		if floatObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, floatObj.Val())
		}
	}
}

func TestFloatMinus(t *testing.T) {
	for i := 0; i < 500; i++ {
		// Generate random floats as strings
		a := randomFloat()
		b := randomFloat()

		// Parse them to float64 for verification
		floatA, err := floatFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}
		floatB, err := floatFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}

		// Expected result using standard math
		expectedResult := floatA - floatB

		// Perform the Minus operation
		floatObj := NewFloat(a)
		originalValue := floatObj.Val() // Store original value for later comparison
		result := floatObj.Minus(b)

		// Convert result to float64 for comparison
		resultFloat, err := floatFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to float: %v", err)
		}

		// Verify the result
		if resultFloat != expectedResult {
			t.Errorf("Minus failed: %f - %f = %f, expected %f", floatA, floatB, resultFloat, expectedResult)
		}

		// Verify that the internal value of the original object hasn't been altered
		if floatObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, floatObj.Val())
		}
	}
}

func TestFloatTimes(t *testing.T) {
	for i := 0; i < 500; i++ {
		// Generate random floats as strings
		a := randomFloat()
		b := randomFloat()

		// Parse them to float64 for verification
		floatA, err := floatFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}
		floatB, err := floatFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}

		// Expected result using standard math
		expectedResult := floatA * floatB

		// Perform the Times operation
		floatObj := NewFloat(a)
		originalValue := floatObj.Val()
		result := floatObj.Times(b)

		// Convert result to float64 for comparison
		resultFloat, err := floatFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to float: %v", err)
		}

		// Verify the result
		if resultFloat != expectedResult {
			t.Errorf("Times failed: %f * %f = %f, expected %f", floatA, floatB, resultFloat, expectedResult)
		}

		// Verify the original object is not altered
		if floatObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, floatObj.Val())
		}
	}
}

func TestFloatDividedBy(t *testing.T) {
	for i := 0; i < 500; i++ {
		// Generate random floats as strings
		a := randomFloat()
		b := randomFloat()

		// Parse them to float64 for verification
		floatA, err := floatFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}
		floatB, err := floatFromString(b)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}

		// Avoid divide by zero
		if floatB == 0 {
			floatB = 1
			b = "1"
		}

		// Expected result using standard math
		expectedResult := floatA / floatB

		// Perform the DividedBy operation
		floatObj := NewFloat(a)
		originalValue := floatObj.Val()
		result := floatObj.DividedBy(b)

		// Convert result to float64 for comparison
		resultFloat, err := floatFromString(result.Val())
		if err != nil {
			t.Fatalf("Failed to convert result string to float: %v", err)
		}

		// Verify the result
		if resultFloat != expectedResult {
			t.Errorf("DividedBy failed: %f / %f = %f, expected %f", floatA, floatB, resultFloat, expectedResult)
		}

		// Verify the original object is not altered
		if floatObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, floatObj.Val())
		}
	}
}

func TestFloatSquareRoot(t *testing.T) {
	for i := 0; i < 500; i++ {
		// Generate random float as string
		a := randomFloat()

		// Parse it to float64 for verification
		floatA, err := floatFromString(a)
		if err != nil {
			t.Fatalf("Failed to convert string to float: %v", err)
		}

		// Ensure the number is non-negative
		if floatA < 0 {
			floatA = -floatA
			a = strconv.FormatFloat(floatA, 'f', -1, 64)
		}

		// Expected result using standard math
		expectedResult := new(big.Float).SetFloat64(floatA).Sqrt(nil)

		// Perform the SquareRoot operation
		floatObj := NewFloat(a)
		originalValue := floatObj.Val()
		result := floatObj.SquareRoot()

		// Convert result to big.Float for comparison
		resultFloat, _, err := big.NewFloat(0).Parse(result.Val(), 10)
		if err != nil {
			t.Fatalf("Failed to convert result string to big.Float: %v", err)
		}

		// Verify the result
		if resultFloat.Cmp(expectedResult) != 0 {
			t.Errorf("SquareRoot failed: sqrt(%f) = %s, expected %s", floatA, result.Val(), expectedResult.String())
		}

		// Verify the original object is not altered
		if floatObj.Val() != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, floatObj.Val())
		}
	}
}

// func TestFloatToThePowerOf(t *testing.T) {
// 	for i := 0; i < 500; i++ {
// 		// Generate random floats as strings for base and exponent
// 		base := randomFloat()
// 		exponent := randomFloat()

// 		// Parse them to float64 for verification
// 		floatBase, err := floatFromString(base)
// 		if err != nil {
// 			t.Fatalf("Failed to convert string to float: %v", err)
// 		}
// 		floatExponent, err := floatFromString(exponent)
// 		if err != nil {
// 			t.Fatalf("Failed to convert string to float: %v", err)
// 		}

// 		// Skip negative base cases (these should panic)
// 		if floatBase < 0 {
// 			continue
// 		}

// 		// Expected result using standard math
// 		expectedResult := new(big.Float).SetFloat64(floatBase).SetPrec(1024)
// 		expectedResult.Exp(big.NewFloat(floatExponent), nil)

// 		// Perform the ToThePowerOf operation
// 		floatObj := NewFloat(base)
// 		originalValue := floatObj.Val()
// 		result := floatObj.ToThePowerOf(exponent)

// 		// Convert result to big.Float for comparison
// 		resultFloat, _, err := big.NewFloat(0).Parse(result.Val(), 10)
// 		if err != nil {
// 			t.Fatalf("Failed to convert result string to big.Float: %v", err)
// 		}

// 		// Verify the result
// 		if resultFloat.Cmp(expectedResult) != 0 {
// 			t.Errorf("ToThePowerOf failed: %f ^ %f = %s, expected %s", floatBase, floatExponent, result.Val(), expectedResult.String())
// 		}

// 		// Verify the original object is not altered
// 		if floatObj.Val() != originalValue {
// 			t.Errorf("Original object value altered: expected %s, got %s", originalValue, floatObj.Val())
// 		}
// 	}
// }
