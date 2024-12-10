package main

import (
	ba "elliptic/pkg/bigarith"
	"elliptic/pkg/ellipticcurve"
	"flag"
	"fmt"
	"os"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error percolated to main, error: %v\n", err)
	}
}

func run() error {
	// Define command-line flags for A, B, and P
	aFlag := flag.String("A", "", "Coefficient A of the elliptic curve")
	bFlag := flag.String("B", "", "Coefficient B of the elliptic curve")
	pFlag := flag.String("P", "", "Prime modulus P of the finite field")
	// visualiseFlag := flag.Bool("visualise", false, "Visualise the points on the elliptic curve")
	flag.Parse()

	// Prompt user for missing values
	a := getInput("A", *aFlag)
	b := getInput("B", *bFlag)
	p := getInput("P", *pFlag)
	fmt.Println()

	// Convert inputs to bigarith.Int
	aBigarithInt := ba.NewInt(a)
	bBigarithInt := ba.NewInt(b)

	switch p {
	case "":
		// Create the elliptic curve over the finite field
		curve := ellipticcurve.NewEllipticCurve(aBigarithInt, bBigarithInt)

		// Calc / display "lowest x value where y = 0" value
		roots, err := curve.SolveCubic()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error solving cubic: %s\nCurve: %v", err, *curve)
			os.Exit(1)
		}

		fmt.Printf("Roots: %v\n", roots)
	default:
		pBigarithInt := ba.NewInt(p)

		// Create the elliptic curve over the finite field
		curve := ellipticcurve.NewFiniteFieldEC(aBigarithInt, bBigarithInt, pBigarithInt)

		// Calculate points on the curve
		// points, err := finiteintfield.CalculatePoints(*curve)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Error calculating points: %v\n", err)
		// 	os.Exit(1)
		// }

		// Output points
		// fmt.Println(finiteintfield.FormatPoints(points))

		// Visualise points if the flag is set
		// if *visualiseFlag {
		// 	// Convert bigarith.Int to int
		// 	// TODO: change function to handle string input for better compatibility
		// 	// TODO: handle this error
		// 	pInt, _ := strconv.Atoi(pBigarithInt.Val())
		// 	visualisation := finiteintfield.VisualisePoints(points, pInt)
		// 	fmt.Println(visualisation)
		// }

		// Calc / display "lowest x value where y = 0" value
		roots, err := curve.SolveCubic()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error solving cubic: %s\nCurve: %v", err, *curve)
			os.Exit(1)
		}

		fmt.Printf("Roots (mod %s): %v\n", p, roots)
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

	return nil
}

// getInput prompts the user for input if the provided value is empty.
func getInput(name, value string) string {
	if value == "" {
		fmt.Printf("Please enter the value for %s: ", name)
		fmt.Scanln(&value)
	}
	return value
}
