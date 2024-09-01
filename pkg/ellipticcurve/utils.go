package ellipticcurve

import (
	"fmt"
	"math/cmplx"
)

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

// TODO: turn solveCubic into something using bigarith instead of cmplx
// because it can't handle big numbers right now!
// half done here:
//
// func solveCubic(a, b, c, d string) ([][2]*big.Int, error) {
// 	cubicSolutions := [][2]*big.Int{}
// 	aCmp, err := bigarith.Cmp(a, "0")
// 	if err != nil {
// 		return cubicSolutions, err
// 	}
// 	if aCmp == 0 {
// 		return nil, fmt.Errorf("the coefficient 'a' must not be zero for a cubic equation")
// 	}

// 	// Convert to normalized cubic t^3 + pt + q = 0

// 	// Create p := (3*a*c - b*b) / (3 * a * a)
// 	acBigInt, err := bigarith.Multiply(a, c)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	threeACBigInt, err := bigarith.Multiply("3", acBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	bSquaredBigInt, err := bigarith.Multiply(b, b)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	aSquaredBigInt, err := bigarith.Multiply(a, a)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	threeASquaredBigInt, err := bigarith.Multiply("3", aSquaredBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	pNumerBigInt, err := bigarith.Subtract(threeACBigInt, bSquaredBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	pBigInt, err := bigarith.Divide(pNumerBigInt, threeASquaredBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	// Create q := (2*b*b*b - 9*a*b*c + 27*a*a*d) / (27 * a * a * a)
// 	bCubedBigInt, err := bigarith.Multiply(bSquaredBigInt, b)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	twoBCubedBigInt, err := bigarith.Multiply("2", bCubedBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	abBigInt, err := bigarith.Multiply(a, b)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	abcBigInt, err := bigarith.Multiply(abBigInt, c)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	nineABCBigInt, err := bigarith.Multiply("9", abcBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	aSquaredDBigInt, err := bigarith.Multiply(aSquaredBigInt, d)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	twentySevenASquaredDBigInt, err := bigarith.Multiply("27", aSquaredDBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	twoBCubedMinusnineABCBigInt, err := bigarith.Subtract(twoBCubedBigInt, nineABCBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	qNumerBigInt, err := bigarith.Add(twoBCubedMinusnineABCBigInt, twentySevenASquaredDBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	aCubedBigInt, err := bigarith.Multiply(aSquaredBigInt, a)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	twentySevenACubedBigInt, err := bigarith.Multiply("27", aCubedBigInt)
// 	if err != nil {
// 		return cubicSolutions, err
// 	}

// 	qBigInt, err := bigarith.Divide(qNumerBigInt, twentySevenACubedBigInt)

// 	// Calculate discriminant: D := cmplx.Rect(-(4*p*p*p + 27*q*q), 0)

// 	dBigInt, err :=

// 	cmplx.Rect(-0.5, 0)
// 	// Solution via Cardano's formula
// 	u := cmplx.Pow(cmplxNum(-0.5)*cmplxNum(q)+cmplx.Sqrt(D)/2, 1.0/3)
// 	v := cmplx.Pow(cmplxNum(-0.5)*cmplxNum(q)-cmplx.Sqrt(D)/2, 1.0/3)

// 	// Three roots
// 	x1 := u + v - complex(b/(3*a), 0)
// 	x2 := -0.5*(u+v) - complex(b/(3*a), 0) + cmplx.Sqrt(3)/2*(u-v)*complex(0, 1)
// 	x3 := -0.5*(u+v) - complex(b/(3*a), 0) - cmplx.Sqrt(3)/2*(u-v)*complex(0, 1)

// 	return []complex128{x1, x2, x3}, nil
// }

func cmplxNum(num float64) complex128 {
	return cmplx.Rect(num, 0)
}
