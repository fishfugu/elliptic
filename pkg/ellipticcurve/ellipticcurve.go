package ellipticcurve

import (
	"fmt"
	"math/big"

	"elliptic/pkg/utils"
)

var (
	zeroInt, oneInt, twoInt *big.Int
	// piRat, twoPiRat, fourPiRat, halfPiRat      *big.Rat
	zeroRat, twoRat, threeRat *big.Rat
	precision, tolerance      uint
	toleranceFractionRat      *big.Rat
)

func init() {
	// Approximation - 2π as a big.Rat
	// piNumer, _ := new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951941511609433057270365759591953092186117381932611793105118548074462379962749567351885752724891227938183011949129833673362", 10)
	// piDenom, _ := new(big.Int).SetString("1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 10)

	zeroInt = new(big.Int).SetInt64(0)
	oneInt = new(big.Int).SetInt64(1)
	twoInt = new(big.Int).SetInt64(2)

	zeroRat = big.NewRat(0, 1)
	// oneRat = big.NewRat(1, 1)
	twoRat = big.NewRat(2, 1)
	threeRat = big.NewRat(3, 1)
	// fourRat = big.NewRat(4, 1)

	// piRat = new(big.Rat).SetFrac(piNumer, piDenom)
	// halfPiRat = new(big.Rat).Quo(piRat, big.NewRat(2, 1))
	// twoPiRat = new(big.Rat).Mul(piRat, big.NewRat(2, 1))
	// fourPiRat = new(big.Rat).Mul(piRat, big.NewRat(2, 1))

	// 2048 bit precision was chosen as approximately 4 times the MAXIMUM number of bits
	// used for. keys in EC Cryptography in the highest security situations:
	// https://en.wikipedia.org/wiki/Key_size?utm_source=chatgpt.com
	// "521-bit keys: Deliver a security level of roughly 256 bits, used in scenarios requiring
	// the highest security assurances."
	precision = 2048
	tolerance = 1024

	// set approximation preciion to 1024 - half of calc precision, still almost double the usual max
	// precision used for. keys in EC Cryptography in the highest security situations:
	// https://en.wikipedia.org/wiki/Key_size?utm_source=chatgpt.com
	tolleranceInt := new(big.Int).SetUint64(uint64(tolerance))
	// Check if input is within tolerance of a whole number
	twoToThePowerOfTollerance := new(big.Int).Exp(twoInt, tolleranceInt, nil)
	toleranceFractionRat = new(big.Rat).SetFrac(oneInt, twoToThePowerOfTollerance)
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
	return &EllipticCurve{
		a: a,
		b: b,
	}
}

// NewFiniteFieldEC creates a new EllipticCurve, defined over a finite field, with given coefficients and modulus.
func NewFiniteFieldEC(a, b, p *big.Int) *FiniteFieldEC {
	return &FiniteFieldEC{
		ec: *NewEllipticCurve(
			new(big.Int).Mod(a, p),
			new(big.Int).Mod(b, p),
		),
		p: p,
	}
}

// GetDetails returns the coefficients A and B of the curve.
func (ec *EllipticCurve) GetDetails() (*big.Int, *big.Int) {
	return ec.a, ec.b
}

// GetDetailsAsRats returns the coefficients A and B of the curve, as big.Rat values.
func (ec *EllipticCurve) GetDetailsAsRats() (*big.Rat, *big.Rat) {
	return new(big.Rat).SetInt(ec.a), new(big.Rat).SetInt(ec.b)
}

// GetDetails returns the coefficients A, B, and the modulus P of the finite field curve.
func (ffec *FiniteFieldEC) GetDetails() (*big.Int, *big.Int, *big.Int) {
	return ffec.ec.a, ffec.ec.b, ffec.p
}

// GetDetailsAsRats returns the coefficients A, B, and the modulus P of the finite field curve, as big.Rat values.
func (ffec *FiniteFieldEC) GetDetailsAsRats() (*big.Rat, *big.Rat, *big.Rat) {
	return new(big.Rat).SetInt(ffec.ec.a), new(big.Rat).SetInt(ffec.ec.b), new(big.Rat).SetInt(ffec.p)
}

