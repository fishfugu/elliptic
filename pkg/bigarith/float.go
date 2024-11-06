package bigarith

import (
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
)

// Float is a bigarith type for Floating Point Arithmetic
type Float struct {
	strVal string // the string representing the floating point value
}

func NewFloat(a string) Float {
	return new(Float).set(a)
}

func (f Float) truncDecPoints() Float {
	// NB: for now - assume max num of decimal points we want is 2 less than standard
	// max precision being used... because the last 2 decimal points is where we
	// start to see the rounding errors in loops etc.
	// log_2(10) approx. 3.321928094887362 < 3.322
	// so for 512 decimal places let's say 512 * 3.322 = 1700.84 < 1800
	bigFloatF, _ := new(big.Float).SetString(f.strVal)
	return NewFloat(bigFloatF.Text('f', numberOfDecimalPoints))
}

// Discovery functions
// the functions that surface info about the Float - don't return a Float itself

// Val returns the string representation of the bigarith.Float value
func (f Float) Val() string {
	return f.strVal
}

// Compare takes a string representation of a floating point number and compares
// it with the current bigarith.Float,
// returns:
// -1 if x <  y, 0 if x == y (incl. -0 == 0, -Inf == -Inf, and +Inf == +Inf), +1 if x >  y
func (f Float) Compare(a string) int {
	return bigFloat(f.strVal).Cmp(bigFloat(a))
}

// Manipulation functions
// functions that manipulate the value of a Float
// should return the new value without modifying the original Float

// set sets the string value that represents the3bigarith.Float value
// and returns a new Float with the updated value
func (f Float) set(a string) Float {
	f.strVal = a
	return f
}

// setBigFloat sets the string value that represents the bigarith.Float value
// by getting a text representation from a big.Float object
// and returns a new Float with the updated value
func (f Float) setBigFloat(a big.Float) Float {
	// number of dec points after the decimial point ("prec") is twice as many digits as needed to write out the integer part
	// aFloat, _ := a.Float64()
	// prec := (math.Log(aFloat) + 10.0) * 2
	return f.set(a.Text('f', numberOfDecimalPoints+4))
}

// setRational sets a bigarith.Float value by dividing the Numerator by the Denominator, in Floating Point arithmatic
func (f Float) setRational(a Rational) Float {
	NewFloat(a.Num().strVal).DividedBy(a.Denom().strVal)
	return NewFloat(a.Num().strVal).DividedBy(a.Denom().strVal)
}

// Neg returns the negated value of the current floating point number as a new bigarith.Float
func (f Float) Neg() Float {
	return f.setBigFloat(*new(big.Float).Neg(bigFloat(f.strVal)))
}

// Plus adds a string representation of a floating point number
// to the current floating point number and returns the result as a new bigarith.Float
func (f Float) Plus(a string) Float {
	return f.setBigFloat(*new(big.Float).Add(bigFloat(f.strVal), bigFloat(a)))
}

// Minus subtracts a string representation of a floating point number
// from the current floating point number and returns the result as a new bigarith.Float
func (f Float) Minus(a string) Float {
	return f.setBigFloat(*new(big.Float).Sub(bigFloat(f.strVal), bigFloat(a)))
}

// Diff subtracts a string representation of a floating point number from the current floating point number
// and takes the absolute value and returns the result as a new bigarith.Float
func (r Float) Diff(a string) Float {
	minusFloat := r.Minus(a)
	if minusFloat.Compare("0") < 0 {
		return minusFloat.Neg()
	}
	return minusFloat
}

// Times multiplies the current floating point number by a string representation of another floating point number
// and returns the result as a new bigarith.Float
func (f Float) Times(a string) Float {
	return f.setBigFloat(*new(big.Float).Mul(bigFloat(f.strVal), bigFloat(a)))
}

// DividedBy divides the current floating point number by a string representation of a floating point number
// and returns the result as a new bigarith.Float
func (f Float) DividedBy(a string) Float {
	return f.setBigFloat(*new(big.Float).Quo(bigFloat(f.strVal), bigFloat(a)))
}

