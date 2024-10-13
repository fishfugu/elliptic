package bigarith

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

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
