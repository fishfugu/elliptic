package bigarith

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// Define a struct to hold the test case inputs and expected output
type testRationalCase struct {
	description string
	expected    Rational
	got         Rational
}

// [ ] TODO: Test - NthRoot
// [ ] TODO: Test - Mod2Pi
// [ ] TODO: Test - ToThePowerOf

const testLoopNumber = 10

// RandomRational returns a random Rational between 0 and 1
func RandomRational() Rational {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a random Rational between 0 and 1
	fracDenom := randomSmallInt(1000000000)                                       // choose a denom between 0 and 1,000,000,000
	fracDenomInt, _ := strconv.Atoi(fracDenom)                                    // convert to int
	fracNum := NewInt(randomSmallInt(int64(fracDenomInt))).Plus(fracDenom).strVal // choose a number even dtributed between 0 and the denom

	return NewRational(fmt.Sprintf("%s/%s", fracNum, fracDenom))
}

// RandomRationalBetween returns a random Rational between the Rational representations of strings a and b
func RandomRationalBetween(aStr, bStr string) Rational {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Calculate ((|b - a|) * f) + a -- where f is a random fraction 0 <= f <= 1
	return NewRational(bStr).Diff(aStr).Times(RandomRational().strVal).Plus(aStr)
}

func TestPi(t *testing.T) {
	// Pi to 510 dec places
	piExpectedString := "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989"
	piExpectedFloat := NewFloat("1").Times(piExpectedString)
	piGotRational := NewRational("1").TimesPi()
	piGoFloat := new(Float).setRational(piGotRational)
	if piExpectedFloat.Diff(piGoFloat.strVal).Compare(testToleranceFloat.strVal) > 0 {
		t.Errorf(`
rational Pi estimate failed.

expected:
%s

got:
%s

diff:
%s

expected diff:
%s

`,
			piExpectedFloat.strVal,
			piGoFloat.strVal,
			piExpectedFloat.Diff(piGoFloat.strVal).strVal,
			testToleranceFloat.strVal,
		)
	}
}

func TestNthRoot(t *testing.T) {
	testRationalCases := []testRationalCase{
		{
			description: "Square Root of 4",
			expected:    NewRational("2"),
			got:         NewRational("4").NthRoot("2"),
		},
		{
			description: "Cube Root of 27",
			expected:    NewRational("3"),
			got:         NewRational("27").NthRoot("3"),
		},
	}
	for _, tc := range testRationalCases {
		if tc.expected.Diff(tc.got.strVal).Compare(testToleranceFloat.strVal) > 0 {
			t.Errorf(`
rational Square Root estimate failed.

expected:
%s

got:
%s

diff:
%s

expected diff:
%s

`,
				tc.expected.strVal,
				tc.got.strVal,
				tc.expected.Diff(tc.got.strVal).strVal,
				testToleranceFloat.strVal,
			)
		}
	}
}

// func TestRationalMod(t *testing.T) {
// 	// TODO: turn this into random, like the ones below
// 	testRationalCases := []testRationalCase{
// 		{
// 			description: "simple - 0 mod anything",
// 			expected:    NewRational("0"),
// 			got:         NewRational("0").Mod("2"),
// 		},
// 		{
// 			description: "2 mod 2",
// 			expected:    NewRational("0"),
// 			got:         NewRational("2").Mod("2"),
// 		},
// 		{
// 			description: "9/2 mod 3", // answer should be "3/2"
// 			expected:    NewRational("3/2"),
// 			got:         NewRational("9/2").Mod("3"),
// 		},
// 		{
// 			description: "3/1 times 1234/1000 plus 5678/10000 mod 1234/1000", // answer should be 5678/10000 (or its equiv)
// 			expected:    NewRational("5678/10000"),
// 			got:         NewRational("1234/1000").Times("3").Plus("5678/10000").Mod("1234/1000"),
// 		},
// 		{
// 			description: "9/2 pi mod 2 pi",
// 			expected:    NewRational("1/2").TimesPi(),
// 			got:         NewRational("9/2").TimesPi().Mod(NewRational("2").TimesPi().strVal),
// 		},
// 	}
// 	for _, tc := range testRationalCases {
// 		if tc.expected.Diff(tc.got.strVal).Compare(testToleranceRational.strVal) > 0 {
// 			t.Errorf(`
// Test for: %s -- failed.

