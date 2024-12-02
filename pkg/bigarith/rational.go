package bigarith

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

// [ ] TODO: make everything a Rational - go through and take out Int/Float funcs and see what errors
// [ ] TODO: go back to efficiency / performacne testing (dropping it for now...)
// [ ] TODO: take out anything not used outside this bigarith package (e.g. don't need all the trig funcs)

// DEFINITION STUFF

// Rational is a bigarith type Rational Arithmetic
type Rational struct {
	strVal    string   // the string representing the rational value
	bigRatVal *big.Rat // cache the parsed big.Rat
}

func NewRational(a string) Rational {
	return new(Rational).set(a)
}

// DISCOVERY FUNCS
// the functions that surface info about the Rational - don't return a Rational itself

// Val returns the string representation of the bigarith.Rational object
func (r Rational) Val() string {
	return r.truncate().strVal
}

// Num returns the Int representation of the bigarith.Rational numerator
func (r Rational) Num() Int {
	return NewInt(r.bigRatVal.Num().String()) // Use cached bigRatVal
}

// Denom returns the Int representation of the bigarith.Rational denominator
func (r Rational) Denom() Int {
	return NewInt(r.bigRatVal.Denom().String()) // Use cached bigRatVal
}

// Compare takes a string representation of a rational number and compares
// it with the current bigarith.Rational,
// returns:
// -1 if x <  y, 0 if x == y (incl. -0 == 0, -Inf == -Inf, and +Inf == +Inf), +1 if x >  y
func (r Rational) Compare(a string) int {
	return r.bigRatVal.Cmp(bigRational(a)) // Use cached bigRatVal
}

// MANIPULATION FUNCS
// functions that manipulate the value of a Rational
// should return the new value without modifying the original Rational

// NON-EXPORTED (internal only)

// set sets the string value that represents the bigarith.Rational value
// and returns a new Rational with the updated value
func (r Rational) set(a string) Rational {
	r.strVal = a
	r.bigRatVal = bigRational(a) // Cache the big.Rat
	return r
}

// truncate reduces the length of the numerator and denominator by the same amount
// essentially this "divides" each by the same factor of 10
// it does so in a way such that, if BOTH the numerator and denominator are longer than "numberOfDecimalPoints" (see in bigarith.go)
// then whichever is shorter will now be that length
// it is a computationally simple way of rounding the Rational when its accuracy is above the assumed maximum
// and allows rounding that may be beneficial (if the the Rational is really close to an integer or much simpler value for e.g.)
func (r Rational) truncate() Rational {
	num := r.Num().Val()
	// deal with negatives
	negStr := ""
	if num[:1] == "-" {
		num = num[1:]
		negStr = "-"
	}
	// deal with the rest
	den := r.Denom().Val()
	lenNum := len(num)
	lenDen := len(den)
	if lenDen <= numberOfDecimalPoints || lenNum <= numberOfDecimalPoints { // if EITHER is shorter or same length then exit
		// both must be longer than numberOfDecimalPoints to continue
		return r
	}
	lenShortest := lenNum // just assume this then test
	if lenNum >= lenDen { // num is longest or the same - use length of den
		lenShortest = lenDen
	}
	lenToTrunc := lenShortest - numberOfDecimalPoints // must be enough to get  shortest to right length, and leave longest same or longer
	// deal with automatic 0
	if lenDen-lenNum >= numberOfDecimalPoints { // number MUST be 0... denominator is juat oob big compared to numerator
		return NewRational("0")
	}
	// do shortening
	numNextChar := num[numberOfDecimalPoints : numberOfDecimalPoints+1]
	num = num[:(lenNum)-lenToTrunc]
	denNextChar := den[numberOfDecimalPoints : numberOfDecimalPoints+1]
	den = den[:(lenDen)-lenToTrunc]
	// make into Ints for rounding calcs
	numInt := NewInt(num)
	denInt := NewInt(den)
	numNextCharInt := NewInt(numNextChar)
	denNextCharInt := NewInt(denNextChar)
	if numNextCharInt.Compare("5") >= 0 {
		switch negStr {
		case "":
			numInt = numInt.Plus("1")
		case "-":
			numInt = numInt.Minus("1")
		}
	}
	if denNextCharInt.Compare("5") >= 0 {
		switch negStr {
		case "":
			denInt = denInt.Plus("1")
		case "-":
			denInt = denInt.Minus("1")
		}
	}
	return r.set(fmt.Sprintf("%s%s/%s", negStr, numInt.Val(), denInt.Val()))
}

