package bigarith

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

// min value
const numberOfDecimalPoints = int(500) // from 48 ... 505 - eveything else seems to error
const precision = uint(numberOfDecimalPoints * 17 / 5)

func (f Float) truncateToMaxNumOfDecimalPoints() Float {
	// NB: for now - assume max num of decimal points we want is 2 less than standard
	// max precision being used... because the last 2 decimal points is where we
	// start to see the rounding errors in loops etc.
	// log_2(10) approx. 3.321928094887362 < 3.322
	// so for 512 decimal places let's say 512 * 3.322 = 1700.84 < 1800
	bigFloatF, _ := new(big.Float).SetPrec(precision).SetString(f.Val())
	return NewFloat(bigFloatF.Text('f', numberOfDecimalPoints))
}

// Float is a bigarith type for Floating Point Arithmetic
type Float struct {
	stringValue string // the string representing the floating point value
}

func NewFloat(a string) Float {
	return new(Float).set(a)
}

// Discovery functions
// the functions that surface info about the Float - don't return a Float itself

// Val returns the string representation of the bigarith.Float value
func (f Float) Val() string {
	return f.stringValue
}

// Compare takes a string representation of a floating point number and compares
// it with the current bigarith.Float,
// returns:
// -1 if x <  y, 0 if x == y (incl. -0 == 0, -Inf == -Inf, and +Inf == +Inf), +1 if x >  y
func (f Float) Compare(a string) int {
	return bigFloat(f.stringValue).Cmp(bigFloat(a))
}

// Manipulation functions
// functions that manipulate the value of a Float
// should return the new value without modifying the original Float

// set sets the string value that represents the bigarith.Float value
// and returns a new Float with the updated value
func (f Float) set(a string) Float {
	f.stringValue = a
	return f
}

// SetBigFloat sets the string value that represents the bigarith.Float value
// by getting a text representation from a big.Float object
// and returns a new Float with the updated value
func (f Float) SetBigFloat(a big.Float) Float {
	// number of dec points after the decimial point ("prec") is twice as many digits as needed to write out the integer part
	// aFloat, _ := a.Float64()
	// prec := (math.Log(aFloat) + 10.0) * 2
	return f.set(a.Text('f', numberOfDecimalPoints*2)) // this seems to need to be higher
}

// Neg returns the negated value of the current floating point number as a new bigarith.Float
func (f Float) Neg() Float {
	return f.SetBigFloat(*new(big.Float).Neg(bigFloat(f.stringValue)))
}

// Plus adds a string representation of a floating point number
// to the current floating point number and returns the result as a new bigarith.Float
func (f Float) Plus(a string) Float {
	return f.SetBigFloat(*new(big.Float).Add(bigFloat(f.stringValue), bigFloat(a)))
}

// Minus subtracts a string representation of a floating point number
// from the current floating point number and returns the result as a new bigarith.Float
func (f Float) Minus(a string) Float {
	return f.SetBigFloat(*new(big.Float).Sub(bigFloat(f.stringValue), bigFloat(a)))
}

// Times multiplies the current floating point number by a string representation of another floating point number
// and returns the result as a new bigarith.Float
func (f Float) Times(a string) Float {
	return f.SetBigFloat(*new(big.Float).Mul(bigFloat(f.stringValue), bigFloat(a)))
}

// DividedBy divides the current floating point number by a string representation of a floating point number
// and returns the result as a new bigarith.Float
func (f Float) DividedBy(a string) Float {
	return f.SetBigFloat(*new(big.Float).Quo(bigFloat(f.stringValue), bigFloat(a)))
}

// SquareRoot calculates the square root of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) SquareRoot() Float {
	return f.SetBigFloat(*new(big.Float).Sqrt(bigFloat(f.stringValue)))
}