// expected:
// %s, over
// %s,

// got:
// %s, over
// %s

// diff:
// %s, over
// %s,

// expected diff:
// %s, over
// %s

// `,
// 				tc.description,
// 				tc.expected.Num().strVal,
// 				tc.expected.Denom().strVal,
// 				tc.got.Num().strVal,
// 				tc.got.Denom().strVal,
// 				tc.expected.Diff(tc.got.strVal).Num().strVal,
// 				tc.expected.Diff(tc.got.strVal).Denom().strVal,
// 				testToleranceRational.Num().strVal,
// 				testToleranceRational.Denom().strVal)
// 		}
// 	}
// }

// func TestRationalArcTan(t *testing.T) {
// 	// Define a list of test cases
// 	// TODO: turn this into random, like the ones below
// 	testRationalCases := []testRationalCase{
// 		// {
// 		// 	description: "ArcTan(0) = 0",
// 		// 	expected:    NewRational("0"),
// 		// 	got:         NewRational("0").ArcTan(),
// 		// },
// 		{
// 			description: "ArcTan(1) = 1/4 Pi",
// 			expected:    NewRational("1/4").TimesPi(),
// 			got:         NewRational("1").arcTan(),
// 		},
// 		// {
// 		// 	description: "ArcTan(-1) = -1/4 Pi",
// 		// 	expected:    NewRational("-1/4").TimesPi(),
// 		// 	got:         NewRational("-1").ArcTan(),
// 		// },
// 		// {
// 		// 	description: "1/sqrt{3}",
// 		// 	expected:    NewRational("1/6").TimesPi(),
// 		// 	got:         NewRational("1").DividedBy(NewRational("3").SquareRoot().strVal).ArcTan(),
// 		// },
// 		// {
// 		// 	description: "-1/sqrt{3}",
// 		// 	expected:    NewRational("-1/6").TimesPi(),
// 		// 	got:         NewRational("-1").DividedBy(NewRational("3").SquareRoot().strVal).ArcTan(),
// 		// },
// 		// {
// 		// 	description: "sqrt{3}",
// 		// 	expected:    NewRational("1/3").TimesPi(),
// 		// 	got:         NewRational("3").SquareRoot().ArcTan(),
// 		// },
// 		// {
// 		// 	description: "-sqrt{3}",
// 		// 	expected:    NewRational("-1/3").TimesPi(),
// 		// 	got:         NewRational("3").SquareRoot().Neg().ArcTan(),
// 		// },
// 	}
// 	for _, tc := range testRationalCases {
// 		if tc.expected.Diff(tc.got.strVal).Compare(testToleranceRational.strVal) > 0 {
// 			t.Errorf(`
// Test for: %s -- failed.
// expected:
// %s, over
// %s,

// got:
// %s, over
// %s

// diff:
// %s, over
// %s,

// expected diff:
// %s, over
// %s

// `,
// 				tc.description,
// 				tc.expected.Num().strVal,
// 				tc.expected.Denom().strVal,
// 				tc.got.Num().strVal,
// 				tc.got.Denom().strVal,
// 				tc.expected.Diff(tc.got.strVal).Num().strVal,
// 				tc.expected.Diff(tc.got.strVal).Denom().strVal,
// 				testToleranceRational.Num().strVal,
// 				testToleranceRational.Denom().strVal)
// 		}
// 	}
// }

