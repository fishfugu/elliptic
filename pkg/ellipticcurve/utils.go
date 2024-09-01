package ellipticcurve

import (
	"fmt"
	"math/cmplx"
)

// finds minimum value of X for an Elliptic Curve where y = 0
// for curve in Weierstrass form, this should be the
// lowest value of x for the whole curve in the real numbers
func FindMinX(FFEC FiniteFieldEC) (float64, error) {
	aBigInt, bBigInt, _ := FFEC.GetDetails()
	A, _ := aBigInt.Float64()
	B, _ := bBigInt.Float64()
	roots, err := solveCubic(1, A, 0, B)
	if err != nil {
		return 0, err
	}
	minX := real(roots[0])
	for _, root := range roots {
		if real(root) < minX {
			minX = real(root)
		}
	}
	return minX, nil
}

// func FindY(FFEC FiniteFieldEC) (float64, error) {
// 	return
// }

func solveCubic(a, b, c, d float64) ([]complex128, error) {
	if a == 0 {
		return nil, fmt.Errorf("the coefficient 'a' must not be zero for a cubic equation")
	}

	// Convert to normalized cubic t^3 + pt + q = 0
	p := (3*a*c - b*b) / (3 * a * a)
	q := (2*b*b*b - 9*a*b*c + 27*a*a*d) / (27 * a * a * a)

	// Calculate discriminant
	D := cmplx.Rect(-(4*p*p*p + 27*q*q), 0)

	cmplx.Rect(-0.5, 0)
	// Solution via Cardano's formula
	u := cmplx.Pow(cmplxNum(-0.5)*cmplxNum(q)+cmplx.Sqrt(D)/2, 1.0/3)
	v := cmplx.Pow(cmplxNum(-0.5)*cmplxNum(q)-cmplx.Sqrt(D)/2, 1.0/3)

	// Three roots
	x1 := u + v - complex(b/(3*a), 0)
	x2 := -0.5*(u+v) - complex(b/(3*a), 0) + cmplx.Sqrt(3)/2*(u-v)*complex(0, 1)
	x3 := -0.5*(u+v) - complex(b/(3*a), 0) - cmplx.Sqrt(3)/2*(u-v)*complex(0, 1)

	return []complex128{x1, x2, x3}, nil
}

func cmplxNum(num float64) complex128 {
	return cmplx.Rect(num, 0)
}
