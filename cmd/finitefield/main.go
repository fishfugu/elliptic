package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/sirupsen/logrus"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"
)

const divider = "==========================="

func main() {
	logger := utils.InitialiseLogger("[FINFIELD/MAIN]")

	err := run(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

func run(logger *logrus.Logger) error {
	logger.Info("starting function run")

	// Define elliptic curves with parameters (a, b, p)
	curves := []struct {
		a, b, p *big.Int
	}{
		{big.NewInt(1), big.NewInt(-1), big.NewInt(17)},
		{big.NewInt(1), big.NewInt(1), big.NewInt(13)},
		{big.NewInt(2), big.NewInt(3), big.NewInt(43)},
	}

	// Process each curve
	for _, curve := range curves {
		// Prepare curve details
		curveInfo := fmt.Sprintf(
			"Elliptic Curve: y^2 = x^3 + %sx + %s\nDefined in a finite field modulo %s\n%s\n",
			curve.a.String(), curve.b.String(), curve.p.String(), divider,
		)
		fmt.Print(curveInfo)

		// Initialise elliptic curve object
		eCurve := ellipticcurve.NewFiniteFieldEC(curve.a, curve.b, curve.p)

		// Calculate points on the curve
		points, err := finiteintfield.CalculatePoints(*eCurve)
		if err != nil {
			utils.LogOnError(logger, err, fmt.Sprintf("error from CalculatePoints: %v\n", err), true)
			continue // Skip to the next curve on error
		}

		// Format and display points
		fmt.Println("Points on the Elliptic Curve:")
		finiteintfield.LogPoints(points)

		// Convert p to int for visualisation
		pInt64 := curve.p.Int64()
		if pInt64 > int64(^uint(0)>>1) { // Check for int overflow
			logger.Errorf("Field size too large for visualisation: p = %s", curve.p.String())
			continue
		}
		pInt := int(pInt64)

		// Visualise the curve points
		fmt.Println(finiteintfield.VisualisePoints(points, pInt))
	}

	return nil
}