// setBigRational sets the string value that represents the bigarith.Rational value
// by getting a text representation from a big.Rational object
// and returns a new Rational with the updated value
func (r Rational) setBigRational(a big.Rat) Rational {
	r.bigRatVal = &a      // Cache the big.Rat directly
	r.strVal = a.String() // Update the strVal if needed
	return r
}

// EXPORTED

// Abs returns the absolute value of the current rational number as a new bigarith.Rational
func (r Rational) Abs() Rational {
	return r.setBigRational(*new(big.Rat).Abs(r.bigRatVal)) // Use cached bigRatVal
}

// Neg returns the negated value of the current rational number as a new bigarith.Rational
func (r Rational) Neg() Rational {
	return r.setBigRational(*new(big.Rat).Neg(r.bigRatVal)) // Use cached bigRatVal
}

// Plus adds a string representation of a rational number
// to the current rational number and returns the result as a new bigarith.Rational
func (r Rational) Plus(a string) Rational {
	return r.setBigRational(*new(big.Rat).Add(r.bigRatVal, bigRational(a))) // Use cached bigRatVal
}

// Minus subtracts a string representation of a rational number
// from the current rational number and returns the result as a new bigarith.Rational
func (r Rational) Minus(a string) Rational {
	return r.setBigRational(*new(big.Rat).Sub(r.bigRatVal, bigRational(a))) // Use cached bigRatVal
}

// Diff subtracts a string representation of a rational number from the current rational number
// and takes the absolute value and returns the result as a new bigarith.Rational
func (r Rational) Diff(a string) Rational {
	// return r.Minus(a).Abs() // Use cached bigRatVal
	return r.Minus(a).Abs() // Use cached bigRatVal
}

// Times multiplies the current rational number by a string representation of another rational number
// and returns the result as a new bigarith.Rational
func (r Rational) Times(a string) Rational {
	return r.setBigRational(*new(big.Rat).Mul(r.bigRatVal, bigRational(a))) // Use cached bigRatVal
}

// DividedBy divides the current rational number by a string representation of a rational number
// and returns the result as a new bigarith.Rational
func (r Rational) DividedBy(a string) Rational {
	return r.setBigRational(*new(big.Rat).Quo(r.bigRatVal, bigRational(a))) // Use the cached bigRatVal for reuse
}

func (r Rational) SquareRootEstimate() Rational {
	// Calculate square roots of the numerator and denominator.
	numSqrt := r.Num().SquareRoot()
	denomSqrt := r.Denom().SquareRoot()

	return numSqrt.DividedBy(denomSqrt.strVal)
}

// SquareRoot calculates the square root of the current rational number
// and returns the result as a new bigarith.Rational
func (r Rational) SquareRoot() Rational {
	num := r.Num()              // Cache numerator
	denom := r.Denom()          // Cache denominator
	denomStr := denom.strVal    // Cache denominator string
	x := r.SquareRootEstimate() // Initial guess: x_0 = 1 (or choose a better approximation)
	prevXStr := "-1"            // Initial value difference check

	initFraction := num.DividedBy(denomStr)

	// Stop if diff = |xNext - x| < tolerance
	for x.Diff(prevXStr).Compare(toleranceRationalStr) >= 0 {
		prevXStr = x.strVal
		x = initFraction.DividedBy(prevXStr).Plus(prevXStr).DividedBy("2")
	}

	return x
}

