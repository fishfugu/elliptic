package main

import (
	"fmt"
	"math/big"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
)

const divider = "==========================="

func main() {
	curves := []struct {
		a, b, p *big.Int
	}{
		{big.NewInt(1), big.NewInt(-1), big.NewInt(17)},
		{big.NewInt(1), big.NewInt(1), big.NewInt(13)},
		{big.NewInt(2), big.NewInt(3), big.NewInt(43)},
	}

	for _, curve := range curves {
		fmt.Printf("Elliptic Curve: y^2 = x^3 + %sx + %s\nDefined in a finite field modulo %s\n%s\n", curve.a.String(), curve.b.String(), curve.p.String(), divider)

		eCurve := ellipticcurve.NewFiniteFieldEC(curve.a, curve.b, curve.p)

		points, err := finiteintfield.CalculatePoints(*eCurve)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Points on the Elliptic Curve:")
		stringPoints := finiteintfield.FormatPoints(points)
		fmt.Println(stringPoints)

		// Convert big.Int to int
		// TODO: change function to handle string input for better compatibility
		visualisation := finiteintfield.VisualisePoints(points, int(curve.p.Int64()))
		fmt.Println(visualisation)
	}
}
