package finiteintfield

import (
	"fmt"
	"math/big"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/utils"
)

var (
	zeroInt, oneInt, twoInt, threeInt *big.Int
)

func init() {
	zeroInt, oneInt, twoInt, threeInt = big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3)
}

// CalculatePoints calculates all points on an elliptic curve y^2 = x^3 + Ax + B over a finite field defined by prime p
// while points on a finite field are usually reported:
// (x, y): 0 <= x < p, 0 <= y < p
// these points are returned in a window such that:
// xWindowShift <= x < p + xWindowShift
// yWindowShift <= y < p + yWindowShift
func CalculatePoints(FFEC *ellipticcurve.FiniteFieldEC, xWindowShift, yWindowShift *big.Int) (points [][2]*big.Int, realPoints [][2]*big.Rat, err error) {
	logger := utils.InitialiseLogger("[CalculatePoints]")
	logger.Debug("starting function CalculatePoints")

	A, B, p := FFEC.GetDetails()
	logger.Debugf("1 CalculatePoints A: %s, B: %s, p: %s", A, B, p)

	// set up y^2 lookup
	ySquaredLookup := make(map[string][]*big.Int)
	ySquaredModP := big.NewInt(1)
	for y := big.NewInt(0); y.Cmp(p) < 0; y.Add(y, oneInt) {
		ySquaredModP.Exp(y, twoInt, p) // y squared mod p
		ySquaredLookup[ySquaredModP.String()] = append(ySquaredLookup[ySquaredModP.String()], new(big.Int).Set(y))
		logger.Debugf("ySquaredLookup: %v, ySquaredModP: %s, y: %v", ySquaredLookup, ySquaredModP.String(), y)
	}

	xCubed, Ax, xCubedPlusAx, xCubedPlusAxPlusB, rhs := new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	logger.Debugf("p: %v, zeroInt.Cmp(p): %v, zeroInt: %v", p, zeroInt.Cmp(p), zeroInt)

	minXWindow := new(big.Int).Add(zeroInt, xWindowShift)
	maxXWindow := new(big.Int).Add(p, xWindowShift)
	minYWindow := new(big.Int).Add(zeroInt, yWindowShift)
	maxYWindow := new(big.Int).Add(p, yWindowShift)
	for x := big.NewInt(0); x.Cmp(p) < 0; x.Add(x, oneInt) {
		// Calculate rhs: x^3 + Ax + B
		xCubed.Exp(x, threeInt, p)             // x^3 mod p
		Ax.Mul(A, x)                           // A times x
		xCubedPlusAx.Add(xCubed, Ax)           // x^3 + Ax
		xCubedPlusAxPlusB.Add(xCubedPlusAx, B) // x^3 + Ax + B
		rhs.Mod(xCubedPlusAxPlusB, p)          // (x^3 + Ax + B) mod p

		logger.Debugf("ySquaredLookup: %v, rhs: %v, ySquaredLookup[new(big.Int).Set(rhs)]: %v", ySquaredLookup, rhs, ySquaredLookup[rhs.String()])

		// Check for y^2 = rhs
		if yList, ok := ySquaredLookup[rhs.String()]; ok {
			// yList contains list of y values for which y^2 = rhs exists
			for _, y := range yList {
				shiftedX := new(big.Int).Set(x)
				// if xWindowShift exists and is not 0
				// shift all the x-values for the points
				// by enough to put them in the right window
				if (xWindowShift != nil) && (xWindowShift.Sign() != 0) {
					for shiftedX.Cmp(minXWindow) < 0 {
						shiftedX.Add(shiftedX, p)
					}
					for shiftedX.Cmp(maxXWindow) >= 0 {
						shiftedX.Sub(shiftedX, p)
					}
				}

				shiftedY := new(big.Int).Set(y)
				// if yWindowShift exists and is not 0
				// shift all the y-values for the points
				// by enough to put them in the right window
				if (yWindowShift != nil) && (yWindowShift.Sign() != 0) {
					for shiftedY.Cmp(minYWindow) < 0 {
						shiftedY.Add(shiftedY, p)
					}
					for shiftedY.Cmp(maxYWindow) >= 0 {
						shiftedY.Sub(shiftedY, p)
					}
				}

				points = append(points, [2]*big.Int{shiftedX, shiftedY})

				// Use original 0 <= x < p
				// But use yWindowShift <= shiftedY < p + yWindowShift
				realY, err := ellipticcurve.FindYOnReals(FFEC.GetEC(), new(big.Int).Set(x), new(big.Int).Set(shiftedY))
				if err != nil {
					return nil, nil, err
				}
				realPoints = append(realPoints, [2]*big.Rat{new(big.Rat).SetInt(x), realY})
			}
		}
	}

	return points, realPoints, nil
}

