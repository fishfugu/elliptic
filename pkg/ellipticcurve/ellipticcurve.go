package ellipticcurve

import (
	"elliptic/pkg/bigarith"
)

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

// TODO: convert all of these into things that use bigarith instead of floats or big.Floats

// finds minimum value of X for an Elliptic Curve where y = 0
// for curve in Weierstrass form, this should be the
// lowest value of x for the whole curve in the real numbers

// solveCubic for form x^3 + Ax + B - real roots only
func (ec EllipticCurve) SolveCubic() ([]string, error) {
	var roots []string

	A, B := ec.GetDetails()

	// Calculate the discriminant - (A/3)^3 + (B/2)^2 = (A^3)/(3^3) + (B^2)/(2^2) = (A^3)/27 + (B^2)/4
	NegativeBOver2 := B.DividedBy("2").Neg()
	Aover3Cubed := A.ToThePowerOf("3", "").DividedBy("27")
	BOver2Squared := B.ToThePowerOf("2", "").DividedBy("4")

	discriminant := Aover3Cubed.Plus(BOver2Squared.Val())
	discriminantCmpToZero := discriminant.Compare("0")
	if discriminantCmpToZero > 0 {
		// One real root, two complex roots
		u := NegativeBOver2.Plus(discriminant.SquareRoot().Val()).NthRoot("3")
		v := NegativeBOver2.Minus(discriminant.SquareRoot().Val()).NthRoot("3")
		root := u.Plus(v.Val()).Val()
		roots = append(roots, root)
	} else if discriminantCmpToZero == 0 {
		// All roots are real, at least two are equal
		// u := math.Cbrt(-B / 2)
		u := NegativeBOver2.NthRoot("3")
		// root1 := u + u
		root1 := u.Times("2").Val()
		// root2 := -u
		root2 := u.Neg().Val()
		roots = append(roots, root1, root2, root2)
	} else {
		// Three real roots (discriminant < 0)
		// Step 1: Calculate r = 2 * sqrt(-A / 3)
		r := A.DividedBy("3").Neg().SquareRoot().Times("2")
		// Step 2: Calculate theta = arccos((-B / 2) * sqrt(27 / A^3))
		srt27OverACubed := bigarith.NewFloat("27").DividedBy(A.ToThePowerOf("3", "").Val()).SquareRoot()
		theta := NegativeBOver2.Times(srt27OverACubed.Val()).ArcCos()
		// Step 3: Calculate the three real roots
		Pi := bigarith.Pi()
		// root1 = r * cos(theta / 3)
		root1 := r.Times(theta.DividedBy("3").Cos().Val()).Val()
		// root2 = r * cos((theta + 2 * Pi) / 3)
		twoPiPlusTheta := Pi.Times("2").Plus(theta.Val())
		root2 := r.Times(twoPiPlusTheta.DividedBy("3").Cos().Val()).Val()
		// root3 = r * cos((theta + 4 * Pi) / 3)
		fourPiPlusTheta := Pi.Times("4").Plus(theta.Val())
		root3 := r.Times(fourPiPlusTheta.DividedBy("3").Cos().Val()).Val()
		// Append the three real roots
		roots = append(roots, root1, root2, root3)
	}

	return roots, nil
}

func (ffec FiniteFieldEC) SolveCubic() ([]string, error) {
	results, err := ffec.ec.SolveCubic()
	// convert each value into its mod p equivalent
	for i, result := range results {
		results[i] = bigarith.NewFloat(result).Mod(ffec.p.Val()).Val()
	}
	return results, err
}

// FindY finds the y value for an EllipticCurve - x^3 + Ax + B
// it returns the positive y value - but the other value is simply the negative of that anyway
func (ec EllipticCurve) FindY(x bigarith.Float) (bigarith.Float, error) {
	A, B := ec.GetDetails()
	Ax := A.Times(x.Val())
	return x.ToThePowerOf("3").Plus(Ax.Val()).Plus(B.Val()).SquareRoot(), nil
}

func (ffec FiniteFieldEC) FindY(x bigarith.Float) (bigarith.Float, error) {
	result, err := ffec.ec.FindY(x)
	return result.Mod(ffec.p.Val()), err
}
