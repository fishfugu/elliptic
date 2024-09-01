package ellipticcurve

import (
	"fmt"
	"math"
	"math/big"
	"math/cmplx"
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
func (ec EllipticCurve) SolveCubic() ([]float64, error) {
	// TODO: do error handling of floats
	// TODO: convert this to strings - or big int...?
	A, _ := ec.a.Float64()
	B, _ := ec.b.Float64()
	det := math.Abs(-math.Pow((A/3.0), 3.0) - math.Pow((B/2.0), 2.0))

	fmt.Printf("det: %v\n", det)

	// Calculate discriminant
	D := cmplx.Rect(-B/2.0, math.Sqrt(det))
	cubeRootD := cmplx.Pow(D, 1.0/3.0)
	R := real(cubeRootD)
	I := imag(cubeRootD)

	// Three roots
	x1 := 2.0 * R
	x2 := -R + (math.Sqrt(3.0) * I)
	x3 := -R - (math.Sqrt(3.0) * I)

	return []float64{x1, x2, x3}, nil
}

func (ffec FiniteFieldEC) SolveCubic() ([]float64, error) {
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
