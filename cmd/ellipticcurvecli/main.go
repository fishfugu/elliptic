package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/sirupsen/logrus"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"
)

func main() {
	logger := utils.InitialiseLogger("[ECCLI/MAIN]")

	err := run(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

func run(logger *logrus.Logger) error {
	logger.Info("starting function run")

	// Define command-line flags for. A, B, and P
	aFlag := flag.String("A", "", "Coefficient A of the elliptic curve")
	bFlag := flag.String("B", "", "Coefficient B of the elliptic curve")
	pFlag := flag.String("P", "", "Prime modulus P of the finite field")
	// visualiseFlag := flag.Bool("visualise", false, "Visualise the points on the elliptic curve")
	flag.Parse()

	// Prompt user for. input EC values
	a := getInput("A", *aFlag)
	b := getInput("B", *bFlag)
	p := getInput("P", *pFlag)
	fmt.Println()

	// Convert inputs to bigarith.Int
	// TODO: deal with errors here
	aBigInt, _ := new(big.Int).SetString(a, 10)
	bBigInt, _ := new(big.Int).SetString(b, 10)

	var roots []*big.Rat
	var err error

	switch p {
	case "":
		// Create the elliptic curve over the finite field
		curve := ellipticcurve.NewEllipticCurve(aBigInt, bBigInt)

		// Calc / display "lowest x value where y = 0" value
		roots, err = curve.SolveCubic()
		if err != nil {
			return err
		}
	default:
		// TODO: deal with error here
		pBigInt, ok := new(big.Int).SetString(p, 10)
		if !ok {
			return fmt.Errorf("error creating big.Int out of prime modulus value")
		}

		// Create the elliptic curve over the finite field
		curve := ellipticcurve.NewFiniteFieldEC(aBigInt, bBigInt, pBigInt)

		// Calculate points on the curve
		points, err := finiteintfield.CalculatePoints(*curve)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calculating points: %v\n", err)
			os.Exit(1)
		}

		// Output points
		finiteintfield.LogPoints(points)

		// Visualise points if the flag is set
		// if *visualiseFlag {
		// 	// Convert bigarith.Int to int
		// 	// TODO: change function to handle string input for. better compatibility
		// 	// TODO: handle this error
		// 	pInt, _ := strconv.Atoi(pBigarithInt.Val())
		// 	visualisation := finiteintfield.VisualisePoints(points, pInt)
		// 	fmt.Println(visualisation)
		// }

		// Calc / display "lowest x value where y = 0" value
		roots, err = curve.SolveCubic()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error solving cubic: %s\nCurve: %v", err, *curve)
			os.Exit(1)
		}

		// fmt.Printf("Roots (mod %s): %v\n", p, roots)
		// fmt.Println()
		// for i, root := range roots { // max 3 roots
		// 	rootFloat := utils.NewFloat().SetRat(root)
		// 	fmt.Printf("Root (mod %s) %d as float: %v\n", p, i+1, rootFloat)
		// }
		// fmt.Println()

		// minX, err := curve.SolveCubic()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Error finding minimum x value where y = 0: %v\n", err)
		// 	os.Exit(1)
		// }
		// fmt.Printf("Minimum x value, where y = 0: %.5f\n", minX)

		// Check y value for. x calculated
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

	fmt.Printf("Roots (mod %s): %v\n", p, roots)
	fmt.Println()
	for i, root := range roots { // max 3 roots
		rootFloat := utils.NewFloat().SetRat(root) // converted to float - just for human readability
		fmt.Printf("Root (mod %s) %d as float: %v\n", p, i+1, rootFloat)
	}
	fmt.Println()

	return nil
}

// getInput prompts the user for. input if the provided value is empty.
func getInput(name, value string) string {
	if value == "" {
		fmt.Printf("Please enter the value of %s, in Elliptic Curve: ", name)
		fmt.Scanln(&value)
	}
	return value
}
