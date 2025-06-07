package ellipticcurve

import (
	"fmt"
	"math/big"

	"elliptic/pkg/utils"
)

// EllipticCurve represents an elliptic curve defined by the equation y^2 = x^3 + Ax + B
// Defined as big.Rat to allow for arbitrary precision, non-integer values for A and B
type EllipticCurve struct {
	a, b *big.Int
}

// FiniteFieldEC represents an elliptic curve over a finite field
// TODO: add in original A, B values (keep them)
// because they make a difference to some "real" calculations
// (and therefor the order of points, even though the points are the same)
type FiniteFieldEC struct {
	ec      *EllipticCurve
	a, b, p *big.Int
}

// Cubic is a cubic defined as y = x^3 + Ax^2 + Bx + C
type Cubic struct {
	a, b *big.Rat
}

// LPoint represents a Point on a Line
type LPoint struct {
	x *big.Rat
}

// ECPoint represents a Point on an Elliptic Curve
type ECPoint struct {
	x     *big.Rat
	isPos bool
}

// Line represents a line as: y = mx + b
type Line struct {
	m, b *big.Rat // y = mx + b
}

// NewEllipticCurve creates a new elliptic curve
func NewEllipticCurve(a, b *big.Int) *EllipticCurve {
	return &EllipticCurve{a: a, b: b}
}

// NewFiniteFieldEC creates a new finite field elliptic curve
// A and B are stored mod p - but the original A and B are stored as details of the EllipticCurve
func NewFiniteFieldEC(a, b, p *big.Int) *FiniteFieldEC {
	modA, modB := new(big.Int).Mod(a, p), new(big.Int).Mod(b, p)
	return &FiniteFieldEC{ec: NewEllipticCurve(a, b), a: modA, b: modB, p: p}
}

// NewCubic creates a new Cubic
func NewCubic(a, b *big.Rat) *Cubic {
	return &Cubic{a: a, b: b}
}

// LPoint...
func NewLPoint(x *big.Rat) LPoint {
	return LPoint{x: x}
}

// ECPoint...
func NewECPoint(x *big.Rat, isPos bool) ECPoint {
	return ECPoint{x: x, isPos: isPos}
}

// Line...
func NewLine(m, b *big.Rat) Line {
	return Line{m: m, b: b}
}