// Mod calculates the modulus of the current rational number with another rational number
// and returns the result as a new bigarith.Rational
func (r Rational) Mod(a string) Rational {
	modValue := bigRational(a)                                                           // Convert the input string to a big.Rat
	divResult := new(big.Rat).Quo(r.bigRatVal, modValue)                                 // Perform r / a
	intPart := new(big.Rat).SetInt(new(big.Int).Div(divResult.Num(), divResult.Denom())) // Find the integer part of r / a
	fractionalPart := new(big.Rat).Sub(divResult, intPart)                               // Subtract the integer part from r / a to get the fractional part
	return r.setBigRational(*new(big.Rat).Mul(fractionalPart, modValue))                 // Multiply the fractional part by a to get the modulus
}

// IsEven checks if the Rational number is an even integer
func (r Rational) IsEven() bool {
	// Check if the denominator is 1, which means the number is an integer
	if r.bigRatVal.Denom().Cmp(big.NewInt(1)) == 0 {
		// Get the numerator as an integer
		numerator := r.bigRatVal.Num()
		// Check if the integer is even by taking modulus 2
		zero := big.NewInt(0)
		two := big.NewInt(2)
		mod := new(big.Int).Mod(numerator, two)
		return mod.Cmp(zero) == 0
	}
	// If not an integer, return false
	return false
}

// NthRoot calculates the nth root of the current rational number
// and returns the result as a new bigarith.Rational
func (r Rational) NthRoot(n string) Rational {
	x := r.nthRootRationalInitialGuess(n)
	prevXStr := x.Plus("1").Val() // Start with a value different from x
	// Assume n is an Int string
	nInt := NewInt(n)                                                                // Convert n to an Int
	nMinus1 := nInt.Minus("1")                                                       // n-1 Newton's method
	for x.Compare("0") != 0 && x.Diff(prevXStr).Compare(toleranceRationalStr) >= 0 { // Iterative method to find the nth root
		// Newton's method nth root: x_next = (1/n) * ((n-1) * x + value / (x^(n-1)))
		xToNMinus1 := x.ToThePowerOf(nMinus1.strVal)                               // x^(n-1)
		rOverXToNMinus1 := r.DividedBy(xToNMinus1.strVal)                          // value / (x^(n-1))
		nextX := x.Times(nMinus1.strVal).Plus(rOverXToNMinus1.strVal).DividedBy(n) // ((n-1) * x + value / (x^(n-1))) * (1/n)
		prevXStr = x.strVal
		x = nextX // Update the guess after the iteration
	}
	xTrunc := x.truncate()
	// if !nInt.IsEven() {
	// 	if r.Compare("0") < 0 {
	// 		return xTrunc.Neg()
	// 	}
	// }
	return xTrunc
}

// initialGuess estimates the initial value for Newton's method without using float64
func (r Rational) nthRootRationalInitialGuess(n string) Rational {
	// Extract the integer part of r and convert it to big.Int
	rNumBigInt := r.Num().bigIntVal
	rDenomBigInt := r.Denom().bigIntVal
	rInt := new(big.Int).Quo(rNumBigInt, rDenomBigInt)

	nRational := NewRational(n)
	nNumBigInt := nRational.Num().bigIntVal
	nDenomBigInt := nRational.Denom().bigIntVal
	nInt, _ := new(big.Int).SetString("1", 10)
	nInt.Quo(nNumBigInt, nDenomBigInt)

	// Bit-length-based approximation
	bitLen := rInt.BitLen()                 // Approximate log2 of r
	approxExp := bitLen / int(nInt.Int64()) // Rough estimate of r^(1/n) in bit length

	// 2^(approxExp) as the initial guess
	initialGuess := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(approxExp)), nil)

	// Convert the initial guess to a Rational type
	return r.set(new(big.Rat).SetInt(initialGuess).String())
}