// DivideInField takes two big ints, divides a by b,
// and returns the result as a string in a finite field defined by a prime modulus p.
// Returns an error if the input strings are not valid integers, or if b has no inverse modulo p,
// (which includes if b is 0).
func DivideInField(a, b, p *big.Int) (*big.Int, error) {
	// First, find the multiplicative inverse of b mod p
	bInv, err := ModularInverse(b, p)
	if err != nil {
		return nil, fmt.Errorf("error finding inverse: %v", err)
	}

	// Calculate a * bInv mod p
	result := new(big.Int).Mul(a, bInv) // a * b^-1
	return result.Mod(result, p), nil
}

// ModularInverse calculates the multiplicative inverse of a modulo p using Fermat's Little Theorem
// and returns it as a string. Returns an error if p is not prime or a and p are not coprime.
func ModularInverse(a, p *big.Int) (*big.Int, error) {
	logger := utils.InitialiseLogger("[ModularInverse]")
	logger.Debug("starting function ModularInverse")

	// Check if p is prime; if not, the multiplicative inverse might not exist
	// https://pkg.go.dev/math/big#Int.ProbablyPrime
	// From ---^ that site: "The probability of returning true for a randomly chosen non-prime is at most ¼ⁿ"
	// i.e. (1/4)^n - where n is the parameter handed in to the function
	// TODO: work out if there's a way to decide what that param should be set to
	if !p.ProbablyPrime(1000) {
		return nil, fmt.Errorf("invalid input: modulus %s is not prime", p)
	}

	// Calculate a^(p-2) mod p
	pMinusTwo := new(big.Int).Sub(p, big.NewInt(2)) // p-2
	return new(big.Int).Exp(a, pMinusTwo, p), nil   // Calculate a^(p-2) mod p
}

// FormatPoints prints out in the log, a list of points for easy reading
func LogPoints(points [][2]*big.Int) {
	logger := utils.InitialiseLogger("[ModularInverse]")
	logger.Debug("starting function ModularInverse")

	logger.Debug("list of points on the curve:")
	for _, point := range points {
		logger.Debugf("Point (x, y): (%s, %s)\n", point[0].String(), point[1].String())
	}
	logger.Debug("END OF list of points on the curve ===============")
}

