package main

import (
	"fmt"
	"math/big"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
)

func main() {
	a, b, p := big.NewInt(1), big.NewInt(-1), big.NewInt(17)
	curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

	points, err := finiteintfield.CalculatePoints(*curve)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Points on the Elliptic Curve:")
	for _, point := range points {
		fmt.Printf("Point (x, y): (%s, %s)\n", point[0], point[1])
	}
}