// Mod2Pi performs x mod 2*Pi
func (r Rational) Mod2Pi() Rational {
	return r.Mod(NewRational("2").TimesPi().strVal)
}

// ToThePowerOf raises the current rational number to the power of another rational number
// and returns the result as a new bigarith.Rational
func (r Rational) ToThePowerOf(exp string) Rational {
	baseRat := bigRational(r.strVal) // Convert the string value of base to big.Rat
	expRat := bigRational(exp)       // Convert the exponent value to big.Rat

	// If exp is an integer (denominator = 1), use integer exponentiation
	if expRat.IsInt() {
		// Fast path integer exponentiation
		intExp := expRat.Num()                                      // exp as an integer (numerator of expRat)
		resultNum := new(big.Int).Exp(baseRat.Num(), intExp, nil)   // a^c
		resultDen := new(big.Int).Exp(baseRat.Denom(), intExp, nil) // b^c

		resultRat := new(big.Rat).SetFrac(resultNum, resultDen) // a^c / b^c
		return r.setBigRational(*resultRat)
	}

	// Handle (a/b)^(c/d) = (a^c / b^c)^(1/d)
	// Compute a^c and b^c
	resultNum := new(big.Int).Exp(baseRat.Num(), expRat.Num(), nil)   // a^c
	resultDen := new(big.Int).Exp(baseRat.Denom(), expRat.Num(), nil) // b^c

	// Now compute the d-th root of both resultNum and resultDen (1/d exponent)
	d := expRat.Denom()
	rootNum, err := nthRoot(resultNum, d)
	if err != nil {
		logrus.Errorf("Error calculating nth root: %v", err)
		return r // or handle the error appropriately
	}
	rootDen, err := nthRoot(resultDen, d)
	if err != nil {
		logrus.Errorf("Error calculating nth root: %v", err)
		return r // or handle the error appropriately
	}

	return r.setBigRational(*new(big.Rat).SetFrac(rootNum, rootDen)) // Combine numerator/denominator
}

// Cos calculates the cosine of the current rational number (in radians) using an optimized Taylor series
// and returns the result as a new bigarith.Rational
func (r Rational) Cos() Rational {
	// Reduce the angle to [0, 2*pi] to improve convergence speed
	r = r.Mod2Pi()

	// Initial values
	cosEstimate := NewRational("0")                        // Start with the first term 0
	previousCosEstimateStr := cosEstimate.Plus("1").strVal // Different initial estimate

	numerator := NewRational("1") // Initial term numerator for cos(x) starts at 1
	denominator := NewInt("1")    // Initial term denominator is 1
	nextFactorial := NewInt("0")  // For calculating factorials
	xSquared := r.Times(r.strVal) // Pre-compute x^2 to reduce redundant calculations
	sign := NewInt("1")           // Alternating sign for terms (starts negative after 1)

	// Loop until the difference between terms is below the tolerance
	for cosEstimate.Diff(previousCosEstimateStr).Compare(toleranceRationalStr) >= 0 {
		// Update previous estimate
		previousCosEstimateStr = cosEstimate.strVal

		// Calculate the current term (sign * numerator / denominator) and add to estimate
		term := numerator.DividedBy(denominator.strVal).Times(sign.strVal)
		cosEstimate = cosEstimate.Plus(term.strVal)

		// Update numerator and denominator for the next term
		numerator = numerator.Times(xSquared.strVal)          // Multiply by x^2
		nextFactorial = nextFactorial.Plus("1")               // add one to nextFactorial
		denominator = denominator.Times(nextFactorial.strVal) // Update denominator for factorial in Taylor expansion
		nextFactorial = nextFactorial.Plus("1")               // do it all again ...
		denominator = denominator.Times(nextFactorial.strVal) // one more time

		// Flip the sign for the next term
		sign = sign.Neg()
	}

	return cosEstimate
}

