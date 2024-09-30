package main

import (
	"fmt"

	"elliptic/pkg/bigarith"
)

func main() {
	// Perform addition
	addResult := bigarith.NewInt("12345678901234567890").Plus(bigarith.NewInt("98765432109876543210").Val())
	fmt.Println("Addition Result:", addResult.Val())

	// Perform multiplication
	multiplyResult := bigarith.NewInt("12345678901234567890").Times(bigarith.NewInt("98765432109876543210").Val())
	fmt.Println("Multiplication Result:", multiplyResult.Val())

	// Perform subtraction
	subtractResult := bigarith.NewInt("12345678901234567890").Minus(bigarith.NewInt("98765432109876543210").Val())
	fmt.Println("Subtraction Result:", subtractResult.Val())

	// Find the prime within the range of [multiplyResult, multiplyResult * 2]
	maxPrime := multiplyResult.Times("2")
	prime := multiplyResult.FindPrime(maxPrime.Val())
	fmt.Println("Prime found:", prime.Val())

	// Perform division in the field defined by the prime modulus
	divideResult := bigarith.NewInt("98765432109876543210").DivideInField("12345678901234567890", prime.Val())
	fmt.Println("Division In Field Result:", divideResult.Val())
}