// SquareRoot calculates the square root of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) SquareRoot() Float {
	return f.setBigFloat(*new(big.Float).Sqrt(bigFloat(f.strVal)))
}

func (f Float) NthRoot(n string) Float {
	oneOverN := NewFloat("1").DividedBy(n)
	logrus.Debugf("Calculating nth root of %s, n = %s", f.strVal, n)
	if f.Compare("0") < 0 {
		logrus.Debugf("Result %s", f.Neg().ToThePowerOf(oneOverN.strVal).Neg().strVal)
		return f.Neg().ToThePowerOf(oneOverN.strVal).Neg()
	} else {
		logrus.Debugf("Result %s", f.ToThePowerOf(oneOverN.strVal).strVal)
		return f.ToThePowerOf(oneOverN.strVal)
	}
}

// Mod calculates the modulus of the current floating point number with another floating point number
// and returns the result as a new bigarith.Float
func (f Float) Mod(a string) Float {
	q := NewFloat(f.strVal).DividedBy(a)
	qInt, _ := bigFloat(q.strVal).Int(nil)
	return f.Minus(NewFloat(qInt.String()).Times(a).strVal)
}

func (f Float) Mod2Pi() Float {
	// Use a more precise Pi value for accurate modulo calculations
	// Perform x mod 2*Pi
	return f.Mod(NewFloat("2").TimesPi().strVal)
}

// ToThePowerOf raises the current floating point number to the power of another floating point number
// and returns the result as a new bigarith.Float
func (f Float) ToThePowerOf(a string) Float {
	if f.Compare("0") == 0 {
		return NewFloat("0")
	}
	if f.Compare("0") < 0 {
		panic(fmt.Sprintf("base value for Float.ToThePowerOf cannot be negative - base = %s - exponent = %s", f.strVal, a))
	}
	logBase := logNewton(bigFloat(f.strVal))
	expLogBase := new(big.Float).Mul(bigFloat(a), logBase)
	// problem is I don't know how accurate you make this ToThePowerOf function...
	// TODO: review this... how does it interact with the for loop counts in logNewton and expSeries
	return f.setBigFloat(*expSeries(expLogBase))
}

// Sin calculates the sine of the current floating point number (in radians) using the Taylor series expansion
// and returns the result as a new bigarith.Float
func (f Float) Sin() Float {
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Simple_algebraic_values
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Power_series_expansion
	logrus.Debugf("Calculating Cos for: %s", f.strVal)
	// take input back to === mod 2 pi
	f = f.Mod(NewFloat("2").TimesPi().strVal)
	if f.Compare("0") == 0 {
		return NewFloat("0")
	}
	if f.Compare(NewFloat("1").TimesPi().DividedBy("6").strVal) == 0 {
		return NewFloat(NewFloat("1").DividedBy("2").strVal)
	}
	if f.Compare(NewFloat("1").TimesPi().DividedBy("2").strVal) == 0 {
		return NewFloat("1")
	}
	sinEstimate := NewFloat("0")
	previousSinEstimate := sinEstimate.Plus("1")
	// setup numerator
	numerator := NewFloat(f.strVal)
	// setup denominator
	multiple := NewInt("1")
	denominator := NewInt("1")
	// setup the sign
	sign := NewInt("1")
	// recalc term
	term := numerator.DividedBy(denominator.strVal)
	for sinEstimate.Compare(previousSinEstimate.strVal) != 0 {
		previousSinEstimate = sinEstimate
		sinEstimate = sinEstimate.Plus(term.strVal)
		logrus.Debugf("sinEstimate: %s", sinEstimate.strVal)
		// recalc values
		// recalc numerator
		numerator = numerator.Times(f.strVal).Times(f.strVal)
		logrus.Debugf("numerator: %s", numerator.strVal)
		// recalc denominator
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.strVal)
		denominator = denominator.Times(multiple.strVal)
		logrus.Debugf("denominator: %s", denominator.strVal)
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.strVal)
		denominator = denominator.Times(multiple.strVal)
		logrus.Debugf("denominator: %s", denominator.strVal)
		// Alternate the sign (positive/negative)
		sign = sign.Neg()
		logrus.Debugf("sign: %s", sign.strVal)
		// recalc term
		term = numerator.DividedBy(denominator.strVal).Times(sign.strVal)
		logrus.Debugf("Latest term: %s", term.strVal)
	}

	sinEstimate = sinEstimate.Plus(term.strVal) // one last time... even after it hasn't changed
	return sinEstimate
}