// ArcCos calculates the inverse cosine (arccos) of the current rational number
// and returns the result as a new bigarith.Rational
func (r Rational) ArcCos() Rational {
	// https://en.wikipedia.org/wiki/Inverse_trigonometric_functions#Relationships_between_trigonometric_functions_and_inverse_trigonometric_functions
	// ArcCos = arcTan(sqrt{1 - x^2} / x)
	// fmt.Printf("ArcCos for: %s\n\n", r.strVal)
	if r.Compare("1") > 0 || r.Compare("-1") < 0 {
		panic(fmt.Sprintf("ArcCos only defined -1 <= x <= 1. Value given: %s", r.strVal))
	}

	// Handle special cases directly
	if r.Compare("1") == 0 {
		return NewRational("0")
	} else if r.Compare("0") == 0 {
		return NewRational("1/2").TimesPi()
	} else if r.Compare("-1") == 0 {
		return NewRational("1").TimesPi()
	}

	xSquared := NewRational(r.strVal).Times(r.strVal)
	sqrt1MinusXSquared := NewRational("1").Minus(xSquared.strVal).SquareRoot()
	result := NewRational(sqrt1MinusXSquared.strVal).DividedBy(r.strVal).arcTan()
	if result.Compare("0") < 0 {
		result = result.Plus(NewRational("1").TimesPi().strVal)
	}
	return result
}

// LATER NON-EXPORTED FUNCS

// ArcTan calculates the inverse tangent (arctan) of the current rational number
// and returns the result as a new bigarith.Rational
func (r Rational) arcTan() Rational {
	// Check for well-known values directly
	if r.Compare("0") == 0 {
		return NewRational("0")
	} else if r.Compare("1") == 0 {
		return NewRational("1").TimesPi().DividedBy("4") // π/4
	} else if r.Compare("-1") == 0 {
		return NewRational("1").TimesPi().DividedBy("-4") // -π/4
	}

	// Use symmetry property
	if r.Compare("0") < 0 {
		return r.Neg().arcTan().Neg() // arctan(-x) = -arctan(x)
	}

	// Initialize terms for the Taylor series expansion
	rSquared := r.Times(r.strVal)
	onePlusXSquared := rSquared.Plus("1")              // denominator of it = 1 + r^2
	term := r.DividedBy(onePlusXSquared.strVal)        // initial term value = r / (1 + r^2)
	arcTanEstimate := NewRational(term.strVal)         // initial estimate
	previousArcTanEstimate := arcTanEstimate.Plus("1") // to ensure at least one loop iteration
	counter := NewRational("2")

	// Taylor series expansion
	for arcTanEstimate.Diff(previousArcTanEstimate.strVal).Compare(toleranceRational.strVal) >= 0 {
		term = term.Times(rSquared.strVal)                          // each term is multiplied by r^2
		term = term.DividedBy(onePlusXSquared.strVal)               // each term is divided by another onePlusXSquared
		term = term.Times(counter.strVal)                           // multiply by even numbers
		counter = counter.Plus("1")                                 // add 1
		term = term.DividedBy(counter.strVal)                       // divide by odd numbers
		counter = counter.Plus("1")                                 // add 1 (to get back to even)
		previousArcTanEstimate = NewRational(arcTanEstimate.strVal) // save previous value comparison in loop conditions
		arcTanEstimate = arcTanEstimate.Plus(term.strVal)           // add on the new term
	}

	return arcTanEstimate
}

func (r Rational) TimesPi() Rational {
	// https://en.wikipedia.org/wiki/Approximations_of_%CF%80#Arctangent
	// TODO: is there a reason I don't just hardcode this instead of calculating it?
	piEstimate := NewRational("1") // Starting estimate
	previousPiEstimate := piEstimate.Plus("1")
	termMultiple := NewInt("1")
	termDivisor := NewInt("3")
	nextTerm := NewRational("1")

	for piEstimate.Diff(previousPiEstimate.strVal).Compare(toleranceRational.strVal) >= 0 {
		nextTerm = nextTerm.Times(termMultiple.strVal).DividedBy(termDivisor.strVal)
		previousPiEstimate = piEstimate
		piEstimate = piEstimate.Plus(nextTerm.strVal)

		// Iterate multiple and divisor
		termMultiple = termMultiple.Plus("1")
		termDivisor = termDivisor.Plus("2")
	}

	// Pi is actually double the current estimate - and multiply it by the original value
	return piEstimate.Times("2").Times(r.strVal)
}

