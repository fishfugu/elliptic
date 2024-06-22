package main

import (
	"fmt"

	"elliptic/pkg/bigarith"
)

func main() {
	result, err := bigarith.Add("12345678901234567890", "98765432109876543210")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Addition Result:", result)
	}

	result, err = bigarith.Multiply("12345678901234567890", "98765432109876543210")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Multiplication Result:", result)
	}
	multiplyResult := result

	result, err = bigarith.Subtract("12345678901234567890", "98765432109876543210")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Subtraction Result:", result)
	}

	maxPrime, err := bigarith.Multiply(multiplyResult, "2")
	if err != nil {
		fmt.Println("Error:", err)
	}
	prime, err := bigarith.FindPrime(multiplyResult, maxPrime)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Prime found:", prime)

	result, err = bigarith.DivideInField("12345678901234567890", "98765432109876543210", prime)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Division In Field Result:", result)
	}
}