func (f Float) NthRoot(n string) Float {
	oneOverN := NewFloat("1").DividedBy(n)
	logrus.Debugf("Calculating nth root of %s, n = %s: %%s", f.Val(), n)
	if f.Compare("0") < 0 {
		logrus.Debugf("Result %s", f.Neg().ToThePowerOf(oneOverN.Val()).Neg().Val())
		return f.Neg().ToThePowerOf(oneOverN.Val()).Neg()
	} else {
		logrus.Debugf("Result %s", f.ToThePowerOf(oneOverN.Val()).Val())
		return f.ToThePowerOf(oneOverN.Val())
	}
}

// Mod calculates the modulus of the current floating point number with another floating point number
// and returns the result as a new bigarith.Float
func (f Float) Mod(a string) Float {
	q := NewFloat(f.Val()).DividedBy(a)
	qInt, _ := bigFloat(q.Val()).Int(nil)
	return f.Minus(NewFloat(qInt.String()).Times(a).Val())
}

func (f Float) Mod2Pi() Float {
	// Use a more precise Pi value for accurate modulo calculations
	// Perform x mod 2*Pi
	return f.Mod(Pi().Times("2").Val())
}

// ToThePowerOf raises the current floating point number to the power of another floating point number
// and returns the result as a new bigarith.Float
func (f Float) ToThePowerOf(a string) Float {
	if f.Compare("0") == 0 {
		return NewFloat("0")
	}
	if f.Compare("0") < 0 {
		panic(fmt.Sprintf("base value for Float.ToThePowerOf cannot be negative - base = %s - exponent = %s", f.Val(), a))
	}
	logBase := logNewton(bigFloat(f.Val()))
	expLogBase := new(big.Float).Mul(bigFloat(a), logBase)
	// problem is I don't know how accurate you make this ToThePowerOf function...
	// TODO: review this... how does it interact with the for loop counts in logNewton and expSeries
	return f.set(expSeries(expLogBase).Text('f', numberOfDecimalPoints))
}

// Sin calculates the sine of the current floating point number (in radians) using the Taylor series expansion
// and returns the result as a new bigarith.Float
func (f Float) Sin() Float {
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Simple_algebraic_values
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Power_series_expansion
	logrus.Debugf("Calculating Cos for: %s", f.Val())
	// take input back to === mod 2 pi
	f = f.Mod(Pi().Times("2").Val())
	if f.Compare("0") == 0 {
		return NewFloat("0")
	}
	if f.Compare(Pi().DividedBy("6").Val()) == 0 {
		return NewFloat(NewFloat("1").DividedBy("2").Val())
	}
	if f.Compare(Pi().DividedBy("2").Val()) == 0 {
		return NewFloat("1")
	}
	sinEstimate := NewFloat("0")
	previousSinEstimate := sinEstimate.Plus("1")
	// setup numerator
	numerator := NewFloat(f.Val())
	// setup denominator
	multiple := NewFloat("1")
	denominator := NewFloat("1")
	// setup the sign
	sign := NewFloat("1")
	// recalc term
	term := numerator.DividedBy(denominator.Val())
	for sinEstimate.truncateToMaxNumOfDecimalPoints().Compare(previousSinEstimate.truncateToMaxNumOfDecimalPoints().Val()) != 0 {
		previousSinEstimate = sinEstimate
		sinEstimate = sinEstimate.Plus(term.Val())
		logrus.Debugf("sinEstimate: %s", sinEstimate.Val())
		// recalc values
		// recalc numerator
		numerator = numerator.Times(f.Val()).Times(f.Val())
		logrus.Debugf("numerator: %s", numerator.Val())
		// recalc denominator
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.Val())
		denominator = denominator.Times(multiple.Val())
		logrus.Debugf("denominator: %s", denominator.Val())
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.Val())
		denominator = denominator.Times(multiple.Val())
		logrus.Debugf("denominator: %s", denominator.Val())
		// Alternate the sign (positive/negative)
		sign = sign.Neg()
		logrus.Debugf("sign: %s", sign.Val())
		// recalc term
		term = numerator.DividedBy(denominator.Val()).Times(sign.Val())
		logrus.Debugf("Latest term: %s", term.Val())
	}

	sinEstimate = sinEstimate.Plus(term.Val()) // one last time... even after it hasn't changed
	return sinEstimate.truncateToMaxNumOfDecimalPoints()
}