// TAKEN OUT FOR NOW - but kept for history / referrence

// Inv returns the inverse value of the current rational number as a new bigarith.Rational
// func (r Rational) Inv() Rational {
// 	return r.setBigRational(*new(big.Rat).Inv(r.bigRatVal)) // Use cached bigRatVal
// }

// NOTE: TRIG FUNCTIONS NOT NEEDED YET - COPIED HERE BEFORE FULLY TESTED!!

// Sin calculates the sine of the current rational number (in radians) using an optimized Taylor series
// and returns the result as a new bigarith.Rational
// func (r Rational) Sin() Rational {
// 	// Reduce the angle to [-pi, pi] to improve convergence speed
// 	r = r.Mod2Pi()

// 	// Initial values
// 	sinEstimate := NewRational("0")
// 	previousSinEstimate := sinEstimate.Plus("1") // Different initial estimate

// 	numerator := NewRational(r.strVal) // Initial term numerator is x
// 	denominator := NewInt("1")         // Initial term denominator is 1
// 	nextFactorial := NewInt("1")       // For calculating Facorials
// 	xSquared := r.Times(r.strVal)      // Pre-compute x^2 to reduce redundant calculations
// 	sign := NewInt("1")                // Alternating sign for terms

// 	// Loop until the difference between terms is below the tolerance
// 	for sinEstimate.Diff(previousSinEstimate.strVal).Compare(toleranceRational.strVal) >= 0 {
// 		// Update previous estimate
// 		previousSinEstimate = sinEstimate

// 		// Calculate the current term (sign * numerator / denominator) and add to estimate
// 		term := numerator.DividedBy(denominator.strVal).Times(sign.strVal)
// 		sinEstimate = sinEstimate.Plus(term.strVal)

// 		// Update numerator and denominator for the next term
// 		numerator = numerator.Times(xSquared.strVal)          // Multiply by x^2
// 		nextFactorial = nextFactorial.Plus("1")               // add one to nextFactorial
// 		denominator = denominator.Times(nextFactorial.strVal) // Update denominator for factorial in Taylor expansion
// 		nextFactorial = nextFactorial.Plus("1")               // do it all again ...
// 		denominator = denominator.Times(nextFactorial.strVal) // one more time

// 		// Flip the sign for the next term
// 		sign = sign.Neg()
// 	}

// 	return sinEstimate
// }

// ArcSin calculates the inverse sine (arcsin) of the current rational number
// and returns the result as a new bigarith.Rational
// func (r Rational) ArcSin() Rational {
// 	// Check if the input is within the valid range including the boundaries
// 	if r.Compare("1") > 0 || r.Compare("-1") < 0 {
// 		panic(fmt.Sprintf("ArcSin is defined for -1 <= x <= 1. Value given: %s", r.strVal))
// 	}

// 	// Handle special cases directly
// 	if r.Compare("1") == 0 {
// 		return NewRational("1").TimesPi().DividedBy("2") // π/2
// 	} else if r.Compare("0") == 0 {
// 		return NewRational("0")
// 	} else if r.Compare("-1") == 0 {
// 		return NewRational("-1").TimesPi().DividedBy("2") // -π/2
// 	}

// 	// General case using the relationship with arctan
// 	xSquared := NewRational(r.strVal).Times(r.strVal)
// 	sqrt1MinusXSquared := NewRational("1").Minus(xSquared.strVal).SquareRoot()
// 	return NewRational(r.strVal).DividedBy(sqrt1MinusXSquared.strVal).arcTan()
// }