// VisualisePoints displays the points on a 2D text-based plane, including axes, reflection line, and scale indicators.
// TODO: turn this into something that accepts and uses a big.math not an int
func VisualisePoints(points [][2]*big.Int, p int) string {
	plane := make([][]rune, p+1)       // +1 to include the x-axis
	tickInterval := utils.Max(1, p/10) // Adjust the interval based on p, avoiding too many ticks

	for i := range plane {
		plane[i] = make([]rune, p+1) // +1 to include the y-axis
		for j := range plane[i] {
			if i == p {
				plane[i][j] = '-' // x-axis at the bottom
				// TODO: get this to put scale under tick points, along x-axis, on extra line
				if j%tickInterval == 0 && j != 0 {
					plane[i][j] = '/' // Ticks on the x-axis
				}
			} else if j == 0 {
				plane[i][j] = '|' // y-axis on the left
				if i%tickInterval == 0 && i != p {
					plane[i][j] = '|' // Ticks on the y-axis
				}
			} else if i == p/2 {
				plane[i][j] = '.' // Reflection line at y = p/2
			} else {
				plane[i][j] = ' ' // Use space for empty
			}
		}
	}
	plane[p][0] = '+'   // Origin at the bottom-left corner
	plane[p/2][0] = '|' // Mark where the y = p/2 line intersects with the y-axis

	// Plot the points
	for _, point := range points {
		x := new(big.Int).Set(point[0])
		y := new(big.Int).Set(point[1])
		xInt, yInt := int(x.Int64()), int(y.Int64())
		plane[p-yInt][xInt] = '*'
	}

	// Construct the visual output
	result := "\n2D Plane Visualisation with Cartesian Axes, Reflection Line, and Scale:\n"
	for i, line := range plane {
		for j, char := range line {
			result += string(char) + " "
			// Add scale numbers at the end of x-axis
			if (i == p) && j == p {
				result += " 0"
			}

			// Add scale numbers at the end of reflection lines
			if (i == p/2) && j == p {
				result += fmt.Sprintf(" %d/2", p)
			}
		}
		// Label on end of line for y-axis scale
		if i%tickInterval == 0 && i != p {
			result += fmt.Sprintf(" %d", p-i)
		}
		result += "\n"
	}

	return result
}

func Div2RoundUp(p *big.Int) *big.Int {
	two := big.NewInt(2)
	remainder := new(big.Int).Mod(p, two)
	result := new(big.Int).Set(p)

	if remainder.Sign() != 0 { // p is odd
		result.Add(result, big.NewInt(1)) // Increment to make it even
	}

	return result.Div(result, two) // Divide by 2
}

// // isPerfectSquare checks if a number is a perfect square.
// func isPerfectSquare(n *big.Int) bool {
// 	if n.Sign() < 0 {
// 		return false
// 	}
// 	sqrt := new(big.Int).Sqrt(n)
// 	square := new(big.Int).Mul(sqrt, sqrt)
// 	return square.Cmp(n) == 0
// }

// func FindRealEquivalentPoint(A, B, p *big.Int, point [2]*big.Int) {
// 	// Define the modular constraints for the point we are testing
// 	mod := new(big.Int).Set(p)
// 	xMod := new(big.Int).Set(point[0])
// 	yMod := new(big.Int).Set(point[1])

// 	// Create temporary big.Int instances for calculations
// 	temp := new(big.Int)
// 	x := new(big.Int).Set(xMod)
// 	for k := big.NewInt(0); k.Cmp(big.NewInt(1000000)) < 0; k.Add(k, big.NewInt(1)) {
// 		// if new(big.Int).Mod(k, big.NewInt(100000)).Cmp(big.NewInt(0)) == 0 {
// 		// 	fmt.Printf("k: %s\n", k)
// 		// }

// 		// Compute x based on modular constraint: x = xMod + mod * k
// 		// just add mod each loop
// 		x := new(big.Int).Add(x, mod)

// 		// Compute y^2 = x^3 + A*x + B
// 		xCubed := new(big.Int).Exp(x, big.NewInt(3), nil)
// 		Ax := new(big.Int).Mul(A, x)
// 		ySquared := new(big.Int).Add(new(big.Int).Add(xCubed, Ax), B)

// 		// Check if ySquared is a perfect square
// 		if isPerfectSquare(ySquared) {
// 			// Compute y
// 			y := new(big.Int).Sqrt(ySquared)

// 			// Check y modular constraint: y % mod == yMod
// 			if temp.Mod(y, mod).Cmp(yMod) == 0 {
// 				fmt.Printf("Smallest integer solution: x = %s, y = %s\n", x.String(), y.String())
// 				return
// 			}
// 		}
// 	}
// }