// Cos calculates the cosine of the current floating point number (in radians) using the Taylor series expansion
// and returns the result as a new bigarith.Float
func (f Float) Cos() Float {
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Simple_algebraic_values
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Power_series_expansion
	logrus.Debugf("Calculating Cos for: %s", f.Val())
	// take input back to === mod 2 pi
	f = f.Mod(Pi().Times("2").Val())
	if f.Compare("0") == 0 {
		return NewFloat("1")
	}
	if f.Compare(Pi().DividedBy("3").Val()) == 0 {
		return NewFloat(NewFloat("1").DividedBy("2").Val())
	}
	if f.Compare(Pi().DividedBy("2").Val()) == 0 {
		return NewFloat(NewFloat("0").Val())
	}
	cosEstimate := NewFloat("0")
	previousCosEstimate := cosEstimate.Plus("-1")
	// setup numerator
	numerator := NewFloat("1")
	logrus.Debugf("original numerator: %s", numerator.Val())
	// setup denominator
	multiple := NewFloat("0")
	denominator := NewFloat("1")
	logrus.Debugf("original denominator: %s", denominator.Val())
	// setup the sign
	sign := NewFloat("1")
	// recalc term
	term := numerator.DividedBy(denominator.Val())
	logrus.Debugf("original term: %s", term.Val())
	for cosEstimate.truncateToMaxNumOfDecimalPoints().Compare(previousCosEstimate.truncateToMaxNumOfDecimalPoints().Val()) != 0 {
		previousCosEstimate = cosEstimate
		cosEstimate = cosEstimate.Plus(term.Val())
		logrus.Debugf("cosEstimate: %s", cosEstimate.Val())
		// recalc values
		// recalc numerator
		numerator = numerator.Times(f.Val()).Times(f.Val())
		logrus.Debugf("numerator: %s", numerator.Val())
		// recalc denominator
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.Val())
		denominator = denominator.Times(multiple.Val())
		logrus.Debugf("denominator: %s", denominator.Val())
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.Val())
		denominator = denominator.Times(multiple.Val())
		logrus.Debugf("denominator: %s", denominator.Val())
		// Alternate the sign (positive/negative)
		sign = sign.Neg()
		logrus.Debugf("sign: %s", sign.Val())
		// recalc term
		term = numerator.DividedBy(denominator.Val()).Times(sign.Val())
		logrus.Debugf("Latest term: %s\n", term.Val())
	}

	cosEstimate = cosEstimate.Plus(term.Val()) // one last time... even after it hasn't changed
	return cosEstimate.truncateToMaxNumOfDecimalPoints()
}

// ArcCos calculates the inverse cosine (arccos) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcCos() Float {
	// https://en.wikipedia.org/wiki/Inverse_trigonometric_functions#Relationships_between_trigonometric_functions_and_inverse_trigonometric_functions
	// ArcCos = ArcTan(sqrt{1 - x^2} / x)
	if f.Compare("1") > 0 || f.Compare("-1") < 0 {
		panic(fmt.Sprintf("ArcCos only defined for -1 <= x <= 1. Value given: %s", f.Val()))
	}
	if f.Compare("1") == 0 {
		return NewFloat("0")
	}
	if f.Compare("0") == 0 {
		return NewFloat(Pi().DividedBy("2").truncateToMaxNumOfDecimalPoints().Val())
	}
	if f.Compare("-1") == 0 {
		return NewFloat(Pi().truncateToMaxNumOfDecimalPoints().Val())
	}
	xSquared := NewFloat(f.Val()).Times(f.Val())
	sqrt1MinusXSquared := NewFloat("1").Minus(xSquared.Val()).SquareRoot()
	result := NewFloat(sqrt1MinusXSquared.Val()).DividedBy(f.Val()).ArcTan()
	if result.Compare("0") < 0 {
		result = result.Plus(Pi().Val())
	}
	return result.truncateToMaxNumOfDecimalPoints()
}

