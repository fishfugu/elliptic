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
	visualizeFlag := flag.Bool("visualise", false, "Visualise the points on the elliptic curve")
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
	if *visualizeFlag {
		visualization := finiteintfield.VisualisePoints(points, int(pInt.Int64()))
		fmt.Println(visualization)
	}
}

// getInput prompts the user for input if the provided value is empty.
func getInput(name, value string) string {
	if value == "" {
		fmt.Printf("Please enter the value for %s: ", name)
		fmt.Scanln(&value)
	}
	return value
}