// Cos calculates the cosine of the current floating point number (in radians) using the Taylor series expansion
// and returns the result as a new bigarith.Float
func (f Float) Cos() Float {
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Simple_algebraic_values
	// https://en.wikipedia.org/wiki/Trigonometric_functions#Power_series_expansion
	logrus.Debugf("Calculating Cos for: %s", f.strVal)
	// take input back to === mod 2 pi
	f = f.Mod(NewFloat("2").TimesPi().strVal)
	if f.Compare("0") == 0 {
		return NewFloat("1")
	}
	if f.Compare(NewFloat("1").TimesPi().DividedBy("3").strVal) == 0 {
		return NewFloat(NewFloat("1").DividedBy("2").strVal)
	}
	if f.Compare(NewFloat("1").TimesPi().DividedBy("2").strVal) == 0 {
		return NewFloat(NewFloat("0").strVal)
	}
	cosEstimate := NewFloat("0")
	previousCosEstimate := cosEstimate.Plus("-1")
	// setup numerator
	numerator := NewFloat("1")
	logrus.Debugf("original numerator: %s", numerator.strVal)
	// setup denominator
	multiple := NewInt("0")
	denominator := NewInt("1")
	logrus.Debugf("original denominator: %s", denominator.strVal)
	// setup the sign
	sign := NewInt("1")
	// recalc term
	term := numerator.DividedBy(denominator.strVal)
	logrus.Debugf("original term: %s", term.strVal)
	for cosEstimate.Compare(previousCosEstimate.strVal) != 0 {
		previousCosEstimate = cosEstimate
		cosEstimate = cosEstimate.Plus(term.strVal)
		logrus.Debugf("cosEstimate: %s", cosEstimate.strVal)
		// recalc values
		// recalc numerator
		numerator = numerator.Times(f.strVal).Times(f.strVal)
		logrus.Debugf("numerator: %s", numerator.strVal)
		// recalc denominator
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.strVal)
		denominator = denominator.Times(multiple.strVal)
		logrus.Debugf("denominator: %s", denominator.strVal)
		multiple = multiple.Plus("1")
		logrus.Debugf("multiple: %s", multiple.strVal)
		denominator = denominator.Times(multiple.strVal)
		logrus.Debugf("denominator: %s", denominator.strVal)
		// Alternate the sign (positive/negative)
		sign = sign.Neg()
		logrus.Debugf("sign: %s", sign.strVal)
		// recalc term
		term = numerator.DividedBy(denominator.strVal).Times(sign.strVal)
		logrus.Debugf("Latest term: %s\n", term.strVal)
	}

	cosEstimate = cosEstimate.Plus(term.strVal) // one last time... even after it hasn't changed
	return cosEstimate
}

// ArcCos calculates the inverse cosine (arccos) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcCos() Float {
	// https://en.wikipedia.org/wiki/Inverse_trigonometric_functions#Relationships_between_trigonometric_functions_and_inverse_trigonometric_functions
	// ArcCos = ArcTan(sqrt{1 - x^2} / x)
	if f.Compare("1") > 0 || f.Compare("-1") < 0 {
		panic(fmt.Sprintf("ArcCos only defined for -1 <= x <= 1. Value given: %s", f.strVal))
	}
	if f.Compare("1") == 0 {
		return NewFloat("0")
	}
	if f.Compare("0") == 0 {
		return NewFloat("1").TimesPi().DividedBy("2")
	}
	if f.Compare("-1") == 0 {
		return NewFloat("1").TimesPi()
	}
	xSquared := NewFloat(f.strVal).Times(f.strVal)
	sqrt1MinusXSquared := NewFloat("1").Minus(xSquared.strVal).SquareRoot()
	result := NewFloat(sqrt1MinusXSquared.strVal).DividedBy(f.strVal).ArcTan()
	if result.Compare("0") < 0 {
		result = result.Plus(NewFloat("1").TimesPi().strVal)
	}
	return result
}

