package ellipticcurve

import (
	"fmt"
	"math/big"

	"elliptic/pkg/utils"
)

// 2048 bit precision was chosen as approximately 4 times the MAXIMUM number of bits
// used for. keys in EC Cryptography in the highest security situations:
// https://en.wikipedia.org/wiki/Key_size?utm_source=chatgpt.com
// "521-bit keys: Deliver a security level of roughly 256 bits, used in scenarios requiring
// the highest security assurances."

var (
	oneInt, twoInt, fourInt                       = big.NewInt(1), big.NewInt(2), big.NewInt(4)
	zeroRat, twoRat, threeRat, fourRat            = big.NewRat(0, 1), big.NewRat(2, 1), big.NewRat(3, 1), big.NewRat(4, 1)
	precision_2048                     uint       = 2048
	tolerance_1024                     uint64     = 1024
	toleranceInt                                  = big.NewInt(int64(tolerance_1024))
	twoToThePowerOfTolerance                      = new(big.Int).Exp(twoInt, toleranceInt, nil)
	toleranceFractionRat               *big.Rat   = new(big.Rat).SetFrac(oneInt, twoToThePowerOfTolerance)
	halfFloat                          *big.Float = utils.NewFloat().SetFloat64(0.5)
)

// EllipticCurve represents an elliptic curve defined by the equation y^2 = x^3 + Ax + B
type EllipticCurve struct {
	a, b *big.Int
}

// FiniteFieldEC represents an elliptic curve over a finite field
type FiniteFieldEC struct {
	ec *EllipticCurve
	p  *big.Int
}

// NewEllipticCurve creates a new elliptic curve
func NewEllipticCurve(a, b *big.Int) *EllipticCurve {
	return &EllipticCurve{a: a, b: b}
}

