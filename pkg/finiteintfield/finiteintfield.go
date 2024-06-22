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

	for x := "0"; ; {
		// Calculate rhs: x^3 + Ax + B
		x3, _ := bigarith.Exp(x, "3", p.String()) // x^3 mod p
		Ax, _ := bigarith.Multiply(A.String(), x)
		Ax, _ = bigarith.Mod(Ax, p.String()) // (A*x) mod p
		rhs, _ := bigarith.Add(x3, Ax)
		rhs, _ = bigarith.Add(rhs, B.String())
		rhs, _ = bigarith.Mod(rhs, p.String()) // (x^3 + Ax + B) mod p

		// Check for y^2 = rhs
		for y := "0"; ; {
			y2, _ := bigarith.Exp(y, "2", p.String()) // y^2 mod p
			if cmp, _ := bigarith.Cmp(y2, rhs); cmp == 0 {
				points = append(points, [2]string{x, y})
			}

			y, _ = bigarith.Add(y, "1")
			if cmp, _ := bigarith.Cmp(y, p.String()); cmp >= 0 {
				break
			}
		}

		x, _ = bigarith.Add(x, "1")
		if cmp, _ := bigarith.Cmp(x, p.String()); cmp >= 0 {
			break
		}
	}

	return points, nil
}

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
				if j%tickInterval == 0 && j != 0 {
					plane[i][j] = '-' // Ticks on the x-axis
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
	result := "\n2D Plane Visualization with Cartesian Axes, Reflection Line, and Scale:\n"
	for i, line := range plane {
		for j, char := range line {
			result += string(char) + " "
			if (i == p || i == p/2) && j == p { // Add scale numbers at the end of axes and reflection lines
				result += fmt.Sprintf(" %d", j)
			}
		}
		if i%tickInterval == 0 && i != p { // Label on y-axis
			result += fmt.Sprintf(" %d", p-i)
		}
		result += "\n"
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
