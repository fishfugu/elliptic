package ellipticcurve

import (
	"elliptic/pkg/bigarith"

	"fmt"
	"math"
	"math/big"
)

// NOTE: these are private / immutable on purpose
// All ECs defined in the Weierstrass form

// EllipticCurve represents an elliptic curve defined by the equation y^2 = x^3 + Ax + B.
type EllipticCurve struct {
	a, b *big.Int // Coefficients of the curve equation.
}

// FiniteFieldEC represents an elliptic curve over a finite field defined by the equation y^2 = x^3 + Ax + B.
type FiniteFieldEC struct {
	ec EllipticCurve
	// TODO: should this be definable by "strings" rather than big.Int?
	// Woould that make it easier to interface with bigarith functions?
	p *big.Int // Coefficients of the curve equation and prime modulus of the field.
}

// NewEllipticCurve creates a new EllipticCurve with given coefficients.
func NewEllipticCurve(a, b *big.Int) *EllipticCurve {
	return &EllipticCurve{a: a, b: b}
}

// NewFiniteFieldEC creates a new EllipticCurve, defined over a finite field, with given coefficients and modulus.
func NewFiniteFieldEC(a, b, p *big.Int) *FiniteFieldEC {
	EC := NewEllipticCurve(a, b)
	return &FiniteFieldEC{ec: *EC, p: p}
}

// GetDetails returns the coefficients A and B of the curve.
func (ec *EllipticCurve) GetDetails() (*big.Int, *big.Int) {
	return ec.a, ec.b
}

// GetDetails returns the coefficients A, B, and the modulus P of the finite field curve.
func (ffec *FiniteFieldEC) GetDetails() (*big.Int, *big.Int, *big.Int) {
	return ffec.ec.a, ffec.ec.b, ffec.p
}

// TODO: convert all of these into things that use bigarith instead of floats or big.Floats

// finds minimum value of X for an Elliptic Curve where y = 0
// for curve in Weierstrass form, this should be the
// lowest value of x for the whole curve in the real numbers

// solveCubic for form x^3 + Ax + B - real roots only
func (ec EllipticCurve) SolveCubic() ([]string, error) {
	var roots []string

	ABigInt, BBigInt := ec.GetDetails()
	A := ABigInt.String()
	B := BBigInt.String()

	// Calculate the discriminant
	AOver3, AOver3Err := bigarith.Divide(A, "3")
	BOver2, BOver2Err := bigarith.Divide(B, "2")
	AOver3Cubed, AOver3CubedErr := bigarith.Exp(AOver3, "3", "")
	BOver2Squared, BOver2SquaredErr := bigarith.Exp(BOver2, "2", "")
	delta, deltaErr := bigarith.Add(AOver3Cubed, BOver2Squared)
	if AOver3Err != nil || BOver2Err != nil || AOver3CubedErr != nil || BOver2SquaredErr != nil || deltaErr != nil {
		return nil, fmt.Errorf("error in some stage of creating delta in SolveCubic - AOver3Err: %v - BOver2Err: %v - AOver3CubedErr: %v - BOver2SquaredErr: %v - deltaErr: %v", AOver3Err, BOver2Err, AOver3CubedErr, BOver2SquaredErr, deltaErr)
	}

	deltaCmpToZero, err := bigarith.Cmp(delta, "0")
	if err != nil {
		return nil, fmt.Errorf("error creating deltaIsGreaterThanZero in SolveCubic")
	}

	if deltaCmpToZero > 0 {
		// One real root, two complex roots
		C, err := bigarith.Exp(delta, "0.5", delta)
		if err != nil {
			return nil, fmt.Errorf("error creating C in SolveCubic")
		}
		NegativeBOver2, err := bigarith.Multiply(BOver2, "-1")
		if err != nil {
			return nil, fmt.Errorf("error creating NegativeBOver2 in SolveCubic")
		}
		NegativeBOver2PlusC, err := bigarith.Add(NegativeBOver2, C)
		if err != nil {
			return nil, fmt.Errorf("error creating NegativeBOver2PlusC in SolveCubic")
		}
		u, err := bigarith.Exp(NegativeBOver2PlusC, "0.333333333333", NegativeBOver2PlusC)
		if err != nil {
			return nil, fmt.Errorf("error creating u in SolveCubic")
		}
		NegativeBOver2MinusC, err := bigarith.Subtract(NegativeBOver2, C)
		if err != nil {
			return nil, fmt.Errorf("error creating NegativeBOver2MinusC in SolveCubic")
		}
		v, err := bigarith.Exp(NegativeBOver2MinusC, "0.333333333333", NegativeBOver2PlusC)
		if err != nil {
			return nil, fmt.Errorf("error creating v in SolveCubic")
		}

		root, err := bigarith.Add(u, v)
		roots = append(roots, root)

	} else if deltaCmpToZero == 0 {
		// All roots are real, at least two are equal
		// u := math.Cbrt(-B / 2)
		// root1 := u + u
		// root2 := -u
		// roots = append(roots, root1, root2, root2)

	} else {
		// Three real roots (delta < 0)
		// theta := math.Acos(-B / 2 * math.Sqrt(27/math.Pow(A, 3)))
		// r := 2 * math.Sqrt(-A/3)
		// root1 := r * math.Cos(theta/3)
		// root2 := r * math.Cos((theta+2*math.Pi)/3)
		// root3 := r * math.Cos((theta+4*math.Pi)/3)
		// roots = append(roots, root1, root2, root3)
	}

	return roots, err
}

func (ffec FiniteFieldEC) SolveCubic() ([]string, error) {
	return ffec.ec.SolveCubic()
}

func (ec EllipticCurve) FindY(x float64) (float64, error) {
	aBigInt, bBigInt := ec.GetDetails()
	// TODO: work out how to use "BigInt Accuracy" values and do error handling here
	A, _ := aBigInt.Float64()
	B, _ := bBigInt.Float64()

	det := math.Pow(x, 3) + (A * x) + B
	if det < 0 {
		return 0, fmt.Errorf("'det' less than 0, cannot find real square root value (%.5f)", det)
	}

	return math.Sqrt(math.Pow(x, 3) + (A * x) + B), nil
}

func (ffec FiniteFieldEC) FindY(x float64) (float64, error) {
	return ffec.ec.FindY(x)
}