// func TestRationalArcCos(t *testing.T) {
// 	// Define a list of test cases
// 	// TODO: turn this into random, like the ones below
// 	testRationalCases := []testRationalCase{
// 		// {
// 		// 	description: "1",
// 		// 	expected:    NewRational("0"),
// 		// 	got:         NewRational("1").ArcCos(),
// 		// },
// 		// {
// 		// 	description: "1/2",
// 		// 	expected:    NewRational("1/3").TimesPi(),
// 		// 	got:         NewRational("1/2").ArcCos(),
// 		// },
// 		// {
// 		// 	description: "0",
// 		// 	expected:    NewRational("1/2").TimesPi(),
// 		// 	got:         NewRational("0").ArcCos(),
// 		// },
// 		// {
// 		// 	description: "-1/2",
// 		// 	expected:    NewRational("2/3").TimesPi(),
// 		// 	got:         NewRational("-1/2").ArcCos(),
// 		// },
// 		// {
// 		// 	description: "-1",
// 		// 	expected:    NewRational("1").TimesPi(),
// 		// 	got:         NewRational("-1").ArcCos(),
// 		// },
// 	}
// 	for _, tc := range testRationalCases {
// 		if tc.expected.Diff(tc.got.strVal).Compare(testToleranceRational.strVal) > 0 {
// 			t.Errorf(`
// Test for: %s -- failed.

// expected:
// %s, over
// %s,

// got:
// %s, over
// %s

// diff:
// %s, over
// %s,

// diff float,
// %s,

// expected diff:
// %s, over
// %s

// `,
// 				tc.description,
// 				tc.expected.Num().strVal,
// 				tc.expected.Denom().strVal,
// 				tc.got.Num().strVal,
// 				tc.got.Denom().strVal,
// 				tc.expected.Diff(tc.got.strVal).Num().strVal,
// 				tc.expected.Diff(tc.got.strVal).Denom().strVal,
// 				NewFloat(tc.expected.Diff(tc.got.strVal).Num().strVal).DividedBy(tc.expected.Diff(tc.got.strVal).Denom().strVal).strVal,
// 				testToleranceRational.Num().strVal,
// 				testToleranceRational.Denom().strVal)
// 		}
// 	}
// }

// NOTE: this only tests to 15 dec. places
// because the big/math library doesn't have Cos and 15 chars seems about the limit
// for big/math
func TestRationalCos(t *testing.T) {
	startTime := time.Now() // Track the start time of the test function
	var totalTime time.Duration

	for i := 0; i < testLoopNumber; i++ {
		a := RandomRationalBetween("0", NewRational("2").TimesPi().strVal)
		aNum := a.Num().strVal
		aDenom := a.Denom().strVal

		// Time the operation
		startFuncTime := time.Now()
		result := a.Cos()
		funcDuration := time.Since(startFuncTime)
		totalTime += funcDuration // Accumulate time spent

		// Get test values
		originalValue := a.strVal

		// float check values
		checkNumBasicFloat, _ := strconv.ParseFloat(aNum, 64)
		checkDenomBasicFloat, _ := strconv.ParseFloat(aDenom, 64)
		checkResultBasicFloat := math.Cos(checkNumBasicFloat / checkDenomBasicFloat)
		checkResultStr := fmt.Sprintf("%0.50f", checkResultBasicFloat)

		// Verify the result
		if checkResultStr[:15] == testToleranceFloat.strVal[:15] {
			t.Errorf("Cos() failed: Cos(%s) = %s = %s, expected %s", a.strVal, "NOT result.strVal", result.bigRatVal.FloatString(numberOfDecimalPoints), checkResultStr)
		}

		// Verify that the internal value of the original object hasn't been altered
		if a.strVal != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, a.strVal)
		}
	}
	// Calculate the average time taken by the SquareRoot method
	averageTime := totalTime / testLoopNumber

	// Report total and average time using t.Logf
	t.Logf("Total test duration: %s", time.Since(startTime))
	t.Logf("Average Cos operation time: %s", averageTime)
}

