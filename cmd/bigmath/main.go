package main

import (
	"fmt"
	"os"

	ba "elliptic/pkg/bigarith"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

func run() error {
	// Perform addition
	addResult := ba.NewInt("12345678901234567890").Plus(ba.NewInt("98765432109876543210").Val())
	fmt.Println("Addition Result:", addResult.Val())

	// Perform multiplication
	multiplyResult := ba.NewInt("12345678901234567890").Times(ba.NewInt("98765432109876543210").Val())
	fmt.Println("Multiplication Result:", multiplyResult.Val())

	// Perform subtraction
	subtractResult := ba.NewInt("12345678901234567890").Minus(ba.NewInt("98765432109876543210").Val())
	fmt.Println("Subtraction Result:", subtractResult.Val())

	// Find the prime within the range of [multiplyResult, multiplyResult * 2]
	maxPrime := multiplyResult.Times("2")
	prime := multiplyResult.FindPrime(maxPrime.Val())
	fmt.Println("Prime found:", prime.Val())

	// Perform division in the field defined by the prime modulus
	divideResult := ba.NewInt("98765432109876543210").DivideInField("12345678901234567890", prime.Val())
	fmt.Println("Division In Field Result:", divideResult.Val())

	return nil
}
