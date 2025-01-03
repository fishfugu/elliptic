package main

import (
	"fmt"
	"math/big"
	"os"

	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

// This used to just be a way to functionally smoke tes the fact that bigarith worked
// I'm replacing it with math.big versions:
// 1) to make sure I know my conversions work
// 2) to keep it here for. math.big testing in the future
func run() error {
	logger := utils.InitialiseLogger("[BIGMATH]")

	// setup bigNums to do calcs with
	bigNum1, ok := new(big.Int).SetString("12345678901234567890", 10)
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'bigNum1' failed")
	bigNum2, ok := new(big.Int).SetString("98765432109876543210", 10)
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'bigNum2' failed")
	two, ok := new(big.Int).SetString("2", 10)
	utils.LogOnFailure(logger, ok, "SetString on big.Int 'two' failed")

	// Perform addition
	addResult := new(big.Int).Add(bigNum1, bigNum2)
	fmt.Println("Addition Result:", addResult.String())

	// Perform multiplication
	multiplyResult := new(big.Int).Mul(bigNum1, bigNum2)
	fmt.Println("Multiplication Result:", multiplyResult.String())

	// Perform subtraction
	subtractResult := new(big.Int).Sub(bigNum1, bigNum2)
	fmt.Println("Subtraction Result:", subtractResult.String())

	// Find the prime within the range of [multiplyResult, multiplyResult * 2]
	maxPrime := new(big.Int).Mul(multiplyResult, two)
	prime := utils.FindPrime(multiplyResult, maxPrime)
	fmt.Println("Prime found:", prime.String())

	// Perform division in the field defined by the prime modulus
	divideResult, err := finiteintfield.DivideInField(bigNum2, bigNum1, prime)
	utils.LogOnError(logger, err, "DivideInField caused an error", false)
	fmt.Println("Division In Field Result:", divideResult.String())

	return nil
}
