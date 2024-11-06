package bigarith

import (
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// Define a struct to hold the test case inputs and expected output
type testFloatCase struct {
	description string
	expected    Float
	got         Float
}

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

func isEqualUpToLastDigit(a, b Float) bool {
	aStr := a.truncDecPoints().Val()
	bStr := b.truncDecPoints().Val()
	if aStr[:len(aStr)-1] == bStr[:len(bStr)-1] {
		return true
	}
	return false
}

func TestFloatPi(t *testing.T) {
	// Pi to 510 dec places
	expectedPi := "3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951941511609433057270365759591953092186117381932611793105118548074462379962749567351885752724891227938183011949129833673362"
	truncExpectedPi := NewFloat(expectedPi).truncDecPoints()
	if !isEqualUpToLastDigit(NewFloat("1").TimesPi(), truncExpectedPi) {
		t.Errorf("Pi failed. Expected: '%s', got: '%s'", truncExpectedPi, NewFloat("1").TimesPi())
	}
}

func TestFloatMod(t *testing.T) {
	testFloatCases := []testFloatCase{
		{
			description: "simple - 0 mod anything",
			expected:    NewFloat("0"),
			got:         NewFloat("0").Mod("2"),
		},
		// {
		// 	description: "2 mod 2",
		// 	expected:    NewFloat("0"),
		// 	got:         NewFloat("2").Mod("2"),
		// },
		// {
		// 	description: "4.5 mod 3",
		// 	expected:    NewFloat("1.5"),
		// 	got:         NewFloat("4.5").Mod("3"),
		// },
		// {
		// 	description: "3 times 1.234 plus 0.5678 mod 1.234",
		// 	expected:    NewFloat("0.5678"),
		// 	got:         NewFloat("1.234").Times("3").Plus("0.5678").Mod("1.234"),
		// },
		// {
		// 	description: "4 1/2 pi mod 2 pi",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("2"),
		// 	got:         NewFloat("4").TimesPi().Plus(NewFloat("1").TimesPi().DividedBy("2").Val()).Mod(NewFloat("2").TimesPi().Val()),
		// },
	}
	for _, tc := range testFloatCases {
		if !isEqualUpToLastDigit(tc.got, tc.expected) {
			t.Errorf("Mod(%s) failed. Expected: '%s', got: '%s'", tc.description, tc.expected, tc.got.Val())
		}
	}
}

func TestArcTan(t *testing.T) {
	// Define a list of test cases
	testFloatCases := []testFloatCase{
		{
			description: "0",
			expected:    NewFloat("0"),
			got:         NewFloat("0").ArcTan(),
		},
		// {
		// 	description: "1",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("4").truncDecPoints(),
		// 	got:         NewFloat("1").ArcTan(),
		// },
		// {
		// 	description: "-1",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("-4").truncDecPoints(),
		// 	got:         NewFloat("-1").ArcTan(),
		// },
		// {
		// 	description: "1/sqrt{3}",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("6").truncDecPoints(),
		// 	got:         NewFloat("1").DividedBy(NewFloat("3").SquareRoot().Val()).ArcTan(),
		// },
		// {
		// 	description: "-1/sqrt{3}",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("-6").truncDecPoints(),
		// 	got:         NewFloat("-1").DividedBy(NewFloat("3").SquareRoot().Val()).ArcTan(),
		// },
		// {
		// 	description: "sqrt{3}",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("3").truncDecPoints(),
		// 	got:         NewFloat("3").SquareRoot().ArcTan(),
		// },
		// {
		// 	description: "-sqrt{3}",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("-3").truncDecPoints(),
		// 	got:         NewFloat("3").SquareRoot().Neg().ArcTan(),
		// },
	}
	for _, tc := range testFloatCases {
		if !isEqualUpToLastDigit(tc.got, tc.expected) {
			t.Errorf("ArcTan(%s) failed. Expected: '%s', got: '%s'", tc.description, tc.expected, tc.got.Val())
		}
	}
}

func TestArcSin(t *testing.T) {
	// Define a list of test cases
	testFloatCases := []testFloatCase{
		{
			description: "1",
			expected:    NewFloat("1").TimesPi().DividedBy("2").truncDecPoints(),
			got:         NewFloat("1").ArcSin(),
		},
		// {
		// 	description: "1/2",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("6").truncDecPoints(),
		// 	got:         NewFloat("1").DividedBy("2").ArcSin(),
		// },
		// {
		// 	description: "0",
		// 	expected:    NewFloat("0"),
		// 	got:         NewFloat("0").ArcSin(),
		// },
		// {
		// 	description: "-1/2",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("-6").truncDecPoints(),
		// 	got:         NewFloat("1").DividedBy("-2").ArcSin(),
		// },
		// {
		// 	description: "-1",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("-2").truncDecPoints(),
		// 	got:         NewFloat("-1").ArcSin(),
		// },
	}
	for _, tc := range testFloatCases {
		if !isEqualUpToLastDigit(tc.got, tc.expected) {
			t.Errorf("ArcSin(%s) failed. Expected: '%s', got: '%s'", tc.description, tc.expected, tc.got.Val())
		}
	}
}

func TestArcCos(t *testing.T) {
	// Define a list of test cases
	testFloatCases := []testFloatCase{
		{
			description: "1",
			expected:    NewFloat("0"),
			got:         NewFloat("1").ArcCos(),
		},
		// {
		// 	description: "1/2",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("3").truncDecPoints(),
		// 	got:         NewFloat("1").DividedBy("2").ArcCos(),
		// },
		// {
		// 	description: "0",
		// 	expected:    NewFloat("1").TimesPi().DividedBy("2").truncDecPoints(),
		// 	got:         NewFloat("0").ArcCos(),
		// },
		// {
		// 	description: "-1/2",
		// 	expected:    NewFloat("1").TimesPi().Times("2").DividedBy("3").truncDecPoints(),
		// 	got:         NewFloat("1").DividedBy("-2").ArcCos(),
		// },
		// {
		// 	description: "-1",
		// 	expected:    NewFloat("1").TimesPi().truncDecPoints(),
		// 	got:         NewFloat("-1").ArcCos(),
		// },
	}
	for _, tc := range testFloatCases {
		if !isEqualUpToLastDigit(tc.got, tc.expected) {
			t.Errorf("ArcCos(%s) failed. Expected: '%s', got: '%s'", tc.description, tc.expected, tc.got.Val())
		}
	}
}

func TestSin(t *testing.T) {
	// Define a list of test cases
	testFloatCases := []testFloatCase{
		{
			description: "0",
			expected:    NewFloat("0"),
			got:         NewFloat("0").Sin(),
		},
		// 	{
		// 		description: "pi/12",
		// 		expected:    NewFloat("6").SquareRoot().Minus(NewFloat("2").SquareRoot().Val()).DividedBy("4").truncDecPoints(),
		// 		got:         NewFloat("1").TimesPi().DividedBy("12").Sin(),
		// 	},
		// 	{
		// 		description: "pi/6",
		// 		expected:    NewFloat("1").DividedBy("2").truncDecPoints(),
		// 		got:         NewFloat("1").TimesPi().DividedBy("6").Sin(),
		// 	},
		// 	{
		// 		description: "pi/4",
		// 		expected:    NewFloat("2").SquareRoot().DividedBy("2").truncDecPoints(),
		// 		got:         NewFloat("1").TimesPi().DividedBy("4").Sin(),
		// 	},
		// 	{
		// 		description: "pi/3",
		// 		expected:    NewFloat("3").SquareRoot().DividedBy("2").truncDecPoints(),
		// 		got:         NewFloat("1").TimesPi().DividedBy("3").Sin(),
		// 	},
		// 	{
		// 		description: "5 pi / 12",
		// 		expected:    NewFloat("6").SquareRoot().Plus(NewFloat("2").SquareRoot().Val()).DividedBy("4").truncDecPoints(),
		// 		got:         NewFloat("1").TimesPi().Times("5").DividedBy("12").Sin(),
		// 	},
		// 	{
		// 		description: "pi/2",
		// 		expected:    NewFloat("1"),
		// 		got:         NewFloat("1").TimesPi().DividedBy("2").Sin(),
		// 	},
		// 	{
		// 		description: "2pi + pi/12",
		// 		expected:    NewFloat("6").SquareRoot().Minus(NewFloat("2").SquareRoot().Val()).DividedBy("4").truncDecPoints(),
		// 		got:         NewFloat("1").TimesPi().DividedBy("12").Plus(NewFloat("2").TimesPi().Val()).Sin(),
		// 	},
	}
	for _, tc := range testFloatCases {
		if !isEqualUpToLastDigit(tc.got, tc.expected) {
			t.Errorf("Sin(%s) failed. Expected: '%s', got: '%s'", tc.description, tc.expected, tc.got.Val())
		}
	}
}

func TestCos(t *testing.T) {
	// Define a list of test cases
	testFloatCases := []testFloatCase{
		{
			description: "0",
			expected:    NewFloat("1"),
			got:         NewFloat("0").Cos(),
		},
		// {
		// 	description: "pi/12",
		// 	expected:    NewFloat("6").SquareRoot().Plus(NewFloat("2").SquareRoot().Val()).DividedBy("4").truncDecPoints(),
		// 	got:         NewFloat("1").TimesPi().DividedBy("12").Cos(),
		// },
		// {
		// 	description: "pi/6",
		// 	expected:    NewFloat("3").SquareRoot().DividedBy("2").truncDecPoints(),
		// 	got:         NewFloat("1").TimesPi().DividedBy("6").Cos(),
		// },
		// {
		// 	description: "pi/4",
		// 	expected:    NewFloat("2").SquareRoot().DividedBy("2").truncDecPoints(),
		// 	got:         NewFloat("1").TimesPi().DividedBy("4").Cos(),
		// },
		// {
		// 	description: "pi/3",
		// 	expected:    NewFloat("1").DividedBy("2").truncDecPoints(),
		// 	got:         NewFloat("1").TimesPi().DividedBy("3").Cos(),
		// },
		// {
		// 	description: "5 pi / 12",
		// 	expected:    NewFloat("6").SquareRoot().Minus(NewFloat("2").SquareRoot().Val()).DividedBy("4").truncDecPoints(),
		// 	got:         NewFloat("1").TimesPi().Times("5").DividedBy("12").Cos(),
		// },
		// {
		// 	description: "pi/2",
		// 	expected:    NewFloat("0"),
		// 	got:         NewFloat("1").TimesPi().DividedBy("2").Cos(),
		// },
		// {
		// 	description: "2pi + pi/12",
		// 	expected:    NewFloat("6").SquareRoot().Plus(NewFloat("2").SquareRoot().Val()).DividedBy("4").truncDecPoints(),
		// 	got:         NewFloat("1").TimesPi().DividedBy("12").Plus(NewFloat("2").TimesPi().Val()).Cos(),
		// },
	}
	for _, tc := range testFloatCases {
		if !isEqualUpToLastDigit(tc.got, tc.expected) {
			t.Errorf("Cos(%s) failed. Expected: '%s', got: '%s'", tc.description, tc.expected, tc.got.Val())
		}
	}
}

func TestFloatPlus(t *testing.T) {
	for i := 0; i < 10; i++ {
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
	for i := 0; i < 10; i++ {
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
	for i := 0; i < 10; i++ {
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
	for i := 0; i < 10; i++ {
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
	for i := 0; i < 10; i++ {
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
		expectedResult := new(big.Float).Sqrt(new(big.Float).SetFloat64(floatA))

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
