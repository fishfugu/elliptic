package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/sirupsen/logrus"

	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"
)

func main() {
	logger := utils.InitialiseLogger("[BIGMATH/MAIN]")
	logger.Debug("starting function main")

	err := run(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

// Smoke test for basic math operations with big.Int
func run(logger *logrus.Logger) error {
	logger.Debug("starting function run")

	// Setup bigNums to do calculations with
	bigNum1, ok := new(big.Int).SetString("12345678901234567890", 10)
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'bigNum1' failed")
	bigNum2, ok := new(big.Int).SetString("98765432109876543210", 10)
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'bigNum2' failed")

	// Perform addition
	addResult := new(big.Int).Add(bigNum1, bigNum2)
	fmt.Println("Addition Result:", addResult.String())

	// Perform multiplication
	multiplyResult := new(big.Int).Mul(bigNum1, bigNum2)
	fmt.Println("Multiplication Result:", multiplyResult.String())

	// Perform subtraction
	subtractResult := new(big.Int).Sub(bigNum1, bigNum2)
	fmt.Println("Subtraction Result:", subtractResult.String())

	// Find the prime within the range [multiplyResult, multiplyResult * 2]
	maxPrime := new(big.Int).Lsh(multiplyResult, 1) // Efficiently multiply by 2
	prime := utils.FindPrime(multiplyResult, maxPrime)
	if prime == nil {
		logger.Errorf("No prime found in range [%s, %s]", multiplyResult.String(), maxPrime.String())
		return fmt.Errorf("failed to find a prime")
	}
	fmt.Println("Prime found:", prime.String())

	// Perform division in the field defined by the prime modulus
	divideResult, err := finiteintfield.DivideInField(bigNum2, bigNum1, prime)
	if err != nil {
		utils.LogOnError(logger, err, "DivideInField caused an error", false)
		return err
	}
	fmt.Println("Division In Field Result:", divideResult.String())

	return nil
}
