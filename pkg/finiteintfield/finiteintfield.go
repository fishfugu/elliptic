package finiteintfield

import (
	"elliptic/pkg/bigarith"
	"elliptic/pkg/ellipticcurve"

	"fmt"
	"math/big"
)

// CalculatePoints calculates all points on an elliptic curve y^2 = x^3 + Ax + B over a finite field defined by prime p
func CalculatePoints(FFEC ellipticcurve.FiniteFieldEC) ([][2]string, error) {
	var points [][2]string

	A, B, p := FFEC.GetDetails()

	// set up y^2 lookup
	ySquaredLookup := make(map[bigarith.Int][]bigarith.Int)
	for y := bigarith.NewInt("0"); y.Compare(p.Val()) < 0; y = y.Plus("1") {
		ySquared := y.ToThePowerOf("2", p.Val())
		ySquaredLookup[ySquared] = append(ySquaredLookup[ySquared], y)
	}

	for x := bigarith.NewInt("0"); x.Compare(p.Val()) < 0; x = x.Plus("1") {
		// Calculate rhs: x^3 + Ax + B
		xCubed := x.ToThePowerOf("3", p.Val()).Mod(p.Val())     // x^3 mod p
		Ax := A.Times(x.Val()).Mod(p.Val())                     // Ax mod p
		rhs := xCubed.Plus(Ax.Val()).Plus(B.Val()).Mod(p.Val()) // (x^3 + Ax + B) mod p

		// Check for y^2 = rhs
		if yList, ok := ySquaredLookup[rhs]; ok {
			// yList contains list of y values for which y^2 = rhs
			for _, y := range yList {
				points = append(points, [2]string{x.Val(), y.Val()})
			}
		}
	}

	return points, nil
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

// FormatPoints formats a list of points for easy command line reading
func FormatPoints(points [][2]string) string {
	result := "List of Points on the Curve:\n"
	for _, point := range points {
		result += fmt.Sprintf("Point (x, y): (%s, %s)\n", point[0], point[1])
	}
	return result
}

// VisualisePoints displays the points on a 2D text-based plane, including axes, reflection line, and scale indicators.
// TODO: turn this into something that accepts and uses a bigarith not an int
func VisualisePoints(points [][2]string, p int) string {
	plane := make([][]rune, p+1) // +1 to include the x-axis
	tickInterval := max(1, p/10) // Adjust the interval based on p, avoiding too many ticks

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
		x, _ := new(big.Int).SetString(point[0], 10)
		y, _ := new(big.Int).SetString(point[1], 10)
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
