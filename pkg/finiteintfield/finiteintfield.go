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
func CalculatePoints(FFEC ellipticcurve.FiniteFieldEC, xWindowShift, yWindowShift *big.Int) ([][2]*big.Int, error) {
	logger := utils.InitialiseLogger("[CalculatePoints]")
	logger.Debug("starting function CalculatePoints")

	var points [][2]*big.Int

	A, B, p := FFEC.GetDetails()

	// set up y^2 lookup
	ySquaredLookup := make(map[string][]*big.Int)
	ySquaredModP := big.NewInt(1)
	for y := big.NewInt(0); y.Cmp(p) < 0; y.Add(y, oneInt) {
		ySquaredModP.Exp(y, twoInt, p) // y squared mod p
		ySquaredLookup[ySquaredModP.String()] = append(ySquaredLookup[ySquaredModP.String()], new(big.Int).Set(y))
		logger.Debugf("XXX -- ySquaredLookup: %v, ySquaredModP: %s, y: %v", ySquaredLookup, ySquaredModP.String(), y)
	}

	xCubed, Ax, xCubedPlusAx, xCubedPlusAxPlusB, rhs := new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	logger.Debugf("p: %v, zeroInt.Cmp(p): %v, zeroInt: %v", p, zeroInt.Cmp(p), zeroInt)

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
				points = append(points, [2]*big.Int{new(big.Int).Set(x), new(big.Int).Set(y)})
			}
		}
	}

	// if xWindowShift exists and is not 0
	// shift all the x-values for the points
	// by enough to put thwm in the right window
	if (xWindowShift != nil) && (xWindowShift.Sign() != 0) {
		minXWindow := new(big.Int).Add(zeroInt, xWindowShift)
		maxXWindow := new(big.Int).Add(p, xWindowShift)
		for _, point := range points {
			for point[0].Cmp(minXWindow) < 0 {
				point[0].Add(point[0], p)
			}
			for point[0].Cmp(maxXWindow) >= 0 {
				point[0].Sub(point[0], p)
			}
		}
	}

	// if xWindowShift exists and is not 0
	// shift all the x-values for the points
	// by enough to put thwm in the right window
	if (yWindowShift != nil) && (yWindowShift.Sign() != 0) {
		minYWindow := new(big.Int).Add(zeroInt, yWindowShift)
		maxYWindow := new(big.Int).Add(p, yWindowShift)
		for _, point := range points {
			for point[1].Cmp(minYWindow) < 0 {
				point[1].Add(point[1], p)
			}
			for point[1].Cmp(maxYWindow) >= 0 {
				point[1].Sub(point[1], p)
			}
		}
	}

	return points, nil
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

// THIS IS A HALF FINSIHED THOUGHT
// I was going to build a function that moved all the points into the range "first position" decribed below in comments
// But I don't think I need it
// I think I need the other way around - i.e. search for points in Z on EC  from y=0 position up the x-axis in integers...
// so... commenting out for now...
// TODO - finish if necessary - or delete with comment linking to half finished version in git history...
// ...
// // Take a list of points and an elliptic curve and convert them to version where
// // finite field box goes from { y: 1/2 p <= y < 1/2 p }, and calculate x = x_0 where y = 0,
// // { x: x_0 <= x < x_0 + p }
// func ConvertPointsToFirstPosition(FFEC ellipticcurve.FiniteFieldEC, initialPoints [][2]string) ([][2]string, error) {
// 	// Work out minimum value for which y = 0 on EC
// 	// x^3 + Ax + B = 0
// 	// x^3 + Ax = -B
// 	// x (x^2 + A) = -B
// 	_, _, pBigInt := FFEC.GetDetails()
// 	p := pBigInt.String()

// 	firstPositionPoints := [][2]string{}

// 	for _, point := range initialPoints {
// 		y := point[1]
// 		halfP, err := bigarith.Divide(p, "2")
// 		if err != nil {
// 			return firstPositionPoints, err
// 		}
// 		negHalfP := fmt.Sprintf("-%s", halfP)

// 		yHalfPCmp, err := bigarith.Cmp(y, halfP)
// 		if err != nil {
// 			return firstPositionPoints, err
// 		}
// 		// if y >= halfP ...
// 		for yHalfPCmp >= 0 {
// 			// keep subtracting p until it's in the right range
// 			newY, err := bigarith.Subtract(y, p)
// 			if err != nil {
// 				return firstPositionPoints, err
// 			}
// 			yHalfPCmp, err = bigarith.Cmp(newY, halfP)
// 			if err != nil {
// 				return firstPositionPoints, err
// 			}
// 			y = newY
// 		}

// 		yNegHalfPCmp, err := bigarith.Cmp(y, negHalfP)
// 		if err != nil {
// 			return firstPositionPoints, err
// 		}
// 		// if y < -halfP ...
// 		for yNegHalfPCmp < 0 {
// 			// keep adding p until it's in the right range
// 			newY, err := bigarith.Add(y, p)
// 			if err != nil {
// 				return firstPositionPoints, err
// 			}
// 			yNegHalfPCmp, err = bigarith.Cmp(newY, negHalfP)
// 			if err != nil {
// 				return firstPositionPoints, err
// 			}
// 			y = newY
// 		}
// 	}
// 	return firstPositionPoints, nil
// }

// FormatPoints prints out in the log, a list of points for easy reading
func LogPoints(points [][2]*big.Int) {
	logger := utils.InitialiseLogger("[ModularInverse]")
	logger.Debug("starting function ModularInverse")

	logger.Debug("list of points on the curve:")
	for _, point := range points {
		logger.Warnf("Point (x, y): (%s, %s)\n", point[0].String(), point[1].String())
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
