package main

import (
	"fmt"
	"strconv"

	ba "elliptic/pkg/bigarith"
	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
)

const divider = "==========================="

func main() {
	curves := []struct {
		a, b, p ba.Int
	}{
		{ba.NewInt("1"), ba.NewInt("-1"), ba.NewInt("17")},
		{ba.NewInt("1"), ba.NewInt("1"), ba.NewInt("13")},
		{ba.NewInt("2"), ba.NewInt("3"), ba.NewInt("43")},
	}

	for _, curve := range curves {
		fmt.Printf("Elliptic Curve: y^2 = x^3 + %sx + %s\nDefined in a finite field modulo %s\n%s\n", curve.a.Val(), curve.b.Val(), curve.p.Val(), divider)

		eCurve := ellipticcurve.NewFiniteFieldEC(curve.a, curve.b, curve.p)

		points, err := finiteintfield.CalculatePoints(*eCurve)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Points on the Elliptic Curve:")
		stringPoints := finiteintfield.FormatPoints(points)
		fmt.Println(stringPoints)

		// Convert bigarith.Int to int
		// TODO: change function to handle string input for better compatibility
		// TODO: handle this error
		pInt, _ := strconv.Atoi(curve.p.Val())
		visualisation := finiteintfield.VisualisePoints(points, pInt)
		fmt.Println(visualisation)
	}
}
