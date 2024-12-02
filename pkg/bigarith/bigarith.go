package bigarith

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// min value
const numberOfDecimalPoints = int(50) // from 48 ... 505 - eveything else seems to error
const precision = uint(numberOfDecimalPoints * 17 / 5)

// Max iterations or precision tolerance
var toleranceRationalStr = fmt.Sprintf("1/1%s", strings.Repeat("0", numberOfDecimalPoints))
var toleranceRational = NewRational(toleranceRationalStr)
var testToleranceRationalStr = fmt.Sprintf("1/1%s", strings.Repeat("0", numberOfDecimalPoints-1))
var testToleranceRational = NewRational(testToleranceRationalStr)

var toleranceFloat = new(Float).setRational(toleranceRational)
var testToleranceFloat = new(Float).setRational(testToleranceRational)

func init() {
	// log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// output to stdout instead of the default stderr.
	logrus.SetOutput(os.Stdout)

	// log the debug severity or above.
	logrus.SetLevel(logrus.ErrorLevel)
}

// Error refactor functions

// bigFloat converts the string representation of the bigarith.Float into a *big.Float
// An internal only function to simplify other code
// NOTE: Panics on bad string representation of big.Float
func bigFloat(a string) *big.Float {
	// in order to create a simplified application of string <-> big.*
	// no errors are handed back up the chain
	// instead simply Panic here with clear reasons why
	x, ok := new(big.Float).SetPrec(precision).SetString(a)
	if !ok {
		errA := a
		if len(errA) > 250 {
			errA = errA[:250] + "..."
		}
		panic(fmt.Sprintf("Error creating big.Float using SetString - base = %d - string value = %s", 10, errA))
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

// bigRational converts the string representation of the bigarith.Rational into a *big.Rational
// An internal only function to simplify other code
// NOTE: Panics on bad string representation of big.Rational
func bigRational(a string) *big.Rat {
	// in order to create a simplified application of string <-> big.*
	// no errors are handed back up the chain
	// instead simply Panic here with clear reasons why
	x, ok := new(big.Rat).SetString(a)
	if !ok {
		panic(fmt.Sprintf("error setting Big Rat: %s", a))
	}
	return x
}

// Logarithm approximation using Newton's method
// internal function only - just works in big.Float
func logNewton(x *big.Float) *big.Float {
	logrus.Debugf("logNewton, x = %s", x.String())
	if x.Cmp(big.NewFloat(0)) <= 0 {
		panic("Logarithm only defined for positive numbers")
	}

	// Initial guess (logarithm approximation)
	approx := bigFloat("0.5")

	// Newton's method to refine the logarithm estimate
	two := bigFloat("2")
	tmp := bigFloat("0")
	for i := 0; i < 1024; i++ { // Don't change lightly - chosen through experimentation on max accuracy / min time
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

	tmp := bigFloat("0")       // Temporary variable for addition step
	for i := 1; i < 300; i++ { // Don't change lightly - chosen through experimentation on max accuracy / min time
		factorial = factorial.Mul(factorial, bigFloat(strconv.Itoa(i)))
		term = term.Mul(term, x)
		tmp.Quo(term, factorial)
		result.Add(result, tmp)
	}

	return result
}

// nthRoot calculates the integer nth root of a big.Int.
// It returns the closest integer approximation of the root.
func nthRoot(x *big.Int, n *big.Int) (*big.Int, error) {
	if n.Cmp(big.NewInt(1)) == 0 { // The 1st root is the number itself
		return x, nil
	}

	// Binary search for nth root
	low := big.NewInt(0)
	high := new(big.Int).Set(x)
	mid := new(big.Int)
	one := big.NewInt(1)

	for low.Cmp(high) <= 0 {
		mid.Add(low, high).Div(mid, big.NewInt(2))

		// mid^n
		midPow := new(big.Int).Exp(mid, n, nil)

		cmp := midPow.Cmp(x)
		if cmp == 0 {
			return mid, nil // Found exact root
		} else if cmp < 0 {
			low.Add(mid, one)
		} else {
			high.Sub(mid, one)
		}
	}

	return high, nil // Return closest approximation
}