func TestRationalRationalPlus(t *testing.T) {
	startTime := time.Now() // Track the start time of the test function
	var totalTime time.Duration

	for i := 0; i < testLoopNumber; i++ {
		// Generate random floats as strings
		aNum := randomInt()
		aDenom := randomInt()
		bNum := randomInt()
		bDenom := randomInt()

		a := NewRational(fmt.Sprintf("%s/%s", aNum, aDenom))
		b := NewRational(fmt.Sprintf("%s/%s", bNum, bDenom))

		// Time the operation
		startFuncTime := time.Now()
		result := a.Plus(b.strVal)
		funcDuration := time.Since(startFuncTime)
		totalTime += funcDuration // Accumulate time spent

		// Get test values
		originalValue := a.strVal
		resultFloat := new(Float).setRational(result)

		// Float check values
		checkAFloat := NewFloat(aNum).DividedBy(aDenom)
		checkBFloat := NewFloat(bNum).DividedBy(bDenom)
		checkResultFloat := checkAFloat.Plus(checkBFloat.strVal)

		// Verify the result
		if resultFloat.Diff(checkResultFloat.strVal).Compare(testToleranceFloat.strVal) > 0 {
			t.Errorf("Plus failed: %s + %s = %s, expected %s", a.strVal, b.strVal, result.strVal, checkResultFloat.strVal)
		}

		// Verify that the internal value of the original object hasn't been altered
		if a.strVal != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, a.strVal)
		}
	}
	// Calculate the average time taken by the SquareRoot method
	averageTime := totalTime / testLoopNumber

	// Report total and average time using t.Logf
	t.Logf("Total test duration: %s", time.Since(startTime))
	t.Logf("Average Plus operation time: %s", averageTime)
}

func TestRationalMinus(t *testing.T) {
	startTime := time.Now() // Track the start time of the test function
	var totalTime time.Duration

	for i := 0; i < testLoopNumber; i++ {
		// Generate random floats as strings
		aNum := randomInt()
		aDenom := randomInt()
		bNum := randomInt()
		bDenom := randomInt()

		a := NewRational(fmt.Sprintf("%s/%s", aNum, aDenom))
		b := NewRational(fmt.Sprintf("%s/%s", bNum, bDenom))

		// Time the operation
		startFuncTime := time.Now()
		result := a.Minus(b.strVal)
		funcDuration := time.Since(startFuncTime)
		totalTime += funcDuration // Accumulate time spent

		// Get test values
		originalValue := a.strVal
		resultFloat := new(Float).setRational(result)

		// Float check values
		checkAFloat := NewFloat(aNum).DividedBy(aDenom)
		checkBFloat := NewFloat(bNum).DividedBy(bDenom)
		checkResultFloat := checkAFloat.Minus(checkBFloat.strVal)

		// Verify the result
		if resultFloat.Diff(checkResultFloat.strVal).Compare(testToleranceFloat.strVal) > 0 {
			t.Errorf("Minus failed: %s - %s = %s, expected %s", a.strVal, b.strVal, result.strVal, checkResultFloat.strVal)
		}

		// Verify that the internal value of the original object hasn't been altered
		if a.strVal != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, a.strVal)
		}
	}
	// Calculate the average time taken by the SquareRoot method
	averageTime := totalTime / testLoopNumber

	// Report total and average time using t.Logf
	t.Logf("Total test duration: %s", time.Since(startTime))
	t.Logf("Average Minus operation time: %s", averageTime)
}

func TestRationalTimes(t *testing.T) {
	startTime := time.Now() // Track the start time of the test function
	var totalTime time.Duration

	for i := 0; i < testLoopNumber; i++ {
		// Generate random floats as strings
		aNum := randomInt()
		aDenom := randomInt()
		bNum := randomInt()
		bDenom := randomInt()

		a := NewRational(fmt.Sprintf("%s/%s", aNum, aDenom))
		b := NewRational(fmt.Sprintf("%s/%s", bNum, bDenom))

		// Time the operation
		startFuncTime := time.Now()
		result := a.Times(b.strVal)
		funcDuration := time.Since(startFuncTime)
		totalTime += funcDuration // Accumulate time spent

		// Get test values
		originalValue := a.strVal
		resultFloat := new(Float).setRational(result)

		// Float check values
		checkAFloat := NewFloat(aNum).DividedBy(aDenom)
		checkBFloat := NewFloat(bNum).DividedBy(bDenom)
		checkResultFloat := checkAFloat.Times(checkBFloat.strVal)

		// Verify the result
		if resultFloat.Diff(checkResultFloat.strVal).Compare(testToleranceFloat.strVal) > 0 {
			t.Errorf("Times failed: %s * %s = %s, expected %s", a.strVal, b.strVal, result.strVal, checkResultFloat.strVal)
		}

		// Verify that the internal value of the original object hasn't been altered
		if a.strVal != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, a.strVal)
		}
	}
	// Calculate the average time taken by the SquareRoot method
	averageTime := totalTime / testLoopNumber

	// Report total and average time using t.Logf
	t.Logf("Total test duration: %s", time.Since(startTime))
	t.Logf("Average Times operation time: %s", averageTime)
}