// NewFiniteFieldEC creates a new finite field elliptic curve
func NewFiniteFieldEC(a, b, p *big.Int) *FiniteFieldEC {
	modA, modB := new(big.Int).Mod(a, p), new(big.Int).Mod(b, p)
	return &FiniteFieldEC{ec: NewEllipticCurve(modA, modB), p: p}
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

// GetEC returns the Elliptic Curve object of the finite field curve.
func (ffec *FiniteFieldEC) GetEC() *EllipticCurve {
	return ffec.ec
}

// GetDetailsAsRats returns the coefficients A, B, and the modulus P of the finite field curve, as big.Rat values.
func (ffec *FiniteFieldEC) GetDetailsAsRats() (*big.Rat, *big.Rat, *big.Rat) {
	return new(big.Rat).SetInt(ffec.ec.a), new(big.Rat).SetInt(ffec.ec.b), new(big.Rat).SetInt(ffec.p)
}

// SolveCubic finds roots of the cubic equation
func (ec EllipticCurve) SolveCubic() ([]*big.Rat, error) {
	logger := utils.InitialiseLogger("[EllipticCurve/SolveCubic]")
	logger.Debug("starting function EllipticCurve/SolveCubic")

	A, B := new(big.Int).Set(ec.a), new(big.Int).Set(ec.b)    // amke sure these don't get edited while working with them
	logger.Debugf("1 EllipticCurve A: %s, B: %s", A, B)       // for text output to screen
	logger.Debugf("2 EllipticCurve A: %s, B: %s", ec.a, ec.b) // for text output to screen
	discriminant := calcDiscriminant(A, B)

	logger.Debugf("3 EllipticCurve A: %s, B: %s", A, B)       // for text output to screen
	logger.Debugf("4 EllipticCurve A: %s, B: %s", ec.a, ec.b) // for text output to screen

	roots := make([]*big.Rat, 0, 3)
	root1, err := newtonCubic(A, B)
	if err != nil {
		return nil, err
	}
	roots = append(roots, root1)

	if discriminant.Sign() == 0 {
		roots = handleDoubleRoot(A, root1)
	} else if discriminant.Sign() < 0 {
		remainingRoots, err := findRemainingRoots(new(big.Rat).SetInt(A), root1)
		if err != nil {
			return nil, err
		}
		roots = append(roots, remainingRoots...)
	}

	return sortRoots(roots), nil
}

func (ffec FiniteFieldEC) SolveCubic(xWindowShift *big.Int) ([]*big.Rat, error) {
	logger := utils.InitialiseLogger("[FiniteFieldEC/SolveCubic]")
	logger.Debug("starting function FiniteFieldEC/SolveCubic")

	_, _, p := ffec.GetDetails()
	pRat := new(big.Rat).SetInt(p)

	logger.Debugf("1 FiniteFieldEC A: %s, B: %s, P: %s", ffec.ec.a, ffec.ec.b, p) // for text output to screen

	roots, err := ffec.ec.SolveCubic()
	logger.Debugf("2 FiniteFieldEC A: %s, B: %s, P: %s", ffec.ec.a, ffec.ec.b, p) // for text output to screen
	// convert each value into its mod p equivalent
	for i, result := range roots {
		roots[i] = modRatInt(result, p)
	}

	// if xWindowShift exists and is not 0
	// shift all the x-values for the points
	// by enough to put thwm in the right window
	if (xWindowShift != nil) && (xWindowShift.Sign() != 0) {
		minXWindow := new(big.Rat).Add(zeroRat, new(big.Rat).SetInt(xWindowShift))
		maxXWindow := new(big.Rat).Add(pRat, new(big.Rat).SetInt(xWindowShift))
		for _, root := range roots {
			for root.Cmp(minXWindow) < 0 {
				root.Add(root, pRat)
			}
			for root.Cmp(maxXWindow) >= 0 {
				root.Sub(root, pRat)
			}
		}
	}

	return roots, err
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

// Utility functions
func calcDiscriminant(A, B *big.Int) *big.Int {
	aCubed := new(big.Int).Mul(A, new(big.Int).Mul(A, A))
	bSquared := new(big.Int).Mul(B, B)

	return new(big.Int).Add(new(big.Int).Mul(aCubed, big.NewInt(4)), new(big.Int).Mul(bSquared, big.NewInt(27)))
}

func handleDoubleRoot(A *big.Int, root1 *big.Rat) []*big.Rat {
	gradient := new(big.Rat).Mul(threeRat, new(big.Rat).Mul(root1, root1))
	gradient.Add(gradient, new(big.Rat).SetInt(A))
	if gradient.Sign() == 0 {
		root3 := new(big.Rat).Neg(new(big.Rat).Mul(root1, twoRat))
		return []*big.Rat{root1, root1, root3}
	}
	root2 := new(big.Rat).Neg(new(big.Rat).Quo(root1, twoRat))
	return []*big.Rat{root1, root2, root2}
}

func sortRoots(roots []*big.Rat) []*big.Rat {
	// If there's only one root, no sorting is needed
	if len(roots) <= 1 {
		return roots
	}

	// Simple insertion sort for up to 3 elements
	for i := 1; i < len(roots); i++ {
		key := roots[i]
		j := i - 1

		// Move elements of roots[0..i-1] that are greater than key
		// to one position ahead of their current position
		for j >= 0 && roots[j].Cmp(key) > 0 {
			roots[j+1] = roots[j]
			j--
		}
		roots[j+1] = key
	}

	return roots
}

// modRatInt computes `a mod b` where `a` is a big.Rat and `b` is a big.Int.
// The result is a big.Rat such that 0 <= result < b (converted to a big.Rat).
func modRatInt(a *big.Rat, b *big.Int) *big.Rat {
	bRat := new(big.Rat).SetInt(b)        // Convert b to big.Rat
	quotient := new(big.Rat).Quo(a, bRat) // Compute quotient
	quotientFloor := new(big.Rat).SetInt(quotient.Num().Div(quotient.Num(), quotient.Denom()))
	remainder := new(big.Rat).Sub(a, new(big.Rat).Mul(quotientFloor, bRat))

	if remainder.Sign() < 0 { // Ensure result in [0, b)
		remainder.Add(remainder, bRat)
	}
	return remainder
}

// sqrtRat computes the square root of a big.Rat with arbitrary precision.
// If the result is not an exact rational number, it computes an approximation with the specified precision.
func sqrtRat(input *big.Rat) (*big.Rat, error) {
	num, den := input.Num(), input.Denom()

	// Check if numerator and denominator are perfect squares
	if sqrtNum, sqrtDen := utils.IntSqrt(num), utils.IntSqrt(den); sqrtNum != nil && sqrtDen != nil {
		return new(big.Rat).SetFrac(sqrtNum, sqrtDen), nil
	}

	// Approximation for non-perfect square roots
	floatInput := new(big.Float).SetPrec(precision_2048).SetRat(input)
	floatSqrt := sqrtFloat(floatInput, precision_2048)

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
	for i := uint(0); i < prec; i++ {
		temp := utils.NewFloat().Quo(a, guess)     // temp = a / guess
		temp2 := utils.NewFloat().Add(guess, temp) // temp2 = (guess + a/guess)
		guess = utils.NewFloat().Mul(temp2, halfFloat)
	}

	return guess
}

// QuickEstimateRoot estimates a root of the cubic equation x^3 + Ax + B = 0
// by solving the linear approximation Ax + B = 0 -> x = -B/A
// Using the linear approximation y = Ax + B
// becuase this line intersects curve at x=0 -> (0, B)
// and has the same gradient = A at that point: y' = 3x^2 + B
// except where B=0, when 0 is def a root - so return 0
// and except where A=0 so y=B never crosses y-axis
// so just guess x = -B
func quickEstimateRoot(A, B *big.Int) *big.Int {
	// Check for B == 0 which means f(x) = x^3 + Ax and so x == 0 is a root
	if B.Sign() == 0 {
		return new(big.Int).SetInt64(0)
	}
	negB := new(big.Int).Neg(B)

	// Check for division by zero (A = 0)
	// Reutrn -B as estimate
	if A.Sign() == 0 {
		return negB
	}
	// Compute -B / A
	return new(big.Int).Quo(negB, fourInt)
}

// newtonCubic finds one root for the cubic in the form x^3 + Ax + B
func newtonCubic(A, B *big.Int) (*big.Rat, error) {
	// Check for B == 0 which means f(x) = x^3 + Ax and so x == 0 is a root
	if B.Sign() == 0 {
		return new(big.Rat).SetInt64(0), nil
	}

	x := new(big.Rat).SetInt(quickEstimateRoot(A, B))
	delta := new(big.Rat).SetInt64(1) // just assume not 0 for now

	for {
		fx := new(big.Rat).Add(
			new(big.Rat).Add(
				new(big.Rat).Mul(x, new(big.Rat).Mul(x, x)), // x^3
				new(big.Rat).Mul(new(big.Rat).SetInt(A), x),
			),
			new(big.Rat).SetInt(B),
		)
		fpx := new(big.Rat).Add(
			new(big.Rat).Mul(threeRat, new(big.Rat).Mul(x, x)), // 3x^2
			new(big.Rat).SetInt(A),
		)

		if fpx.Sign() != 0 { // Avoid division by zero
			delta.Quo(fx, fpx) // only do division if fpx not 0
		}

		// either way, subtract last delta again, from x, and keep going
		// and check for Rat simplification
		x = approximateRat(new(big.Rat).Sub(x, delta))

		if new(big.Rat).Abs(delta).Cmp(toleranceFractionRat) < 0 { // Convergence check
			break
		}
	}
	return x, nil
}

// solveQuadratic calculates the roots of a quadratic equation of the form
// x^2 + px + q = 0 and returns the two roots as big.Rat.
func solveQuadratic(p, q *big.Rat) ([]*big.Rat, error) {
	discriminant := new(big.Rat).Sub(new(big.Rat).Mul(p, p), new(big.Rat).Mul(fourRat, q))
	if discriminant.Sign() < 0 {
		return nil, fmt.Errorf("no real roots; discriminant is negative")
	}

	sqrtDiscriminant, err := sqrtRat(discriminant)
	if err != nil {
		return nil, fmt.Errorf("failed to compute square root: %v", err)
	}

	negP := new(big.Rat).Neg(p)
	root1 := new(big.Rat).Quo(new(big.Rat).Add(negP, sqrtDiscriminant), twoRat)
	root2 := new(big.Rat).Quo(new(big.Rat).Sub(negP, sqrtDiscriminant), twoRat)

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
	nearestInt := new(big.Rat).SetInt(new(big.Int).Div(input.Num(), input.Denom()))
	if new(big.Rat).Abs(new(big.Rat).Sub(input, nearestInt)).Cmp(toleranceFractionRat) <= 0 {
		return nearestInt
	}

	return bestRationalApproximation(input) // Approximation if not near an integer
}

// bestRationalApproximation finds the best rational approximation of input within the given tolerance.
func bestRationalApproximation(input *big.Rat) *big.Rat {
	logger := utils.InitialiseLogger("[bestRationalApproximation]")
	logger.Debug("starting function bestRationalApproximation")

	// Continued fraction expansion to find the best rational approximation
	// Initialise variables
	a := new(big.Int)
	num0, num1, den0, den1 := big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(0)
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
		if diff.Abs(diff).Cmp(toleranceFractionRat) <= 0 {
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

// FindY finds the y value - an EllipticCurve - x^3 + Ax + B - in the Reals
// it returns the positive y value - but the other value is simply the negative of that anyway
func FindYOnReals(ec *EllipticCurve, x *big.Int, finiteFieldY *big.Int) (*big.Rat, error) {
	logger := utils.InitialiseLogger("[FindYOnReals]")
	logger.Debug("starting function FindYOnReals")

	A, B := ec.GetDetails()
	Ax := new(big.Int).Mul(A, x)
	xSquared := new(big.Int).Mul(x, x)
	xCubed := new(big.Int).Mul(xSquared, x)
	xCubedPlusAx := new(big.Int).Add(xCubed, Ax)
	xCubedPlusAxPlusB := new(big.Int).Add(xCubedPlusAx, B)
	sqrtXCubedPlusAxPlusB, err := sqrtRat(new(big.Rat).SetInt(xCubedPlusAxPlusB))
	if err != nil {
		return nil, err
	}
	if finiteFieldY.Sign() < 0 {
		sqrtXCubedPlusAxPlusB.Neg(sqrtXCubedPlusAxPlusB)
	}

	return sqrtXCubedPlusAxPlusB, nil
}
