package main

import (
	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"flag"
	"fmt"
	"math/big"
	"os"
)

func main() {
	// Define command-line flags for A, B, and P
	aFlag := flag.String("A", "", "Coefficient A of the elliptic curve")
	bFlag := flag.String("B", "", "Coefficient B of the elliptic curve")
	pFlag := flag.String("P", "", "Prime modulus P of the finite field")
	visualiseFlag := flag.Bool("visualise", false, "Visualise the points on the elliptic curve")
	flag.Parse()

	// Prompt user for missing values
	a := getInput("A", *aFlag)
	b := getInput("B", *bFlag)
	p := getInput("P", *pFlag)

	// Convert inputs to big.Int
	aInt, _ := new(big.Int).SetString(a, 10)
	bInt, _ := new(big.Int).SetString(b, 10)
	pInt, _ := new(big.Int).SetString(p, 10)

	// Create the elliptic curve over the finite field
	curve := ellipticcurve.NewFiniteFieldEC(aInt, bInt, pInt)

	// Calculate points on the curve
	points, err := finiteintfield.CalculatePoints(*curve)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating points: %v\n", err)
		os.Exit(1)
	}

	// Output points
	fmt.Println(finiteintfield.FormatPoints(points))

	// Visualize points if the flag is set
	if *visualiseFlag {
		visualisation := finiteintfield.VisualisePoints(points, int(pInt.Int64()))
		fmt.Println(visualisation)
	}

	// Calc / display "lowest x value where y = 0" value
	roots, err := curve.SolveCubic()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error solving cubic: %s\nCurve: %v", err, *curve)
		os.Exit(1)
	}

	fmt.Printf("Roots: %v\n", roots)
	// minX, err := curve.SolveCubic()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error finding minimum x value where y = 0: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Minimum x value, where y = 0: %.5f\n", minX)

	// Check y vaqlue for x calculated
	// TODO: work out how to use "accuracy" and do error handling here
	// minXFloat, _ := minX.Float64()
	// y, err := curve.FindY(minXFloat)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error calculating y - checking minX value: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Check y value (should be '0'): %.5f\n", y)

	// Calc / display where tangent to EC goes through minX

}

// getInput prompts the user for input if the provided value is empty.
func getInput(name, value string) string {
	if value == "" {
		fmt.Printf("Please enter the value for %s: ", name)
		fmt.Scanln(&value)
	}
	return value
}