func TestRationalDividedBy(t *testing.T) {
	startTime := time.Now() // Track the start time of the test function
	var totalTime time.Duration

	for i := 0; i < testLoopNumber; i++ {
		// Generate random floats as strings
		aNum := randomInt()
		aDenom := randomInt()
		bNum := randomInt()
		bDenom := randomInt()

		a := NewRational(fmt.Sprintf("%s/%s", aNum, aDenom))
		b := NewRational(fmt.Sprintf("%s/%s", bNum, bDenom))

		// Time the operation
		startFuncTime := time.Now()
		result := a.DividedBy(b.strVal)
		funcDuration := time.Since(startFuncTime)
		totalTime += funcDuration // Accumulate time spent

		// Get test values
		originalValue := a.strVal
		resultFloat := new(Float).setRational(result)

		// Float check values
		checkAFloat := NewFloat(aNum).DividedBy(aDenom)
		checkBFloat := NewFloat(bNum).DividedBy(bDenom)
		checkResultFloat := checkAFloat.DividedBy(checkBFloat.strVal)

		// Verify the result
		if resultFloat.Diff(checkResultFloat.strVal).Compare(testToleranceFloat.strVal) > 0 {
			t.Errorf("DividedBy failed: %s / %s = %s, expected %s", a.strVal, b.strVal, result.strVal, checkResultFloat.strVal)
		}

		// Verify that the internal value of the original object hasn't been altered
		if a.strVal != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, a.strVal)
		}
	}

	// Calculate the average time taken by the SquareRoot method
	averageTime := totalTime / testLoopNumber

	// Report total and average time using t.Logf
	t.Logf("Total test duration: %s", time.Since(startTime))
	t.Logf("Average DividedBy operation time: %s", averageTime)
}

func TestRationalSquareRoot(t *testing.T) {
	startTime := time.Now() // Track the start time of the test function
	var totalTime time.Duration

	for i := 0; i < testLoopNumber; i++ {
		// Generate random floats as strings
		aNum := randomInt()
		aDenom := randomInt()

		a := NewRational(fmt.Sprintf("%s/%s", aNum, aDenom))

		// Time the operation
		startFuncTime := time.Now()
		result := a.SquareRoot()
		funcDuration := time.Since(startFuncTime)
		totalTime += funcDuration // Accumulate time spent

		// Get test values
		originalValue := a.strVal
		resultFloat := new(Float).setRational(result)

		// Float check values
		checkAFloat := NewFloat(aNum).DividedBy(aDenom)
		checkResultFloat := checkAFloat.SquareRoot()

		// Verify the result
		if resultFloat.Diff(checkResultFloat.strVal).Compare(testToleranceFloat.strVal) > 0 {
			t.Errorf("Sqrt failed: sqrt(%s) = %s, expected %s", a.strVal, result.strVal, checkResultFloat.strVal)
		}

		// Verify that the internal value of the original object hasn't been altered
		if a.strVal != originalValue {
			t.Errorf("Original object value altered: expected %s, got %s", originalValue, a.strVal)
		}
	}

	// Calculate the average time taken by the SquareRoot method
	averageTime := totalTime / testLoopNumber

	t.Logf("Total test duration: %s", time.Since(startTime))
	t.Logf("Average SquareRoot operation time: %s", averageTime)
}
