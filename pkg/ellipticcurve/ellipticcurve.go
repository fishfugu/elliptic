package ellipticcurve

import (
	"elliptic/pkg/bigarith"
	"fmt"
	"os"

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
// where A and B are both integers represented bigarith.Int objects
type EllipticCurve struct {
	a, b bigarith.Int // Coefficients of the curve equation.
}

// FiniteFieldEC represents an elliptic curve over a finite field defined by the equation y^2 = x^3 + Ax + B.
type FiniteFieldEC struct {
	ec EllipticCurve
	// TODO: should this be definable by "strings" rather than big.Int?
	// Woould that make it easier to interface with bigarith functions?
	p bigarith.Int // Coefficients of the curve equation and prime modulus of the field.
}

// NewEllipticCurve creates a new EllipticCurve with given coefficients.
func NewEllipticCurve(a, b bigarith.Int) *EllipticCurve {
	return &EllipticCurve{a: a, b: b}
}

// NewFiniteFieldEC creates a new EllipticCurve, defined over a finite field, with given coefficients and modulus.
func NewFiniteFieldEC(a, b, p bigarith.Int) *FiniteFieldEC {
	EC := NewEllipticCurve(a.Mod(p.Val()), b.Mod(p.Val()))
	return &FiniteFieldEC{ec: *EC, p: p}
}

// GetDetails returns the coefficients A and B of the curve.
func (ec *EllipticCurve) GetDetails() (bigarith.Int, bigarith.Int) {
	return ec.a, ec.b
}

// GetDetails returns the coefficients A, B, and the modulus P of the finite field curve.
func (ffec *FiniteFieldEC) GetDetails() (bigarith.Int, bigarith.Int, bigarith.Int) {
	return ffec.ec.a, ffec.ec.b, ffec.p
}

// finds minimum value of X for an Elliptic Curve where y = 0
// for curve in Weierstrass form, this should be the
// lowest value of x for the whole curve in the real numbers

// solveCubic for form x^3 + Ax + B - real roots only
func (ec EllipticCurve) SolveCubic() ([]string, error) {
	var roots []string

	A, B := ec.GetDetails()

	// Calculate the discriminant - (A/3)^3 + (B/2)^2
	// = (A^3/27) + (B^2/4) = (4 A^3 / 108) + (27 B^2 / 108) = (4 A^3 + 27 B^2) / 108
	NegativeBOver2 := B.DividedBy("2").Neg()
	Aover3Cubed := A.ToThePowerOf("3", "").DividedBy("27")
	BOver2Squared := B.ToThePowerOf("2", "").DividedBy("4")

	discriminant := Aover3Cubed.Plus(BOver2Squared.Val())

	discriminantCmpToZero := discriminant.Compare("0")
	if discriminantCmpToZero > 0 {
		// One real root, two complex roots
		sqrtDiscriminant := discriminant.SquareRoot()
		NegativeBOver2PlusSqrtDiscriminant := NegativeBOver2.Plus(sqrtDiscriminant.Val())
		u := NegativeBOver2PlusSqrtDiscriminant.NthRoot("3")

		NegativeBOver2MinusSqrtDiscriminant := NegativeBOver2.Minus(sqrtDiscriminant.Val())
		v := NegativeBOver2MinusSqrtDiscriminant.NthRoot("3")

		root := u.Plus(v.Val()).Val()
		roots = append(roots, root)
	} else if discriminantCmpToZero == 0 {
		// All roots are real, at least two are equal
		u := NegativeBOver2.NthRoot("3")
		root1 := u.Times("2").Val()
		root2 := u.Neg().Val()
		roots = append(roots, root1, root2, root2)
	} else {
		// Three real roots (discriminant < 0)
		fmt.Printf("discriminant: %s", discriminant.Val())

		// r = \sqrt{\frac{-A^3}{27}}
		r := A.ToThePowerOf("3", "").DividedBy("27").Neg().SquareRoot()
		logrus.Debugf("r: %s", r.Val())

		// \cos(\theta) = -\frac{B}{2r}
		// \theta) = \arccos{ -\frac{B}{2r} }
		theta := B.DividedBy(r.Val()).DividedBy("2").Neg().ArcCos()
		logrus.Debugf("theta: %s", theta.Val())

		// 2. **Check for Valid Range**:
		// If the value computed for \(\cos(\theta)\) is outside the range, it indicates a numerical issue or a mistake in the conversion. For three real roots, the correct trigonometric approach should always yield a valid angle \(\theta\).

		// 3. **Calculate the Three Real Roots**:
		// Once you have the correct value for \(\theta\), compute the roots using:
		// \[
		// x_k = 2\sqrt{\frac{-A}{3}} \cos\left( \frac{\theta + 2k\pi}{3} \right) \quad \text{for } k = 0, 1, 2
		// \]

		// ### Numerical Stability and Adjustments

		// One possible issue is numerical instability or rounding errors that push the value of \(\cos(\theta)\) outside the acceptable range. To handle this situation:

		// - **Clamping**: Ensure that the computed value of \(\cos(\theta)\) is clamped to lie within \([-1, 1]\):
		// \[
		// \cos(\theta) = \max(-1, \min(1, \cos(\theta)))
		// \]
		// - This adjustment ensures that you always get a valid angle for \(\arccos(\cos(\theta))\), which will lead to the correct roots.

		// ### Explanation of Why It Happens
		// The reason this occurs is that the formula for \(\cos(\theta)\) can be sensitive to numerical precision, especially when dealing with values close to the limits of floating-point representation. The cubic equation itself may have perfectly valid real roots, but the trigonometric calculation must be treated carefully to avoid errors.
		// Using the corrected approach to calculate \(\cos(\theta)\) should lead to the three real roots \(-6\), \(2\), and \(4\) as expected for the values \(A = -28\) and \(B = 48\).

		// Calculate 2 Pi / 3 and 4 Pi over 3
		twoPi := bigarith.NewRational("2").TimesPi()
		fourPi := bigarith.NewRational("4").TimesPi()

		//	\[
		//		\cos(\theta) = -\frac{3B}{2A} \sqrt{\frac{27}{-A}} // - ( ( 3 * 48 ) / ( 2 * ( -28 ) ) ) sqrt( 27 / - (-28) ) =
		//	\]
		//	theta = arccos ( -\frac{3B}{2A} \sqrt{\frac{27}{-A}} )

		// all the final root values are calculated by multiplying by: 2 * sqrt{ - A / 3 }

		sqrtMinusAOver3Times2 := bigarith.NewInt(A.Val()).DividedBy("3").Neg().SquareRoot().Times("2")
		// (theta plus 2kpi) / 3

		cosThetaOver3 := theta.DividedBy("3").Cos()
		cosThetaPlus2PiOver3 := theta.Plus(twoPi.Val())
		cosThetaPlus4PiOver3 := theta.Plus(fourPi.Val())

		// The three real roots \(x_1\), \(x_2\), and \(x_3\) can then be computed as:
		// \[
		// x_k = 2\sqrt{\frac{-A}{3}} \cos\left( \frac{\theta + 2k\pi}{3} \right)
		// \]
		// where \(k = 0, 1, 2\) for the three distinct roots.
		// 	return roots, nil
		// }

		root1 := sqrtMinusAOver3Times2.Times(cosThetaOver3.Val()).Val()
		root2 := sqrtMinusAOver3Times2.Times(cosThetaPlus2PiOver3.Val()).Val()
		root3 := sqrtMinusAOver3Times2.Times(cosThetaPlus4PiOver3.Val()).Val()

		roots = append(roots, root1, root2, root3)
		logrus.Debug("Roots calculated for all real and all distinct roots: ", roots)
		logrus.Debugf("Rootz: %s", roots)
	}

	return roots, nil
}

func (ffec FiniteFieldEC) SolveCubic() ([]string, error) {
	results, err := ffec.ec.SolveCubic()
	// convert each value into its mod p equivalent
	for i, result := range results {
		results[i] = bigarith.NewRational(result).Mod(ffec.p.Val()).Val()
	}
	return results, err
}

// FindY finds the y value for an EllipticCurve - x^3 + Ax + B
// it returns the positive y value - but the other value is simply the negative of that anyway
func (ec EllipticCurve) FindY(x bigarith.Rational) (bigarith.Rational, error) {
	A, B := ec.GetDetails()
	Ax := A.Times(x.Val())
	return x.ToThePowerOf("3").Plus(Ax.Val()).Plus(B.Val()).SquareRoot(), nil
}

func (ffec FiniteFieldEC) FindY(x bigarith.Rational) (bigarith.Rational, error) {
	result, err := ffec.ec.FindY(x)
	return result.Mod(ffec.p.Val()), err
}
