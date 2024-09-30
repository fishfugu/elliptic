package bigarith

import (
	"fmt"
	"math/big"
	"strconv"
)

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
	return f.set(a.Text('f', 2048))
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
	if f.Compare("0") < 0 {
		return f.Neg().ToThePowerOf(oneOverN.Val()).Neg()
	} else {
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

// ToThePowerOf raises the current floating point number to the power of another floating point number
// and returns the result as a new bigarith.Float
func (f Float) ToThePowerOf(a string) Float {
	if f.Compare("0") <= 0 {
		panic(fmt.Sprintf("base value for Float.ToThePowerOf cannot be negative - base = %s - exponent = %s", f.Val(), a))
	}
	logBase := logNewton(bigFloat(f.Val()))
	expLogBase := new(big.Float).Mul(bigFloat(a), logBase)
	return f.SetBigFloat(*expSeries(expLogBase))
}

// Sin calculates the sine of the current floating point number (in radians) using the Taylor series expansion
// and returns the result as a new bigarith.Float
func (f Float) Sin() Float {
	f = f.Mod2Pi() // Modulo 2π to normalise the input
	result := NewFloat("0")
	aSquared := NewFloat(f.Val()).Times(f.Val()) // Square of the current value
	term := NewFloat(f.Val())                    // First term is just the current value
	factorial := NewFloat("1")                   // Initialise factorial for Taylor expansion
	sign := NewFloat("1")                        // Initialise sign for alternating series

	for i := 1; i < 100; i++ {
		if i > 1 {
			// Update factorial to (2*i - 1) * (2*i)
			factorial = factorial.Times(strconv.Itoa(2*i - 1)).Times(strconv.Itoa(2 * i))
			// Calculate the next term in the series
			term = term.Times(aSquared.Val()).DividedBy(factorial.Val()).Times(sign.Val())
		}
		// Add the term to the result
		result = result.Plus(term.Val())
		// Alternate the sign (positive/negative)
		sign = sign.Neg()
	}
	// Return the final sine value
	return f.set(result.Val())
}

// Cos calculates the cosine of the current floating point number (in radians) using the Taylor series expansion
// and returns the result as a new bigarith.Float
func (f Float) Cos() Float {
	f = f.Mod2Pi()                               // Modulo 2π to normalise the input
	result := NewFloat("1")                      // First term in the Taylor series for cos is 1
	aSquared := NewFloat(f.Val()).Times(f.Val()) // Square of the current value
	term := NewFloat("1")                        // Initialise the first term as 1
	factorial := NewFloat("1")                   // Initialise factorial for Taylor expansion
	sign := NewFloat("1")                        // Initialise sign for alternating series

	for i := 1; i < 100; i++ {
		// Update factorial for the current term
		factorial = factorial.Times(strconv.Itoa(2*i - 1)).Times(strconv.Itoa(2 * i))
		// Calculate the next term in the series
		term = term.Times(aSquared.Val()).DividedBy(factorial.Val()).Times(sign.Val())
		// Add the term to the result
		result = result.Plus(term.Val())
		// Alternate the sign (positive/negative)
		sign = sign.Neg()
	}
	// Return the final cosine value
	return f.set(result.Val())
}

// ArcCos calculates the inverse cosine (arccos) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcCos() Float {
	// Check if the input is within the valid range [-1, 1]
	if f.Compare("-1") < 0 || f.Compare("1") > 0 {
		panic(fmt.Sprintf("Input for ArcCos must be in the range [-1, 1]. Got: %s", f.Val()))
	}

	// ArcCos(x) = Pi/2 - ArcSin(x)
	arcSinResult := f.ArcSin()
	piOver2 := Pi().DividedBy("2")

	return piOver2.Minus(arcSinResult.Val())
}

// ArcSin calculates the inverse sine (arcsin) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcSin() Float {
	// Check if the input is within the valid range [-1, 1]
	if f.Compare("-1") < 0 || f.Compare("1") > 0 {
		panic(fmt.Sprintf("Input for ArcSin must be in the range [-1, 1]. Got: %s", f.Val()))
	}

	// ArcSin(x) = arctan(x / sqrt(1 - x^2))
	one := NewFloat("1")
	xSquared := f.Times(f.Val())                          // x^2
	oneMinusXSquared := one.Minus(xSquared.Val())         // 1 - x^2
	sqrtOneMinusXSquared := oneMinusXSquared.SquareRoot() // sqrt(1 - x^2)

	// x / sqrt(1 - x^2)
	arcTanInput := f.DividedBy(sqrtOneMinusXSquared.Val())
	return arcTanInput.ArcTan() // Use ArcTan (if you implement it) or approximate
}

// ArcTan calculates the inverse tangent (arctan) of the current floating point number
// and returns the result as a new bigarith.Float
func (f Float) ArcTan() Float {
	// Taylor series expansion for arctan(x) when |x| <= 1
	// arctan(x) = x - (x^3)/3 + (x^5)/5 - (x^7)/7 + ...
	result := NewFloat(f.Val())
	numerator := NewFloat(f.Val())
	sign := NewFloat("-1")

	for i := 3; i < 256; i += 2 {
		numerator = numerator.Times(f.Val()).Times(f.Val())            // numerator = x^i
		term := numerator.DividedBy(strconv.Itoa(i)).Times(sign.Val()) // term = +/- x^i / i = +/- numerator / i
		result = result.Plus(term.Val())                               // accumulate
		sign = sign.Neg()                                              // alternate sign
	}

	return result
}

func Pi() Float {
	// Constants in the Chudnovsky algorithm
	// C = 426880 * sqrt(10005)
	C := NewFloat("426880")
	K := NewFloat("13591409")
	M := NewFloat("1")
	L := NewFloat("13591409")
	X := NewFloat("1")
	S := NewFloat("13591409") // same as L

	// Multiply C by sqrt(10005)
	C.Times(NewFloat("10005").SquareRoot().Val())

	// Set factorials
	factK := NewInt("1")

	// Iterations to improve precision (increase this for more digits)
	iterations := 50

	// Use Chudnovsky's series to compute Pi
	for i := 1; i < iterations; i++ {
		// Update K, M, L, X
		K.Times("545140134")
		M.Times(strconv.Itoa(i * i * i))

		L.Plus("545140134")
		X.Times("-262537412640768000")

		// Calculate term and add to the sum S
		term := NewFloat(M.Val()).Times(L.Val()).DividedBy(X.Val()).DividedBy(factK.Val())
		S.Plus(term.Val())

		// Update the factorial
		factK.Times(strconv.Itoa(6 * i))
	}

	// Calculate pi = C / S
	return NewFloat(C.Val()).DividedBy(S.Val())
}

func (f Float) Mod2Pi() Float {
	// Use a more precise Pi value for accurate modulo calculations
	// Perform x mod 2*Pi
	return f.Mod(Pi().Times("2").Val())
}