// solveCubic - form x^3 + Ax + B
func (ec EllipticCurve) SolveCubic() ([]*big.Rat, error) {
	logger := utils.InitialiseLogger("[SolveCubic]")
	var roots []*big.Rat

	A, B := ec.GetDetailsAsRats()

	// convert A and B to big.Rat

	// Discriminant
	// 4a^3 + 27b^2
	aCubed := new(big.Rat).Mul(A, new(big.Rat).Mul(A, A))
	bSquared := new(big.Rat).Mul(B, B)
	fourACubed := new(big.Rat).Mul(new(big.Rat).SetInt64(4), aCubed)
	twentySevenBCubed := new(big.Rat).Mul(new(big.Rat).SetInt64(27), bSquared)
	discriminant := new(big.Rat).Add(fourACubed, twentySevenBCubed)

	logger.Debugf("Discriminant: %s", discriminant.FloatString(10))

	// Find one root using newtonCubic
	root1, err := newtonCubic(A, B)
	if err != nil {
		utils.LogOnError(logger, err, "SolveCubic/newtonCubic", false)
		return nil, err
	}
	roots = append(roots, approximateRat(root1))

	logger.Debugf("First root: %s", approximateRat(root1).FloatString(400))

	if discriminant.Cmp(zeroRat) == 0 {
		// must have a double root
		// see if root1 is double
		xSquared := new(big.Rat).Mul(root1, root1)
		threeXSquared := new(big.Rat).Mul(threeRat, xSquared)
		gradient := approximateRat(new(big.Rat).Add(threeXSquared, A))

		logger.Debugf("Gradient: %s", gradient.FloatString(400))
		logger.Debugf("Cmp: %d", gradient.Cmp(zeroRat))

		if gradient.Cmp(zeroRat) == 0 {
			// root2 == root1 - AND root3 == - 2 * root1
			roots = append(roots, approximateRat(root1))
			root3 := new(big.Rat).Neg(new(big.Rat).Mul(root1, twoRat))
			roots = append(roots, approximateRat(root3))
		} else {
			// root2 == root3 == - root1 / 2
			root2 := new(big.Rat).Neg(new(big.Rat).Quo(root1, twoRat))
			roots = append(roots, approximateRat(root2))
			roots = append(roots, approximateRat(root2))
		}
	}

	if discriminant.Cmp(zeroRat) < 0 {
		// must have 3 distinct roots
		// Find the remaining two roots
		remainingRoots, err := findRemainingRoots(A, root1)
		if err != nil {
			utils.LogOnError(logger, err, "SolveCubic/findRemainingRoots", false)
			return nil, err
		}

		logger.Debugf("Remianing roots: 1) %s, 2) %s", remainingRoots[0].FloatString(10), remainingRoots[1].FloatString(10))
		for i, root := range remainingRoots {
			remainingRoots[i] = approximateRat(root)
		}
		roots = append(roots, remainingRoots...)
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
	A, B := ec.GetDetailsAsRats()
	Ax := new(big.Rat).Mul(A, x)
	xSquared := new(big.Rat).Mul(x, x)
	xCubed := new(big.Rat).Mul(xSquared, x)
	xCubedPlusAx := new(big.Rat).Add(xCubed, Ax)
	xCubedPlusAxPlusB := new(big.Rat).Add(xCubedPlusAx, B)
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

// sqrtRat computes the square root of a big.Rat with arbitrary precision.
// If the result is not an exact rational number, it computes an approximation with the specified precision.
func sqrtRat(input *big.Rat) (*big.Rat, error) {
	logger := utils.InitialiseLogger("[sqrtRat]")
	logger.Debug("starting function sqrtRat")

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

	// Use Newton's method for. arbitrary precision square root
	// see init for ideas behind precision level
	floatInput := new(big.Float).SetPrec(precision).SetRat(input)
	floatSqrt := sqrtFloat(floatInput, precision)

	// Convert the result back to big.Rat
	result := new(big.Rat)
	floatSqrt.Rat(result)
	return result, nil
}

// sqrtFloat computes the square root of a big.Float using Newton's method with the specified precision.
func sqrtFloat(a *big.Float, prec uint) *big.Float {
	logger := utils.InitialiseLogger("[sqrtFloat]")
	logger.Debug("starting function sqrtFloat")

	// Initial guess: x0 = a / 2
	guess := utils.NewFloat().Quo(a, big.NewFloat(2))

	// Iteratively refine the guess
	half := utils.NewFloat().SetFloat64(0.5)
	for i := uint(0); i < prec; i++ {
		temp := utils.NewFloat().Quo(a, guess)  // a / guess
		guess.Add(guess, temp).Mul(guess, half) // (guess + a/guess) / 2
	}

	return guess
}

// intSqrt computes the integer square root of a big.Int if it is a perfect square.
// Otherwise, it returns nil.
func intSqrt(x *big.Int) *big.Int {
	logger := utils.InitialiseLogger("[intSqrt]")
	logger.Debug("starting function intSqrt")

	// Use binary search to find the integer square root
	low := big.NewInt(0)
	high := new(big.Int).Set(x)
	mid := new(big.Int)
	square := new(big.Int)

	for low.Cmp(high) <= 0 {
		logger.Debug(fmt.Sprintf("Low: %s, High: %s", low.String(), high.String()))
		mid = new(big.Int).Rsh(new(big.Int).Add(low, high), 1) // mid = (low + high) / 2
		logger.Debug(fmt.Sprintf("Mid: %s:", mid.String()))
		square = new(big.Int).Mul(mid, mid) // square = mid^2

		cmp := square.Cmp(x)
		logger.Debug(fmt.Sprintf("X: %s, Square: %s", x.String(), square.String()))

		if cmp == 0 {
			return mid // Perfect square
		} else if cmp < 0 {
			low.Add(mid, big.NewInt(1))
		} else {
			high.Sub(mid, big.NewInt(1))
		}
	}

	return nil // Not a perfect square
}

// QuickEstimateRoot estimates a root of the cubic equation x^3 + Ax + B = 0
// by solving the linear approximation Ax + B = 0
// Using the linear approximation y = Ax + B
// to estimate the root is a quick and efficient heuristic. By finding where
// Ax+B=0, the x-intercept of this linear equation provides a very reasonable
// starting point for one of the roots of the cubic.
func quickEstimateRoot(A, B *big.Rat) *big.Rat {
	// Check for B == 0 which means f(x) = x^3 + Ax and so x == 0 is a root
	if B.Sign() == 0 {
		return new(big.Rat).SetInt64(0)
	}
	negB := new(big.Rat).Neg(B)
	// Check for division by zero (A = 0)
	// Reutrn -B as estimate
	if A.Sign() == 0 {
		return negB
	}
	// Compute -B / A
	return new(big.Rat).Quo(negB, A)
}

// newtonCubic finds one root for the cubic in the form x^3 + Ax + B
func newtonCubic(A, B *big.Rat) (*big.Rat, error) {
	logger := utils.InitialiseLogger("[newtonCubic]")
	logger.Info("starting function newtonCubic")

	// for x^3 + Ax + B we have r_1 + r_2 + r_3 = 0
	// so in some sense one root must be >= 0 and another must be <= 0
	// so 0 may not be a bad initial guess
	// let's try this instead:
	initialGuess := new(big.Rat).SetInt64(0)
	var err error
	if A.Sign() != 0 {
		initialGuess = quickEstimateRoot(A, B)
		if err != nil {
			utils.LogOnError(logger, err, "newtonCubic", true)
		}
	}
	x := new(big.Rat).Set(initialGuess)
	delta := new(big.Rat).SetInt64(1)
	for {
		logger.Debugf("Starting Loop for x: %s", x.FloatString(1000))

		// f(x) = x^3 + Ax + B
		fx := new(big.Rat).Add(
			new(big.Rat).Add(
				new(big.Rat).Mul(x, new(big.Rat).Mul(x, x)), // x^3
				new(big.Rat).Mul(A, x),                      // Ax
			),
			B,
		) // x^3 + Ax + B
		logger.Debugf("FX: %s", fx.FloatString(1000))

		// f'(x) = 3x^2 + A
		fpx := new(big.Rat).Add(
			new(big.Rat).Mul(threeRat, new(big.Rat).Mul(x, x)), // 3x^2
			A,
		) // 3x^2 + A
		logger.Debugf("FPX: %s", fpx.FloatString(1000))

		if fpx.Cmp(zeroRat) != 0 {
			delta = new(big.Rat).Quo(fx, fpx) // f(x) / f'(x)
			logger.Debugf("Delta: %s", delta.FloatString(1000))
		}

		x.Sub(x, delta) // x = x - (f(x) / f'(x))
		logger.Debugf("X: %s", x.FloatString(1000))

		// now check for Rat simplification
		x = approximateRat(x)

		deltaAbs := new(big.Rat).Abs(delta)
		logger.Debugf("1 DeltaAbs: %s, toleranceFractionRat: %s", deltaAbs.FloatString(1000), toleranceFractionRat.FloatString(1000))

		if deltaAbs.Cmp(toleranceFractionRat) == -1 { // Convergence tolerance
			break
		}
		logger.Debugf("2 DeltaAbs: %s, toleranceFractionRat: %s", deltaAbs.FloatString(1000), toleranceFractionRat.FloatString(1000))
	}
	return x, nil
}

// solveQuadratic calculates the roots of a quadratic equation of the form
// x^2 + px + q = 0 and returns the two roots as big.Rat.
func solveQuadratic(p, q *big.Rat) ([]*big.Rat, error) {
	logger := utils.InitialiseLogger("[solveQuadratic]")
	logger.Info("starting function solveQuadratic")

	// Discriminant: D = p^2 - 4q
	pSquared := new(big.Rat).Mul(p, p)                // p^2
	fourQ := new(big.Rat).Mul(big.NewRat(4, 1), q)    // 4q
	discriminant := new(big.Rat).Sub(pSquared, fourQ) // D = p^2 - 4q

	logger.Debugf("discriminant: %s, p: %s, q: %s", discriminant.FloatString(400), p.FloatString(400), q.FloatString(400))

	// Check if discriminant is negative (no real roots)
	if discriminant.Sign() < 0 {
		return nil, fmt.Errorf("no real roots; discriminant is negative - p: %s, q: %s", p.FloatString(10), q.FloatString(10))
	}

	// Square root of the discriminant
	sqrtDiscriminant, err := sqrtRat(discriminant)
	if err != nil {
		return nil, fmt.Errorf("failed to compute square root: %v", err)
	}

	// Roots: x = (-p ± sqrt(D)) / 2
	negP := new(big.Rat).Neg(p)                       // -p
	root1 := new(big.Rat).Add(negP, sqrtDiscriminant) // -p + sqrt(D)
	root1.Quo(root1, twoRat)                          // (-p + sqrt(D)) / 2

	root2 := new(big.Rat).Sub(negP, sqrtDiscriminant) // -p - sqrt(D)
	root2.Quo(root2, twoRat)                          // (-p - sqrt(D)) / 2

	return []*big.Rat{root1, root2}, nil
}

// findRemainingRoots calculates the remaining two roots of a cubic equation
// given one known root.
func findRemainingRoots(A, root1 *big.Rat) ([]*big.Rat, error) {
	// Calculate the quadratic coefficients
	// p = root1, q = A + root1^2
	root1Squared := new(big.Rat).Mul(root1, root1) // root1^2
	p := root1                                     // p = root1
	q := new(big.Rat).Add(A, root1Squared)

	// Solve the quadratic equation
	return solveQuadratic(p, q)
}

// ApproximateRat approximates a big.Rat to the nearest integer or simpler rational number within a given tolerance.
// ApproximateRat is a concenssion to the fact that we HAVE to estimate any imperfect square roots / cube roots
// because they are irrational numbers... but we want to use arbitrarily precise Rationals wherever poss
// for. all calculations
// and so ApproximateRat is ONLY applied at the version last moments - to a specifically chosen precision
// see comments where it's applied as to how that precisions was estimated
// TODO: build a function that can find the minimum required precision that passes all tests used
func approximateRat(input *big.Rat) *big.Rat {
	logger := utils.InitialiseLogger("[ApproximateRat]")
	logger.Debug("starting function ApproximateRat")

	// see init for details of tolerance calculations
	nearestInt := new(big.Int).Div(input.Num(), input.Denom())
	nearestIntRat := new(big.Rat).SetInt(nearestInt)
	diff := new(big.Rat).Sub(input, nearestIntRat)
	if diff.Abs(diff).Cmp(toleranceFractionRat) <= 0 {
		return nearestIntRat
	}

	// If not near an integer, find the best rational approximation
	return bestRationalApproximation(input, toleranceFractionRat)
}

// bestRationalApproximation finds the best rational approximation of input within the given tolerance.
func bestRationalApproximation(input *big.Rat, tolerance *big.Rat) *big.Rat {
	logger := utils.InitialiseLogger("[bestRationalApproximation]")
	logger.Debug("starting function bestRationalApproximation")

	// Continued fraction expansion to find the best rational approximation
	// Initialize variables
	a := new(big.Int)
	num0 := big.NewInt(0)
	num1 := big.NewInt(1)
	den0 := big.NewInt(1)
	den1 := big.NewInt(0)
	approx := new(big.Rat)

	x := new(big.Rat).Set(input)
	for {
		// a = floor(x)
		a.Div(x.Num(), x.Denom())

		// Update numerator and denominator
		num2 := new(big.Int).Add(new(big.Int).Mul(a, num1), num0)
		den2 := new(big.Int).Add(new(big.Int).Mul(a, den1), den0)

		// Create the new approximation
		approx.SetFrac(num2, den2)

		// Check if the approximation is within the tolerance
		diff := new(big.Rat).Sub(input, approx)
		if diff.Abs(diff).Cmp(tolerance) <= 0 {
			return approx
		}

		// Prepare for. the next iteration
		num0, num1 = num1, num2
		den0, den1 = den1, den2

		// Update x to the fractional part of its reciprocal
		x.Sub(x, new(big.Rat).SetInt(a))
		if x.Sign() == 0 {
			break
		}
		x.Inv(x)
	}

	// If no better approximation is found, return the input itself
	return input
}
