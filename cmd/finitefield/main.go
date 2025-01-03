package main

import (
	"fmt"
	"math/big"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"

	"github.com/sirupsen/logrus"
)

const divider = "==========================="

func main() {
	logger := utils.InitialiseLogger("[FINFIELD]")

	err := run(logger)
	utils.LogOnError(logger, err, fmt.Sprintf("error percolated to main, error: %v\n", err), false)
}

func run(logger *logrus.Logger) error {
	curves := []struct {
		a, b, p *big.Int
	}{
		{new(big.Int).SetInt64(1), new(big.Int).SetInt64(-1), new(big.Int).SetInt64(17)},
		{new(big.Int).SetInt64(1), new(big.Int).SetInt64(1), new(big.Int).SetInt64(13)},
		{new(big.Int).SetInt64(2), new(big.Int).SetInt64(3), new(big.Int).SetInt64(43)},
	}

	for _, curve := range curves {
		fmt.Printf("Elliptic Curve: y^2 = x^3 + %sx + %s\nDefined in a finite field modulo %s\n%s\n", curve.a.String(), curve.b.String(), curve.p.String(), divider)

		eCurve := ellipticcurve.NewFiniteFieldEC(curve.a, curve.b, curve.p)
 
		points, err := finiteintfield.CalculatePoints(*eCurve)
		utils.LogOnError(logger, err, fmt.Sprintf("error from CalculatePoints: %v\n", err), true)

		fmt.Println("Points on the Elliptic Curve:")
		stringPoints := finiteintfield.FormatPoints(points)
		fmt.Println(stringPoints)

		// Convert bigarith.Int to int
		// TODO: change function to handle string input for better compatibility
		// TODO: handle this error
		pInt := int(curve.p.Int64())
		visualisation := finiteintfield.VisualisePoints(points, pInt)
		fmt.Println(visualisation)
	}

	return nil
}
