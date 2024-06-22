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
}