// ArcSin calculates the inverse sine (arcsin) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcSin() Float {
	// https://en.wikipedia.org/wiki/Inverse_trigonometric_functions#Relationships_between_trigonometric_functions_and_inverse_trigonometric_functions
	// ArcSin = ArcTan(x / sqrt{1 - x^2})
	if f.Compare("1") > 0 || f.Compare("-1") < 0 {
		panic(fmt.Sprintf("ArcSin only defined for -1 <= x <= 1. Value given: %s", f.strVal))
	}
	if f.Compare("1") == 0 {
		return NewFloat("1").TimesPi().DividedBy("2")
	}
	if f.Compare("0") == 0 {
		return NewFloat("0")
	}
	if f.Compare("-1") == 0 {
		return NewFloat("1").TimesPi().DividedBy("-2")
	}
	xSquared := NewFloat(f.strVal).Times(f.strVal)
	sqrt1MinusXSquared := NewFloat("1").Minus(xSquared.strVal).SquareRoot()
	return NewFloat(f.strVal).DividedBy(sqrt1MinusXSquared.strVal).ArcTan()
}

// ArcTan calculates the inverse tangent (arctan) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcTan() Float {
	// ArcTan is helpful because its domain is all x... and it helps define other inverse trig funcs
	// see: /doc/img/arctna.png
	// Taylor series expansion for arctan(x)
	// https://en.wikipedia.org/wiki/Arctangent_series#Accelerated_series
	onePlusXSquared := f.Times(f.strVal).Plus("1")     // Right Hand Multiple, the denominator of it = 1 + f^2
	term := f.DividedBy(onePlusXSquared.strVal)        // initial term value = f / 1 + f^2
	arcTanEstimate := NewFloat("0").Plus(term.strVal)  // intial estimate is initial term
	previousArcTanEstimate := arcTanEstimate.Plus("1") // something different to arcTanEstimate
	counter := NewFloat("2")                           // counter for multiples and division
	for arcTanEstimate.Compare(previousArcTanEstimate.strVal) != 0 {
		term = term.Times(f.strVal).Times(f.strVal)       // each term is multiplied by f^2
		term = term.DividedBy(onePlusXSquared.strVal)     // each term is divided by another onePlusXSquared
		term = term.Times(counter.strVal)                 // multiply by even numbers
		counter = counter.Plus("1")                       // add 1
		term = term.DividedBy(counter.strVal)             // divide by odd numbers
		counter = counter.Plus("1")                       // add 1 (to get back to even)
		previousArcTanEstimate = arcTanEstimate           // save previous value for comparison in loop conditions
		arcTanEstimate = arcTanEstimate.Plus(term.strVal) // add on the new term
	}

	return arcTanEstimate
}

func (r Float) TimesPi() Float {
	// https://en.wikipedia.org/wiki/Approximations_of_%CF%80#Arctangent
	// TODO: is there a reason I don't just hardcode this instead of calculating it?
	piEstimate := NewFloat("1")
	previousPiEstimate := piEstimate.Plus("1")
	termMultiple := NewInt("1")
	termDivisor := NewInt("3")

	nextTerm := NewFloat("1")
	for piEstimate.Compare(previousPiEstimate.strVal) != 0 {
		nextTerm = nextTerm.Times(termMultiple.strVal).DividedBy(termDivisor.strVal)
		previousPiEstimate = piEstimate
		piEstimate = piEstimate.Plus(nextTerm.strVal)

		// iterate multiple and divisor
		termMultiple = termMultiple.Plus("1")
		termDivisor = termDivisor.Plus("2")
	}
	// Pi is actually double the current estimate
	piEstimate = piEstimate.Times("2")

	return piEstimate
}