// ArcSin calculates the inverse sine (arcsin) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcSin() Float {
	// https://en.wikipedia.org/wiki/Inverse_trigonometric_functions#Relationships_between_trigonometric_functions_and_inverse_trigonometric_functions
	// ArcSin = ArcTan(x / sqrt{1 - x^2})
	if f.Compare("1") > 0 || f.Compare("-1") < 0 {
		panic(fmt.Sprintf("ArcSin only defined for -1 <= x <= 1. Value given: %s", f.Val()))
	}
	if f.Compare("1") == 0 {
		return NewFloat(Pi().DividedBy("2").truncateToMaxNumOfDecimalPoints().Val())
	}
	if f.Compare("0") == 0 {
		return NewFloat("0")
	}
	if f.Compare("-1") == 0 {
		return NewFloat(Pi().DividedBy("-2").truncateToMaxNumOfDecimalPoints().Val())
	}
	xSquared := NewFloat(f.Val()).Times(f.Val())
	sqrt1MinusXSquared := NewFloat("1").Minus(xSquared.Val()).SquareRoot()
	return NewFloat(f.Val()).DividedBy(sqrt1MinusXSquared.Val()).ArcTan()
}

// ArcTan calculates the inverse tangent (arctan) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcTan() Float {
	// ArcTan is helpful because its domain is all x... and it helps define other inverse trig funcs
	// see: /doc/img/arctna.png
	// Taylor series expansion for arctan(x)
	// https://en.wikipedia.org/wiki/Arctangent_series#Accelerated_series
	onePlusXSquared := f.Times(f.Val()).Plus("1")      // Right Hand Multiple, the denominator of it = 1 + f^2
	term := f.DividedBy(onePlusXSquared.Val())         // initial term value = f / 1 + f^2
	arcTanEstimate := NewFloat("0").Plus(term.Val())   // intial estimate is initial term
	previousArcTanEstimate := arcTanEstimate.Plus("1") // something different to arcTanEstimate
	counter := NewFloat("2")                           // counter for multiples and division
	for arcTanEstimate.truncateToMaxNumOfDecimalPoints().Compare(previousArcTanEstimate.truncateToMaxNumOfDecimalPoints().Val()) != 0 {
		term = term.Times(f.Val()).Times(f.Val())        // each term is multiplied by f^2
		term = term.DividedBy(onePlusXSquared.Val())     // each term is divided by another onePlusXSquared
		term = term.Times(counter.Val())                 // multiply by even numbers
		counter = counter.Plus("1")                      // add 1
		term = term.DividedBy(counter.Val())             // divide by odd numbers
		counter = counter.Plus("1")                      // add 1 (to get back to even)
		previousArcTanEstimate = arcTanEstimate          // save previous value for comparison in loop conditions
		arcTanEstimate = arcTanEstimate.Plus(term.Val()) // add on the new term
	}

	return arcTanEstimate.truncateToMaxNumOfDecimalPoints()
}

func Pi() Float {
	// https://en.wikipedia.org/wiki/Approximations_of_%CF%80#Arctangent
	// TODO: is there a reason I don't just hardcode this instead of calculating it?
	piEstimate := NewFloat("1")
	previousPiEstimate := piEstimate.Plus("1")
	termMultiple := NewInt("1")
	termDivisor := NewInt("3")

	nextTerm := NewFloat("1")
	for piEstimate.truncateToMaxNumOfDecimalPoints().Compare(previousPiEstimate.truncateToMaxNumOfDecimalPoints().Val()) != 0 {
		nextTerm = nextTerm.Times(termMultiple.Val()).DividedBy(termDivisor.Val())
		previousPiEstimate = piEstimate
		piEstimate = piEstimate.Plus(nextTerm.Val())

		// iterate multiple and divisor
		termMultiple = termMultiple.Plus("1")
		termDivisor = termDivisor.Plus("2")
	}
	// Pi is actually double the current estimate
	piEstimate = piEstimate.Times("2")

	return piEstimate.truncateToMaxNumOfDecimalPoints()
}
