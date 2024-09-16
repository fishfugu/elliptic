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

	// Calculate the discriminant - (A/3)^3 + (B/2)^2 = (A^3)/(3^3) + (B^2)/(2^2) = (A^3)/27 + (B^2)/4
	ACubed, ACubedErr := bigarith.Exp(A, "3", "")
	fmt.Printf("ACubed: %s\n", ACubed)
	BSquared, BSquaredErr := bigarith.Exp(B, "2", "")
	ACubedOver27, ACubedOver27Err := bigarith.Divide(ACubed, "27")
	BSquaredOver4, BSquaredOver4Err := bigarith.Divide(BSquared, "4")
	delta, deltaErr := bigarith.AddFloat(ACubedOver27, BSquaredOver4)
	if ACubedErr != nil || BSquaredErr != nil || ACubedOver27Err != nil || BSquaredOver4Err != nil || deltaErr != nil {
		return nil,
			fmt.Errorf(`error in some stage of creating delta in SolveCubic
				ACubedErr: %v
				BSquaredErr: %v
				ACubedOver9Err: %v
				BSquaredOver4Err: %v
				deltaErr: %v
`,
				ACubedErr,
				BSquaredErr,
				ACubedOver27Err,
				BSquaredOver4Err,
				deltaErr,
			)
	}
	// fmt.Printf("delta: %s\n", delta)

	deltaCmpToZero, err := bigarith.CmpFloat(delta, "0")
	if err != nil {
		return nil, fmt.Errorf("error creating deltaIsGreaterThanZero in SolveCubic - %v", err)
	}

	NegativeBOver2, err := bigarith.Divide(B, "-2")
	if err != nil {
		return nil, fmt.Errorf("error creating NegativeBOver2 in SolveCubic - %v", err)
	}

	if deltaCmpToZero > 0 {
		// One real root, two complex roots
		C, err := bigarith.Sqrt(delta, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating C in SolveCubic - %v", err)
		}
		NegativeBOver2PlusC, err := bigarith.AddFloat(NegativeBOver2, C)
		if err != nil {
			return nil, fmt.Errorf("error creating NegativeBOver2PlusC in SolveCubic - %v", err)
		}
		u, err := bigarith.CubeRoot(NegativeBOver2PlusC, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating u in SolveCubic - %v", err)
		}
		NegativeBOver2MinusC, err := bigarith.SubtractFloat(NegativeBOver2, C)
		if err != nil {
			return nil, fmt.Errorf("error creating NegativeBOver2MinusC in SolveCubic - %v", err)
		}
		v, err := bigarith.CubeRoot(NegativeBOver2MinusC, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating v in SolveCubic - %v", err)
		}

		root, err := bigarith.AddFloat(u, v)
		if err != nil {
			return nil, fmt.Errorf("error adding u and v in SolveCubic - %v", err)
		}
		roots = append(roots, root)

	} else if deltaCmpToZero == 0 {
		// All roots are real, at least two are equal
		// u := math.Cbrt(-B / 2)
		u, err := bigarith.CubeRoot(NegativeBOver2, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating u in SolveCubic - %v", err)
		}
		// root1 := u + u
		root1, err := bigarith.AddFloat(u, u)
		if err != nil {
			return nil, fmt.Errorf("error creating u + u in SolveCubic - %v", err)
		}
		// root2 := -u
		root2, err := bigarith.Neg(u)
		if err != nil {
			return nil, fmt.Errorf("error creating -u in SolveCubic - %v", err)
		}
		roots = append(roots, root1, root2, root2)

	} else {
		// Three real roots (delta < 0)

		// C = \sqrt{\frac{4A^3}{27}}
		FourACubed, err := bigarith.Exp(ACubed, "4", "")
		if err != nil {
			return nil, fmt.Errorf("error creating FourACubed in SolveCubic - %v", err)
		}
		FourACubedOver27, err := bigarith.Divide(FourACubed, "27")
		if err != nil {
			return nil, fmt.Errorf("error creating FourACubedOver27 in SolveCubic - %v", err)
		}
		C, err := bigarith.Sqrt(FourACubedOver27, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating SqrtFourACubedOver27 in SolveCubic - %v", err)
		}

		// \theta = \cos^{-1}\left(\frac{-B}{2C}\right)
		MinusTwoC, err := bigarith.Multiply("-2", C)
		if err != nil {
			return nil, fmt.Errorf("error creating MinusTwoC in SolveCubic - %v", err)
		}
		MinusBOverTwoC, err := bigarith.Divide(B, MinusTwoC)
		if err != nil {
			return nil, fmt.Errorf("error creating MinusBOverTwoC in SolveCubic - %v", err)
		}
		theta, err := bigarith.Acos(MinusBOverTwoC, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating theta in SolveCubic - %v", err)
		}

		// theta / 3
		thetaOver3, err := bigarith.Divide(theta, "3")
		if err != nil {
			return nil, fmt.Errorf("error creating thetaOver3 in SolveCubic - %v", err)
		}

		// gamma = \sqrt{\frac{-A}{3}}
		MinusAOver3, err := bigarith.Multiply(A, "-3")
		if err != nil {
			return nil, fmt.Errorf("error creating MinusAOver3 in SolveCubic - %v", err)
		}
		gamma, err := bigarith.Sqrt(MinusAOver3, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating gamma in SolveCubic - %v", err)
		}

		// x_1 = 2 gamma \cos\left(\frac{\theta}{3}\right)
		CosThetaOver3, err := bigarith.Cos(thetaOver3, 2048)
		if err != nil {
			return nil, fmt.Errorf("error creating CosThetaOver3 in SolveCubic - %v", err)
		}
		x_1, err := bigarith.Multiply(gamma, CosThetaOver3)
		if err != nil {
			return nil, fmt.Errorf("error creating x_1 in SolveCubic - %v", err)
		}

		roots = append(roots, x_1)

		// // x_2 = 2 gamma \cos\left(\frac{\theta}{3} + \frac{2\pi}{3}\right)
		// TwoPi, err := bigarith.MultiplyFloat("2", bigarith.Pi(2048).String())
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating TwoPi in SolveCubic - %v", err)
		// }
		// TwoPiOver3, err := bigarith.Divide(TwoPi, "3")
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating TwoPiOver3 in SolveCubic - %v", err)
		// }
		// thetaOver3PlusTwoPiOver3, err := bigarith.AddFloat(thetaOver3, TwoPiOver3)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating thetaOver3PlusTwoPiOver3 in SolveCubic - %v", err)
		// }
		// CosThetaOver3PlusTwoPiOver3, err := bigarith.Cos(thetaOver3PlusTwoPiOver3, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating thetaOver3PlusTwoPiOver3 in SolveCubic - %v", err)
		// }
		// // x_3 = 2 gamma \cos\left(\frac{\theta}{3} + \frac{4\pi}{3}\right)

		// // theta := math.Acos(-B / 2 * math.Sqrt(27/math.Pow(A, 3)))
		// Twenty7OverACubed, err := bigarith.Divide("27", ACubed)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating Twenty7OverACubed in SolveCubic - %v", err)
		// }
		// fmt.Printf("Twenty7OverACubed: %s\n", Twenty7OverACubed)
		// SqrtTwenty7OverACubed, err := bigarith.Sqrt(Twenty7OverACubed, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating SqrtTwenty7OverACubed in SolveCubic - %v", err)
		// }
		// fmt.Printf("SqrtTwenty7OverACubed: %s\n", SqrtTwenty7OverACubed)
		// NegativeBOver2MultipliedBySqrtTwenty7OverACubed, err := bigarith.MultiplyFloat(SqrtTwenty7OverACubed, NegativeBOver2)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating NegativeBOver2MultipliedBySqrtTwenty7OverACubed in SolveCubic - %v", err)
		// }
		// fmt.Printf("NegativeBOver2MultipliedBySqrtTwenty7OverACubed: %s\n", NegativeBOver2MultipliedBySqrtTwenty7OverACubed)
		// theta, err := bigarith.Acos(NegativeBOver2MultipliedBySqrtTwenty7OverACubed, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating theta in SolveCubic - %v", err)
		// }
		// fmt.Printf("theta: %s\n", theta)
		// // r := 2 * math.Sqrt(-A/3)
		// MinusAOver3, err := bigarith.Divide(A, "-3")
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating MinusAOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("MinusAOver3: %s\n", MinusAOver3)
		// SqrtMinusAOver3, err := bigarith.Sqrt(MinusAOver3, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating SqrtMinusAOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("SqrtMinusAOver3: %s\n", SqrtMinusAOver3)
		// r, err := bigarith.MultiplyFloat("2", SqrtMinusAOver3)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating r in SolveCubic - %v", err)
		// }
		// fmt.Printf("r: %s\n", r)
		// // root1 := r * math.Cos(theta/3)
		// ThetaOver3, err := bigarith.Divide(theta, "3")
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating ThetaOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("ThetaOver3: %s\n", ThetaOver3)
		// CosThetaOver3, err := bigarith.Cos(ThetaOver3, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating CosThetaOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("CosThetaOver3: %s\n", CosThetaOver3)
		// root1, err := bigarith.MultiplyFloat(r, CosThetaOver3)
		// // root2 := r * math.Cos((theta+2*math.Pi)/3)
		// TwoPi, err := bigarith.MultiplyFloat("2", bigarith.Pi(2048).String())
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating TwoPi in SolveCubic - %v", err)
		// }
		// fmt.Printf("TwoPi: %s\n", TwoPi)
		// ThetaPlusTwoPi, err := bigarith.AddFloat(theta, TwoPi)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating ThetaPlusTwoPi in SolveCubic - %v", err)
		// }
		// fmt.Printf("ThetaPlusTwoPi: %s\n", ThetaPlusTwoPi)
		// ThetaPlusTwoPiOver3, err := bigarith.Divide(ThetaPlusTwoPi, "3")
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating ThetaPlusTwoPiOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("ThetaPlusTwoPiOver3: %s\n", ThetaPlusTwoPiOver3)
		// CosThetaPlusTwoPiOver3, err := bigarith.Cos(ThetaPlusTwoPiOver3, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating CosThetaPlusTwoPiOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("CosThetaPlusTwoPiOver3: %s\n", CosThetaPlusTwoPiOver3)
		// root2, err := bigarith.MultiplyFloat(r, CosThetaPlusTwoPiOver3)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating root2 in SolveCubic - %v", err)
		// }
		// fmt.Printf("root2: %s\n", root2)
		// // root3 := r * math.Cos((theta+4*math.Pi)/3)
		// FourPi, err := bigarith.MultiplyFloat("4", bigarith.Pi(2048).String())
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating FourPi in SolveCubic - %v", err)
		// }
		// fmt.Printf("FourPi: %s\n", FourPi)
		// ThetaPlusFourPi, err := bigarith.AddFloat(theta, FourPi)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating ThetaPlusFourPi in SolveCubic - %v", err)
		// }
		// fmt.Printf("ThetaPlusFourPi: %s\n", ThetaPlusFourPi)
		// ThetaPlusFourPiOver3, err := bigarith.Divide(ThetaPlusFourPi, "3")
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating ThetaPlusFourPiOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("ThetaPlusFourPiOver3: %s\n", ThetaPlusFourPiOver3)
		// CosThetaPlusFourPiOver3, err := bigarith.Cos(ThetaPlusFourPiOver3, 2048)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating CosThetaPlusFourPiOver3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("CosThetaPlusFourPiOver3: %s\n", CosThetaPlusFourPiOver3)
		// root3, err := bigarith.MultiplyFloat(r, CosThetaPlusFourPiOver3)
		// if err != nil {
		// 	return nil, fmt.Errorf("error creating root3 in SolveCubic - %v", err)
		// }
		// fmt.Printf("root3: %s\n", root3)

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
