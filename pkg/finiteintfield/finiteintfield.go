package finiteintfield

import (
	"elliptic/pkg/bigarith"
	"fmt"
	"math/big"
)

// TODO: make an elliptic curve in a finite field some kind of object / type / interface

// CalculatePoints calculates all points on an elliptic curve y^2 = x^3 + Ax + B over a finite field defined by prime p
func CalculatePoints(A, B, p string) ([][2]string, error) {
	var points [][2]string

	for x := "0"; ; {
		// Calculate rhs: x^3 + Ax + B
		x3, _ := bigarith.Exp(x, "3", p) // x^3 mod p
		Ax, _ := bigarith.Multiply(A, x)
		Ax, _ = bigarith.Mod(Ax, p) // (A*x) mod p
		rhs, _ := bigarith.Add(x3, Ax)
		rhs, _ = bigarith.Add(rhs, B)
		rhs, _ = bigarith.Mod(rhs, p) // (x^3 + Ax + B) mod p

		// Check for y^2 = rhs
		for y := "0"; ; {
			y2, _ := bigarith.Exp(y, "2", p) // y^2 mod p
			if cmp, _ := bigarith.Cmp(y2, rhs); cmp == 0 {
				points = append(points, [2]string{x, y})
			}

			y, _ = bigarith.Add(y, "1")
			if cmp, _ := bigarith.Cmp(y, p); cmp >= 0 {
				break
			}
		}

		x, _ = bigarith.Add(x, "1")
		if cmp, _ := bigarith.Cmp(x, p); cmp >= 0 {
			break
		}
	}

	return points, nil
}

// FormatPoints formats a list of points for easy command line reading
func FormatPoints(points [][2]string) string {
	result := "List of Points on the Curve:\n"
	for _, point := range points {
		result += fmt.Sprintf("(%s, %s)\n", point[0], point[1])
	}
	return result
}

// VisualizePoints displays the points on a 2D text-based plane
// TODO: turn this into something that accepts and uses a bigsrith not an int
func VisualizePoints(points [][2]string, p int) string {
	plane := make([][]rune, p)
	for i := range plane {
		plane[i] = make([]rune, p)
		for j := range plane[i] {
			plane[i][j] = '.'
		}
	}

	for _, point := range points {
		x, _ := new(big.Int).SetString(point[0], 10)
		y, _ := new(big.Int).SetString(point[1], 10)
		plane[y.Int64()][x.Int64()] = '*'
	}

	result := "\n2D Plane Visualization:\n"
	for _, line := range plane {
		for _, char := range line {
			result += string(char) + " "
		}
		result += "\n"
	}

	return result
}
