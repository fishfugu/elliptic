package ellipticcurve

import (
	"fmt"
	"math/big"
	"os"

	"elliptic/pkg/utils"

	"github.com/sirupsen/logrus"
)

func init() {
	// log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// output to stdout instead of the default stderr.
	logrus.SetOutput(os.Stdout)

	// log the debug severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}

// NOTE: these are private / immutable on purpose
// All ECs defined in the Weierstrass form

// EllipticCurve represents an elliptic curve defined by the equation y^2 = x^3 + Ax + B
// where A and B are both integers represented *big.Int objects
type EllipticCurve struct {
	a, b *big.Int // Coefficients of the curve equation.
}

// FiniteFieldEC represents an elliptic curve over a finite field defined by the equation y^2 = x^3 + Ax + B.
type FiniteFieldEC struct {
	ec EllipticCurve
	// TODO: should this be definable by "strings" rather than *big.Int?
	// Woould that make it easier to interface with bigarith functions?
	p *big.Int // Coefficients of the curve equation and prime modulus of the field.
}

// NewEllipticCurve creates a new EllipticCurve with given coefficients.
func NewEllipticCurve(a, b *big.Int) *EllipticCurve {
	return &EllipticCurve{a: a, b: b}
}

// NewFiniteFieldEC creates a new EllipticCurve, defined over a finite field, with given coefficients and modulus.
func NewFiniteFieldEC(a, b, p *big.Int) *FiniteFieldEC {
	aModP := new(big.Int).Mod(a, p)
	bModP := new(big.Int).Mod(b, p)
	EC := NewEllipticCurve(aModP, bModP)
	return &FiniteFieldEC{ec: *EC, p: p}
}

// GetDetails returns the coefficients A and B of the curve.
func (ec *EllipticCurve) GetDetails() (*big.Int, *big.Int) {
	return ec.a, ec.b
}

// GetDetailsAsRats returns the coefficients A and B of the curve, as big.Rat values.
func (ec *EllipticCurve) GetDetailsAsRats() (*big.Rat, *big.Rat) {
	A := new(big.Rat).SetInt(ec.a)
	B := new(big.Rat).SetInt(ec.b)
	return A, B
}

// GetDetails returns the coefficients A, B, and the modulus P of the finite field curve.
func (ffec *FiniteFieldEC) GetDetails() (*big.Int, *big.Int, *big.Int) {
	return ffec.ec.a, ffec.ec.b, ffec.p
}

// GetDetailsAsRats returns the coefficients A, B, and the modulus P of the finite field curve, as big.Rat values.
func (ffec *FiniteFieldEC) GetDetailsAsRats() (*big.Rat, *big.Rat, *big.Rat) {
	A := new(big.Rat).SetInt(ffec.ec.a)
	B := new(big.Rat).SetInt(ffec.ec.b)
	P := new(big.Rat).SetInt(ffec.p)
	return A, B, P
}

func pi() *big.Rat {
	// Approximation - 2π as a big.Rat
	piNumer, _ := new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951941511609433057270365759591953092186117381932611793105118548074462379962749567351885752724891227938183011949129833673362", 10)
	piDenom, _ := new(big.Int).SetString("1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 10)
	return new(big.Rat).SetFrac(piNumer, piDenom)
}

func twoPiRat() *big.Rat {
	// Approximation - 2π as a big.Rat
	return new(big.Rat).Mul(pi(), big.NewRat(2, 1))
}

func piOver2Rat() *big.Rat {
	// π/2 as a rational approximation
	return new(big.Rat).Quo(pi(), big.NewRat(2, 1))
}

// PowRatInt computes a^b - big.Rat a and big.Int b
func powRatInt(base *big.Rat, exp *big.Int) *big.Rat {
	result := new(big.Rat).SetInt64(1) // Start with 1 as the identity - multiplication

	if exp.Sign() == 0 {
		// Any number to the power of 0 is 1
		return result
	}

	absExp := new(big.Int).Abs(exp) // Get the absolute value of the exponent

	tempBase := new(big.Rat).Set(base) // Copy the base to avoid modifying the input
	for absExp.BitLen() > 0 {
		if absExp.Bit(0) == 1 {
			result.Mul(result, tempBase)
		}
		tempBase.Mul(tempBase, tempBase)
		absExp.Rsh(absExp, 1)
	}

	if exp.Sign() < 0 {
		// If the exponent is negative, invert the result
		return new(big.Rat).Inv(result)
	}

	return result
}

// modRatInt computes `a mod b` where `a` is a big.Rat and `b` is a big.Int.
// The result is a big.Rat such that 0 <= result < b (converted to a big.Rat).
func modRatInt(a *big.Rat, b *big.Int) *big.Rat {
	// Convert b to big.Rat
	bRat := new(big.Rat).SetInt(b)

	// Compute the integer division: q = floor(a / b)
	quotient := new(big.Rat).Quo(a, bRat)
	quotientFloor := new(big.Rat).SetFrac(quotient.Num(), quotient.Denom())
	quotientFloor.SetInt(quotientFloor.Num().Quo(quotientFloor.Num(), quotientFloor.Denom()))

	// Compute the remainder: r = a - q * b
	remainder := new(big.Rat).Sub(a, new(big.Rat).Mul(quotientFloor, bRat))

	// Ensure the result is in the range [0, b)
	if remainder.Cmp(big.NewRat(0, 1)) < 0 {
		remainder.Add(remainder, bRat)
	}

	return remainder
}

// modRat computes `a mod b` - two big.Rat values.
// The result is a big.Rat such that 0 <= result < b.
func modRat(a, b *big.Rat) *big.Rat {
	// Compute the integer division: q = floor(a / b)
	quotient := new(big.Rat).Quo(a, b)
	quotientFloor := new(big.Rat).SetFrac(quotient.Num(), quotient.Denom())
	quotientFloor.SetInt(quotientFloor.Num().Quo(quotientFloor.Num(), quotientFloor.Denom()))

	// Compute the remainder: r = a - q * b
	remainder := new(big.Rat).Sub(a, new(big.Rat).Mul(quotientFloor, b))

	// Ensure the result is in the range [0, b)
	if remainder.Cmp(big.NewRat(0, 1)) < 0 {
		remainder.Add(remainder, b)
	}

	return remainder
}

// sqrtRat computes the square root of a big.Rat and returns another big.Rat.
// If the result is not an exact rational number, it returns an approximation.
func sqrtRat(input *big.Rat) (*big.Rat, error) {
	// Separate numerator and denominator
	num := input.Num()   // Numerator
	den := input.Denom() // Denominator

	// Check if numerator and denominator are perfect squares
	sqrtNum := intSqrt(num)
	sqrtDen := intSqrt(den)

	if sqrtNum != nil && sqrtDen != nil {
		// Exact rational square root
		return new(big.Rat).SetFrac(sqrtNum, sqrtDen), nil
	}

	// Approximation using big.Float
	floatInput := new(big.Float).SetRat(input)
	floatSqrt := new(big.Float).Sqrt(floatInput)

	// Convert the result back to big.Rat
	result := new(big.Rat)
	if _, acc := floatSqrt.Float64(); acc == big.Exact {
		floatSqrt.Rat(result)
	} else {
		floatSqrt64, _ := floatSqrt.Float64()
		result.SetFloat64(floatSqrt64)
	}

	return result, nil
}

// intSqrt computes the integer square root of a big.Int if it exists.
// Returns nil if the number is not a perfect square.
func intSqrt(n *big.Int) *big.Int {
	root := new(big.Int).Sqrt(n)
	square := new(big.Int).Mul(root, root)
	if square.Cmp(n) == 0 {
		return root
	}
	return nil
}

// cbrtRat computes the cube root of a big.Rat and returns another big.Rat.
// If the result is not an exact rational number, it returns an approximation.
func cbrtRat(input *big.Rat) (*big.Rat, error) {
	// Separate numerator and denominator
	num := input.Num()   // Numerator
	den := input.Denom() // Denominator

	// Check if numerator and denominator are perfect cubes
	cbrtNum := intCbrt(num)
	cbrtDen := intCbrt(den)

	if cbrtNum != nil && cbrtDen != nil {
		// Exact rational cube root
		return new(big.Rat).SetFrac(cbrtNum, cbrtDen), nil
	}

	// Approximation using big.Float
	floatInput := new(big.Float).SetRat(input)
	floatCbrt := cbrtFloat(floatInput)

	// Convert the result back to big.Rat
	result := new(big.Rat)
	floatCbrt.Rat(result)
	return result, nil
}

// cbrtFloat approximates the cube root of a big.Float using Newton's method.
func cbrtFloat(x *big.Float) *big.Float {
	zeroFloat := new(big.Float).SetInt64(0)
	twoFloat := new(big.Float).SetInt64(2)
	threeFloat := new(big.Float).SetInt64(3) // Initial guess: x^(1/3) ≈ x^(1/2) - but use Abs Val of operand

	guess := new(big.Float).Sqrt(new(big.Float).Abs(x))
	guess = new(big.Float).Quo(guess, twoFloat)
	if new(big.Float).Set(x).Cmp(zeroFloat) < 0 { // if operand is neg - make neg again
		guess = new(big.Float).Neg(guess)
	}

	// Newton's method iteration: x_{n+1} = (2*x_n + a / (x_n^2)) / 3
	temp := new(big.Float).SetInt64(0)
	for i := 0; i < int(100000); i++ {
		guessSquared := new(big.Float).Mul(guess, guess)
		twoTimesGuess := new(big.Float).Mul(guess, twoFloat)

		temp = new(big.Float).Quo(x, guessSquared) // temp = x / guess^2

		twoTimesGuessPLusTemp := new(big.Float).Add(twoTimesGuess, temp)

		guess = new(big.Float).Quo(twoTimesGuessPLusTemp, threeFloat) // guess = (2*guess + temp) / 3
	}
	return guess
}

// intCbrt computes the integer cube root of a big.Int if it is a perfect cube.
// Otherwise, it returns nil.
func intCbrt(x *big.Int) *big.Int {
	low := big.NewInt(0)
	high := new(big.Int).Set(x)
	mid := new(big.Int)
	cube := new(big.Int)

	for low.Cmp(high) <= 0 {
		mid.Add(low, high).Rsh(mid, 1)    // mid = (low + high) / 2
		cube.Exp(mid, big.NewInt(3), nil) // cube = mid^3

		cmp := cube.Cmp(x)
		if cmp == 0 {
			return mid // Perfect cube found
		} else if cmp < 0 {
			low.Add(mid, big.NewInt(1))
		} else {
			high.Sub(mid, big.NewInt(1))
		}
	}
	return nil
}

// cosRat calculates the cosine of a `big.Rat` using a Taylor series expansion.
// It returns the result as a `big.Rat`.
func cosRat(input *big.Rat) (*big.Rat, error) {
	// Reduce input modulo 2π
	inputMod := modRat(input, twoPiRat())

	// Initialize result: cos(0) = 1
	result := new(big.Rat).SetInt64(1)

	// Compute terms of the Taylor series
	xPower := big.NewRat(1, 1) // x^(2n), starts as 1
	factorial := big.NewInt(1) // (2n)!
	sign := big.NewRat(1, 1)   // Alternating sign, starts as +1

	for n := 1; n < 100000; n++ {
		// Update xPower: xPower *= inputMod^2
		xPower.Mul(xPower, inputMod)
		xPower.Mul(xPower, inputMod)

		// Update factorial: factorial *= (2n-1)*(2n)
		factorial.Mul(factorial, big.NewInt(int64(2*n-1)))
		factorial.Mul(factorial, big.NewInt(int64(2*n)))

		// Compute term: term = xPower / factorial
		term := new(big.Rat).Set(xPower)
		term.Quo(term, new(big.Rat).SetInt(factorial))

		// Alternate the sign
		sign.Neg(sign)
		term.Mul(term, sign)

		// Add the term to the result
		result.Add(result, term)
	}

	return result, nil
}

// arcCosRat computes an approximation of arcCos(x) using a Taylor series expansion.
// Input and output are `big.Rat`.
func arcCosRat(input *big.Rat) (*big.Rat, error) {
	// Check domain of input
	if input.Cmp(big.NewRat(-1, 1)) < 0 || input.Cmp(big.NewRat(1, 1)) > 0 {
		return nil, fmt.Errorf("input out of domain - arcCos: %s", input.String())
	}

	// Initialize the result with π/2
	result := new(big.Rat).Set(piOver2Rat())

	// Taylor series computation
	xPower := new(big.Rat).Set(input) // x^(2n+1), starts as x
	factorial := big.NewInt(1)        // (2n)!
	twoPower := big.NewInt(1)         // 2^(2n)
	coeff := big.NewInt(1)            // Binomial coefficient
	//sign := big.NewInt(-1)            // Alternating sign

	// number of loops == number of terms in Taylr series
	for n := 0; n < 10000; n++ {
		// Compute term: (2n)! / (2^(2n) * (n!)^2 * (2n+1))
		term := new(big.Rat).SetInt(factorial)
		term.Quo(term, new(big.Rat).SetInt(twoPower))                 // Divide by 2^(2n)
		term.Quo(term, new(big.Rat).SetInt(coeff))                    // Divide by (n!)^2
		term.Quo(term, new(big.Rat).SetInt(big.NewInt(int64(2*n+1)))) // Divide by (2n+1)
		term.Mul(term, xPower)                                        // Multiply by x^(2n+1)

		// Add or subtract the term
		if n%2 == 0 {
			result.Sub(result, term) // Subtract - even n
		} else {
			result.Add(result, term) // Add - odd n
		}

		// Update xPower - the next iteration: xPower *= x^2
		xPower.Mul(xPower, input).Mul(xPower, input)

		// Update factorial - the next iteration: factorial *= (2n+2)*(2n+3)
		factorial.Mul(factorial, big.NewInt(int64(2*n+2)))
		factorial.Mul(factorial, big.NewInt(int64(2*n+3)))

		// Update twoPower: twoPower *= 4
		twoPower.Mul(twoPower, big.NewInt(4))

		// Update coeff: coeff *= (n+1)*(n+1)
		coeff.Mul(coeff, big.NewInt(int64(n+1)))
	}

	return result, nil
}

// finds minimum value of X - an Elliptic Curve where y = 0
// curve in Weierstrass form, this should be the
// lowest value of x - the whole curve in the real numbers

// solveCubic - form x^3 + Ax + B - real roots only
func (ec EllipticCurve) SolveCubic() ([]*big.Rat, error) {
	var roots []*big.Rat

	logger := utils.InitialiseLogger("[SolveCubic]")

	twoRat, ok := new(big.Rat).SetString("2")
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'two' failed")
	threeRat, ok := new(big.Rat).SetString("3")
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'two' failed")
	threeInt, ok := new(big.Int).SetString("3", 10)
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'two' failed")

	A, B := ec.GetDetailsAsRats()

	// convert A and B to big.Rat

	// Calculate the discriminant - (A/3)^3 + (B/2)^2
	// = (A^3/27) + (B^2/4) = (4 A^3 / 108) + (27 B^2 / 108) = (4 A^3 + 27 B^2) / 108
	ZeroRat := new(big.Rat).SetInt64(0)

	NegativeB := new(big.Rat).Neg(B)
	NegativeBOver2 := new(big.Rat).Quo(NegativeB, twoRat)
	AOver3 := new(big.Rat).Quo(A, threeRat)
	Aover3Cubed := powRatInt(AOver3, threeInt)
	BOver2 := new(big.Rat).Quo(B, twoRat)
	BOver2Squared := new(big.Rat).Mul(BOver2, BOver2)

	discriminant := new(big.Rat).Add(Aover3Cubed, BOver2Squared)

	discriminantCmpToZero := discriminant.Cmp(ZeroRat)
	if discriminantCmpToZero > 0 {
		// One real root, two complex roots
		sqrtDiscriminant, err := sqrtRat(discriminant)
		if err != nil {
			return nil, err
		}
		NegativeBOver2PlusSqrtDiscriminant := new(big.Rat).Add(NegativeBOver2, sqrtDiscriminant)
		u, err := cbrtRat(NegativeBOver2PlusSqrtDiscriminant)
		if err != nil {
			return nil, err
		}

		NegativeBOver2MinusSqrtDiscriminant := new(big.Rat).Sub(NegativeBOver2, sqrtDiscriminant)
		v, err := cbrtRat(NegativeBOver2MinusSqrtDiscriminant)
		if err != nil {
			return nil, err
		}

		root := new(big.Rat).Add(u, v)
		roots = append(roots, root)
	} else if discriminantCmpToZero == 0 {
		// All roots are real, at least two are equal
		u, err := cbrtRat(NegativeBOver2)
		if err != nil {
			return nil, err
		}
		root1 := new(big.Rat).Mul(u, twoRat)
		root2 := new(big.Rat).Neg(u)
		roots = append(roots, root1, root2, root2)
	} else {
		// Three real roots (discriminant < 0)
		fmt.Printf("discriminant: %s", discriminant.String())

		// r = \sqrt{\frac{-A^3}{27}}
		aSquared := new(big.Rat).Mul(A, A)
		aCubed := new(big.Rat).Mul(aSquared, A)
		negACubed := new(big.Rat).Neg(aCubed)
		twentySeven, _ := new(big.Rat).SetString("27")
		negACubedOver27 := new(big.Rat).Quo(negACubed, twentySeven)
		r, err := sqrtRat(negACubedOver27)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("r: %s", r.String())

		// \cos(\theta) = -\frac{B}{2r}
		// \theta) = \arccos{ -\frac{B}{2r} }
		twoR := new(big.Rat).Mul(twoRat, r)
		bOverTwoR := new(big.Rat).Quo(B, twoR)
		negBOverTwoR := new(big.Rat).Neg(bOverTwoR)

		theta, err := arcCosRat(negBOverTwoR)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("theta: %s", theta.String())

		// all the final root values are calculated by multiplying by: 2 * sqrt{ - A / 3 }

		negAOver3 := new(big.Rat).Neg(AOver3)
		sqrtNegAOVer3, err := sqrtRat(negAOver3)
		if err != nil {
			return nil, err
		}
		sqrtMinusAOver3Times2 := new(big.Rat).Mul(sqrtNegAOVer3, twoRat)
		// (theta plus 2kpi) / 3

		thetaOver3 := new(big.Rat).Quo(theta, threeRat)
		cosThetaOver3, err := cosRat(thetaOver3)

		thetaPlus2Pi := new(big.Rat).Add(theta, twoPiRat())
		thetaPlus2PiOver3 := new(big.Rat).Quo(thetaPlus2Pi, threeRat)
		cosThetaPlus2PiOver3, err := cosRat(thetaPlus2PiOver3)
		if err != nil {
			return nil, err
		}

		fourPiRat := new(big.Rat).Mul(twoPiRat(), twoRat)
		thetaPlus4Pi := new(big.Rat).Add(theta, fourPiRat)
		thetaPlus4PiOver3 := new(big.Rat).Quo(thetaPlus4Pi, threeRat)
		cosThetaPlus4PiOver3, err := cosRat(thetaPlus4PiOver3)
		if err != nil {
			return nil, err
		}

		// The three real roots \(x_1\), \(x_2\), and \(x_3\) can then be computed as:
		// \[
		// x_k = 2\sqrt{\frac{-A}{3}} \cos\left( \frac{\theta + 2k\pi}{3} \right)
		// \]
		// where \(k = 0, 1, 2\) fr the three distinct roots.
		// 	return roots, nil
		// }

		root1 := new(big.Rat).Mul(sqrtMinusAOver3Times2, cosThetaOver3)
		root2 := new(big.Rat).Mul(sqrtMinusAOver3Times2, cosThetaPlus2PiOver3)
		root3 := new(big.Rat).Mul(sqrtMinusAOver3Times2, cosThetaPlus4PiOver3)

		roots = append(roots, root1, root2, root3)
		logrus.Debug("Roots calculated - all real and all distinct roots: ", roots)
		logrus.Debugf("Rootz: %s", roots)
	}

	return roots, nil
}

func (ffec FiniteFieldEC) SolveCubic() ([]*big.Rat, error) {
	results, err := ffec.ec.SolveCubic()
	// convert each value into its mod p equivalent
	for i, result := range results {
		results[i] = modRatInt(result, ffec.p)
	}
	return results, err
}

// FindY finds the y value - an EllipticCurve - x^3 + Ax + B
// it returns the positive y value - but the other value is simply the negative of that anyway
func (ec EllipticCurve) FindY(x *big.Rat) (*big.Rat, error) {
	A, B := ec.GetDetails()
	ARat := new(big.Rat).SetInt(A)
	BRat := new(big.Rat).SetInt(B)
	Ax := new(big.Rat).Mul(ARat, x)
	xSquared := new(big.Rat).Mul(x, x)
	xCubed := new(big.Rat).Mul(xSquared, x)
	xCubedPlusAx := new(big.Rat).Add(xCubed, Ax)
	xCubedPlusAxPlusB := new(big.Rat).Add(xCubedPlusAx, BRat)
	sqrtXCubedPlusAxPlusB, err := sqrtRat(xCubedPlusAxPlusB)
	if err != nil {
		return nil, err
	}
	return sqrtXCubedPlusAxPlusB, nil
}

func (ffec FiniteFieldEC) FindY(x *big.Rat) (*big.Rat, error) {
	result, err := ffec.ec.FindY(x)
	if err != nil {
		return nil, err
	}
	return modRatInt(result, ffec.p), nil
}