// NewCubic creates a new Cubic
// y = x^3 + Ax^2 + Bx + C
// converting it to y = x^3 + Ax + B form
func NewCubicWithXSquaredComponent(a, b, c *big.Rat) *Cubic {
	// Convert to depressed cubic form: t^3 + pt + q = 0
	p := new(big.Rat).Sub(b, new(big.Rat).Mul(new(big.Rat).SetFrac(utils.OneInt, big.NewInt(3)), a))
	q := new(big.Rat).Add(new(big.Rat).Mul(new(big.Rat).SetFrac(utils.TwoInt, big.NewInt(27)), new(big.Rat).Mul(a, a)), c)
	return &Cubic{a: p, b: q}
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
// NOTE: A and B have been converted mod p.
func (ffec *FiniteFieldEC) GetDetails() (*big.Int, *big.Int, *big.Int) {
	return ffec.a, ffec.b, ffec.p
}

// GetDetailsAsRats returns the coefficients A, B, and the modulus P of the finite field curve, as big.Rat values.
func (ffec *FiniteFieldEC) GetDetailsAsRats() (*big.Rat, *big.Rat, *big.Rat) {
	return new(big.Rat).SetInt(ffec.ec.a), new(big.Rat).SetInt(ffec.ec.b), new(big.Rat).SetInt(ffec.p)
}

// GetEC returns the Elliptic Curve object of the finite field curve.
func (ffec *FiniteFieldEC) GetEC() *EllipticCurve {
	return ffec.ec
}

func (l *Line) GetDetails() (*big.Rat, *big.Rat) {
	return l.m, l.b
}

// SolveCubic finds roots of the cubic equation defined in an Elliptic Curve
func (ec EllipticCurve) SolveCubic() ([]*big.Rat, error) {
	logger := utils.InitialiseLogger("[EllipticCurve/SolveCubic]")
	logger.Debug("starting function EllipticCurve/SolveCubic")

	A, B := new(big.Int).Set(ec.a), new(big.Int).Set(ec.b) // make sure these don't get edited while working with them
	discriminant := calcDiscriminant(A, B)
	roots := make([]*big.Rat, 0, 3)
	root1, err := newtonCubic(A, B)
	if err != nil {
		return nil, err
	}
	roots = append(roots, root1)

	if discriminant.Sign() == 0 {
		roots = handleDoubleRoot(new(big.Rat).SetInt(A), root1)
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
		minXWindow := new(big.Rat).Add(utils.ZeroRat, new(big.Rat).SetInt(xWindowShift))
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

// SolveCubic finds the real roots of the cubic equation defined in the Cubic object
func (c Cubic) SolveCubic() ([]*big.Rat, error) {
	logger := utils.InitialiseLogger("[EllipticCurve/SolveCubic]")
	logger.Debug("starting function EllipticCurve/SolveCubic")

	A, B := new(big.Rat).Set(c.a), new(big.Rat).Set(c.b) // make sure these don't get edited while working with them
	discriminant := calcDiscriminantRat(A, B)
	roots := make([]*big.Rat, 0, 3)
	root1, err := newtonCubicRat(A, B)
	if err != nil {
		return nil, err
	}
	roots = append(roots, root1)

	if discriminant.Sign() == 0 {
		roots = handleDoubleRoot(A, root1)
	} else if discriminant.Sign() < 0 {
		remainingRoots, err := findRemainingRoots(new(big.Rat).Set(A), root1)
		if err != nil {
			return nil, err
		}
		roots = append(roots, remainingRoots...)
	}

	return sortRoots(roots), nil
}

// FindY finds the y value - on a Line - y = mx + b
func (l Line) FindY(lPoint LPoint) (*big.Rat, error) {
	m, b := l.GetDetails()
	mx := new(big.Rat).Mul(m, lPoint.x)
	y := new(big.Rat).Add(mx, b)
	return y, nil
}

// FindY finds the y value - an EllipticCurve - x^3 + Ax + B
// it returns the positive y value - but the other value is simply the negative of that anyway
func (ec EllipticCurve) FindY(ecPoint ECPoint) (*big.Rat, error) {
	A, B := ec.GetDetailsAsRats()
	xRat := new(big.Rat).Set(ecPoint.x)
	Ax := new(big.Rat).Mul(A, xRat)
	xSquared := new(big.Rat).Mul(xRat, xRat)
	xCubed := new(big.Rat).Mul(xSquared, xRat)
	xCubedPlusAx := new(big.Rat).Add(xCubed, Ax)
	xCubedPlusAxPlusB := new(big.Rat).Add(xCubedPlusAx, B)
	sqrtXCubedPlusAxPlusB, err := utils.SqrtRat(xCubedPlusAxPlusB)
	if err != nil {
		return nil, err
	}
	if !ecPoint.isPos {
		sqrtXCubedPlusAxPlusB.Neg(sqrtXCubedPlusAxPlusB)
	}
	return sqrtXCubedPlusAxPlusB, nil
}

func (ffec FiniteFieldEC) FindY(ecPoint ECPoint) (*big.Rat, error) {
	result, err := ffec.ec.FindY(ecPoint)
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

// Utility functions
func calcDiscriminantRat(A, B *big.Rat) *big.Rat {
	aCubed := new(big.Rat).Mul(A, new(big.Rat).Mul(A, A))
	bSquared := new(big.Rat).Mul(B, B)

	return new(big.Rat).Add(new(big.Rat).Mul(aCubed, big.NewRat(4, 1)), new(big.Rat).Mul(bSquared, big.NewRat(27, 1)))
}

func handleDoubleRoot(A *big.Rat, root1 *big.Rat) []*big.Rat {
	gradient := new(big.Rat).Mul(utils.ThreeRat, new(big.Rat).Mul(root1, root1))
	gradient.Add(gradient, new(big.Rat).Set(A))
	if gradient.Sign() == 0 {
		root3 := new(big.Rat).Neg(new(big.Rat).Mul(root1, utils.TwoRat))
		return []*big.Rat{root1, root1, root3}
	}
	root2 := new(big.Rat).Neg(new(big.Rat).Quo(root1, utils.TwoRat))
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
	return new(big.Int).Quo(negB, utils.FourInt)
}

func quickEstimateRootRat(A, B *big.Rat) *big.Rat {
	// Check for B == 0 which means f(x) = x^3 + Ax and so x == 0 is a root
	if B.Sign() == 0 {
		return utils.ZeroRat
	}
	negB := new(big.Rat).Neg(B)

	// Check for division by zero (A = 0)
	// Reutrn -B as estimate
	if A.Sign() == 0 {
		return negB
	}
	// Compute -B / A
	return new(big.Rat).Quo(negB, utils.FourRat)
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
			new(big.Rat).Mul(utils.ThreeRat, new(big.Rat).Mul(x, x)), // 3x^2
			new(big.Rat).SetInt(A),
		)

		if fpx.Sign() != 0 { // Avoid division by zero
			delta.Quo(fx, fpx) // only do division if fpx not 0
		}

		// either way, subtract last delta again, from x, and keep going
		// and check for Rat simplification
		x = approximateRat(new(big.Rat).Sub(x, delta))

		if new(big.Rat).Abs(delta).Cmp(utils.ToleranceFractionRat) < 0 { // Convergence check
			break
		}
	}
	return new(big.Rat).Set(x), nil
}

// newtonCubicRat finds one root for the cubic in the form x^3 + Ax + B
func newtonCubicRat(A, B *big.Rat) (*big.Rat, error) {
	// Check for B == 0 which means f(x) = x^3 + Ax and so x == 0 is a root
	if B.Sign() == 0 {
		return new(big.Rat).SetInt64(0), nil
	}

	x := quickEstimateRootRat(A, B)
	delta := new(big.Rat).SetInt64(1) // just assume not 0 for now

	for {
		fx := new(big.Rat).Add(
			new(big.Rat).Add(
				new(big.Rat).Mul(x, new(big.Rat).Mul(x, x)), // x^3
				new(big.Rat).Mul(new(big.Rat).Set(A), x),
			),
			new(big.Rat).Set(B),
		)
		fpx := new(big.Rat).Add(
			new(big.Rat).Mul(utils.ThreeRat, new(big.Rat).Mul(x, x)), // 3x^2
			new(big.Rat).Set(A),
		)

		if fpx.Sign() != 0 { // Avoid division by zero
			delta.Quo(fx, fpx) // only do division if fpx not 0
		}

		// either way, subtract last delta again, from x, and keep going
		// and check for Rat simplification
		x = approximateRat(new(big.Rat).Sub(x, delta))

		if new(big.Rat).Abs(delta).Cmp(utils.ToleranceFractionRat) < 0 { // Convergence check
			break
		}
	}
	return x, nil
}

// solveQuadratic calculates the roots of a quadratic equation of the form
// x^2 + px + q = 0 and returns the two roots as big.Rat.
func solveQuadratic(p, q *big.Rat) ([]*big.Rat, error) {
	discriminant := new(big.Rat).Sub(new(big.Rat).Mul(p, p), new(big.Rat).Mul(utils.FourRat, q))
	if discriminant.Sign() < 0 {
		return nil, fmt.Errorf("no real roots; discriminant is negative")
	}

	sqrtDiscriminant, err := utils.SqrtRat(discriminant)
	if err != nil {
		return nil, fmt.Errorf("failed to compute square root: %v", err)
	}

	negP := new(big.Rat).Neg(p)
	root1 := new(big.Rat).Quo(new(big.Rat).Add(negP, sqrtDiscriminant), utils.TwoRat)
	root2 := new(big.Rat).Quo(new(big.Rat).Sub(negP, sqrtDiscriminant), utils.TwoRat)

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
	if new(big.Rat).Abs(new(big.Rat).Sub(input, nearestInt)).Cmp(utils.ToleranceFractionRat) <= 0 {
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
		if diff.Abs(diff).Cmp(utils.ToleranceFractionRat) <= 0 {
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

// FindY finds the y value - for an EllipticCurve - x^3 + Ax + B - in the Reals
// it returns the positive y value - unless finiteFieldY is negative, in which case it returns the negative value
// if adjusts all x values so they are minWindow <= x = min(y(0)) < maxWindow
func FindYOnReals(ec *FiniteFieldEC, xInit *big.Int, isPos bool) (*big.Rat, error) {
	logger := utils.InitialiseLogger("[FindYOnReals]")
	logger.Debug("starting function FindYOnReals")

	_, _, p := ec.GetDetails()
	pRat := new(big.Rat).SetInt(p)

	A, B := ec.GetEC().GetDetails()

	logger.Debug("getting left most root")
	minRoot := LeftmostRoot(A, B, p)

	xRat := new(big.Rat).SetInt(xInit)
	x := new(big.Int).Set(xInit)
	logger.Debug("testing for lower than window")
	for xRat.Cmp(minRoot) <= 0 {
		xRat.Add(xRat, pRat)
		x.Add(x, p)
	}
	maxWindow := new(big.Rat).Add(minRoot, new(big.Rat).SetInt(p))
	logger.Debug("testing for higher than window")
	for xRat.Cmp(maxWindow) > 0 {
		xRat.Sub(xRat, pRat)
		x.Sub(x, p)
	}

	logger.Debug("finding Y on Reals")
	Ax := new(big.Int).Mul(A, x)
	xSquared := new(big.Int).Mul(x, x)
	xCubed := new(big.Int).Mul(xSquared, x)
	xCubedPlusAx := new(big.Int).Add(xCubed, Ax)
	xCubedPlusAxPlusB := new(big.Int).Add(xCubedPlusAx, B)
	sqrtXCubedPlusAxPlusB, err := utils.SqrtRat(new(big.Rat).SetInt(xCubedPlusAxPlusB))
	if err != nil {
		return nil, err
	}
	if !isPos {
		sqrtXCubedPlusAxPlusB.Neg(sqrtXCubedPlusAxPlusB)
	}

	sqrtXCubedPlusAxPlusBFloat, _ := sqrtXCubedPlusAxPlusB.Float32()
	logger.Debugf("A: %s, B: %s, xInit: %s, xCubedPlusAxPlusB: %s, sqrtXCubedPlusAxPlusB: %f", A, B, xInit, xCubedPlusAxPlusB, sqrtXCubedPlusAxPlusBFloat)
	return approximateRat(sqrtXCubedPlusAxPlusB), nil
}

// LeftmostRoot finds the smallest real root of y = x^3 + Ax + B = 0
func LeftmostRoot(A, B, p *big.Int) *big.Rat {
	logger := utils.InitialiseLogger("[LeftmostRoot]")
	logger.Debug("starting function LeftmostRoot")

	// Create big.Rat representations of A and B
	aRat := new(big.Rat).SetInt(A)
	bRat := new(big.Rat).SetInt(B)

	// Initial guess for the root (heuristic based on dominant term)
	x := new(big.Rat).SetInt(B) // x_0 = B
	x.Neg(x)
	// TODO: ignore error?
	x, _ = utils.SqrtRat(x) // Rough heuristic for initial guess

	// Temporary variables
	delta := new(big.Rat)
	xSquared := new(big.Rat)
	xCubed := new(big.Rat)
	ax := new(big.Rat)
	fx := new(big.Rat)
	dfdx := new(big.Rat)

	// Newton-Raphson iteration
	logger.Debug("Newton-Raphson iteration")
	for i := 0; i < 100; i++ { // Limit iterations to avoid infinite loop
		logger.Debugf("i: %d", i)

		// Compute f(x) = x^3 + A*x + B
		xSquared.Mul(x, x)      // x^2
		xCubed.Mul(xSquared, x) // x^3
		ax.Mul(aRat, x)         // A*x
		fx.Add(xCubed, ax)      // x^3 + A*x
		fx.Add(fx, bRat)        // x^3 + A*x + B

		// Compute f'(x) = 3*x^2 + A
		dfdx.Mul(utils.ThreeRat, xSquared) // 3*x^2
		dfdx.Add(dfdx, aRat)               // 3*x^2 + A

		// Compute delta = f(x) / f'(x)
		delta.Quo(fx, dfdx)

		// Update x = x - delta
		x = approximateRat(x.Sub(x, delta))
		logger.Debugf("x: %s, delta: %s", x, delta)

		// Check for convergence |delta| < toleranceFractionRat
		if delta.Abs(delta).Cmp(utils.ToleranceFractionRat) < 0 {
			break
		}
	}

	return x
}

// WindowShiftReal moves a finite field from 0 <= (x, y) < p based on xOffset, yOffset
func WindowShiftReal(point [2]*big.Rat, xOffset, yOffset *big.Rat, p *big.Int) [2]*big.Rat {
	pRat := new(big.Rat).SetInt(p)

	minXWindow := new(big.Rat).Add(utils.ZeroRat, xOffset)
	maxXWindow := new(big.Rat).Add(pRat, xOffset)
	minYWindow := new(big.Rat).Add(utils.ZeroRat, yOffset)
	maxYWindow := new(big.Rat).Add(pRat, yOffset)

	shiftedX := new(big.Rat).Set(point[0]) // x value
	shiftedY := new(big.Rat).Set(point[1]) // y value
	if (xOffset != nil) && (xOffset.Sign() != 0) {
		for shiftedX.Cmp(minXWindow) < 0 { // if shiftedX too low, increase it
			shiftedX.Add(shiftedX, pRat)
		}
		for shiftedX.Cmp(maxXWindow) >= 0 { // if shiftedX too high, decrease it
			shiftedX.Sub(shiftedX, pRat)
		}
	}
	if (yOffset != nil) && (yOffset.Sign() != 0) {
		for shiftedY.Cmp(minYWindow) < 0 { // if shiftedY too low, increase it
			shiftedY.Add(shiftedY, pRat)
		}
		for shiftedY.Cmp(maxYWindow) >= 0 { // if shiftedY too high, decrease it
			shiftedY.Sub(shiftedY, pRat)
		}
	}
	return [2]*big.Rat{shiftedX, shiftedY}
}

// CalculateLine returns a Line that passes through p1 and p2
// or the line tanget to the Curve, at p1, if p1 == p2
// (x1, y1) (x2, y2)
// y = ((y2 - y1)/(x2 - x1)) x + b
// when 1st point is (0, y3)
func CalculateLine(ec *EllipticCurve, p1, p2 ECPoint) (*Line, error) {
	var slope *big.Rat

	// TODO: should this be stored in the objewct not have to be calculated each time
	// ... it could be stored as immutible
	y1, err := ec.FindY(p1)
	if err != nil {
		return nil, err
	}
	y2, err := ec.FindY(p2)
	if err != nil {
		return nil, err
	}
	if p1.x.Cmp(p2.x) == 0 && p1.isPos == p2.isPos {
		// Tangent line: y' = (3x^2 + A) / 2y
		aRat, _ := ec.GetDetailsAsRats()
		xSquared := new(big.Rat).Mul(p1.x, p1.x)
		threeXSquared := new(big.Rat).Mul(utils.ThreeRat, xSquared)
		numerator := new(big.Rat).Add(threeXSquared, aRat) // 3x^2 + A
		denominator := new(big.Rat).Mul(utils.TwoRat, y1)
		slope = new(big.Rat).Quo(numerator, denominator)
	} else {
		// Line between two points
		deltaX := new(big.Rat).Sub(p2.x, p1.x)   // x2 - x1
		deltaY := new(big.Rat).Sub(y2, y1)       // y2 - y1
		slope = new(big.Rat).Quo(deltaY, deltaX) // (y2 - y1) / (x2 - x1)
	}
	// b = y1 - m.x1
	b := new(big.Rat).Sub(
		y1,
		new(big.Rat).Mul(slope, p1.x),
	)

	return &Line{
		m: slope,
		b: b,
	}, nil
}

func SubtractLineFromEllipticCurve(ec *EllipticCurve, line *Line) (*Cubic, error) {
	// Get coefficients of the elliptic curve
	A := new(big.Rat).SetInt(ec.a)
	B := new(big.Rat).SetInt(ec.b)

	// Get coefficients of the line
	m := line.m
	b := line.b

	// Calculate the coefficients of the resulting cubic equation
	aPrime := new(big.Rat).Neg(new(big.Rat).Mul(m, m))                                    // -m^2
	bPrime := new(big.Rat).Sub(A, new(big.Rat).Mul(new(big.Rat).Mul(utils.TwoRat, m), b)) // A - 2mb
	cPrime := new(big.Rat).Sub(B, new(big.Rat).Mul(b, b))                                 // B - b^2

	// Create the cubic equation
	return NewCubicWithXSquaredComponent(aPrime, bPrime, cPrime), nil
}

// func Intersection(ec *EllipticCurve, line *Line) (*big.Int, error) {
// 	aRat := new(big.Rat).SetInt(ec.a)
// 	bRat := new(big.Rat).SetInt(ec.b)
// 	mRat := new(big.Rat).SetInt(line.m)
// 	bRatLine := new(big.Rat).SetInt(line.b)

// 	// Substituting line equation into the curve equation
// 	cubicA := new(big.Rat).SetInt(oneInt) // x^3 coefficient is 1
// 	cubicB := new(big.Rat).Mul(mRat, mRat)
// 	cubicC := new(big.Rat).Mul(twoRat, mRat)
// 	cubicD := new(big.Rat).Add(bRatLine, aRat)

// 	// Solve cubic for x
// 	roots, err := SolveCubic(cubicA, cubicB, cubicC, cubicD)
// 	if err != nil || len(roots) != 3 {
// 		return nil, fmt.Errorf("failed to solve cubic equation")
// 	}

// 	// Returning the third root (not p1 or p2)
// 	return roots[2].Num(), nil
// }

func AddPoints(ec *EllipticCurve, p1, p2 ECPoint) (*ECPoint, error) {
	line, err := CalculateLine(ec, p1, p2)
	if err != nil {
		return nil, err
	}

	newCubic, err := SubtractLineFromEllipticCurve(ec, line)
	if err != nil {
		return nil, err
	}

	rootsNewCubic, err := newCubic.SolveCubic()
	if err != nil {
		return nil, err
	}
	if len(rootsNewCubic) != 3 {
		return nil, fmt.Errorf("issue with number of roots to new cubic. found %d roots. roots: %+v. newCubic.a: %s, newCubic.b: %s", len(rootsNewCubic), rootsNewCubic, newCubic.a.FloatString(10), newCubic.b.FloatString(10))
	}
	countUnmatched := 0
	p3 := &ECPoint{}
	for _, root := range rootsNewCubic {
		if !((root == p1.x) || (root == p2.x)) {
			countUnmatched++
			y, err := line.FindY(LPoint{x: root})
			if err != nil {
				return nil, err
			}
			p3 = &ECPoint{x: root, isPos: y.Cmp(utils.ZeroRat) >= 0}
		}
	}
	if countUnmatched != 1 {
		return nil, fmt.Errorf("wrong number of unmatched roots. Found %d unmatched roots. Roots: %+v", countUnmatched, rootsNewCubic)
	}

	return p3, nil
}
