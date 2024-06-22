package bigarith

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define a suite for all arithmetic tests
type ArithTestSuite struct {
	suite.Suite
}

// TestAdd tests the addition of two large numbers, including edge cases.
func (suite *ArithTestSuite) TestAdd() {
	tests := []struct {
		a, b   string
		want   string
		hasErr bool
	}{
		{"12345678901234567890", "98765432109876543210", "111111111011111111100", false},
		{"-12345678901234567890", "12345678901234567890", "0", false},
		{"12345678901234567890", "-98765432109876543210", "-86419753208641975320", false},
		{"0", "0", "0", false},
		{"123", "abc", "", true},
	}

	for _, tt := range tests {
		got, err := Add(tt.a, tt.b)
		correctResult := ""
		if tt.hasErr {
			assert.Error(suite.T(), err, "Add should return an error but did not")
			correctResult = fmt.Sprintf("produces an error of %s", err)
		} else {
			assert.NoError(suite.T(), err, "Add should not return an error")
			assert.Equal(suite.T(), tt.want, got, "Add did not return expected result")
			correctResult = fmt.Sprintf("= %s", got)
		}
		fmt.Printf("\tCONFIRMED: %s plus %s %s\n", tt.a, tt.b, correctResult)
	}
}

// TestSubtract tests the subtraction of two numbers, including edge cases.
func (suite *ArithTestSuite) TestSubtract() {
	tests := []struct {
		a, b   string
		want   string
		hasErr bool
	}{
		{"98765432109876543210", "12345678901234567890", "86419753208641975320", false},
		{"12345678901234567890", "98765432109876543210", "-86419753208641975320", false},
		{"-100", "100", "-200", false},
		{"100", "-100", "200", false},
		{"abc", "100", "", true},
		{"100", "abc", "", true},
	}

	for _, tt := range tests {
		got, err := Subtract(tt.a, tt.b)
		correctResult := ""
		if tt.hasErr {
			assert.Error(suite.T(), err, "Subtract should return an error but did not")
			correctResult = fmt.Sprintf("produces an error of %s", err)
		} else {
			assert.NoError(suite.T(), err, "Subtract should not return an error")
			assert.Equal(suite.T(), tt.want, got, "Subtract did not return expected result")
			correctResult = fmt.Sprintf("= %s", got)
		}
		fmt.Printf("\tCONFIRMED: %s minus %s %s\n", tt.a, tt.b, correctResult)
	}
}

// TestMultiply tests the multiplication of two numbers, including edge cases.
func (suite *ArithTestSuite) TestMultiply() {
	tests := []struct {
		a, b    string
		want    string
		wantErr string
		hasErr  bool
	}{
		{"123456789", "987654321", "121932631112635269", "", false},
		{"-123456789", "987654321", "-121932631112635269", "", false},
		{"123456789", "-987654321", "-121932631112635269", "", false},
		{"0", "123456789", "0", "", false},
		{"123456789", "0", "0", "", false},
		{"abc", "123", "", "invalid input: a = abc, b = 123 - cannot create all the integers required, from this input", true},
	}

	for _, tt := range tests {
		got, err := Multiply(tt.a, tt.b)
		strResult := ""
		strWant := ""
		if tt.hasErr {
			strResult = fmt.Sprintf("produces an error of '%s' (with a result of %s)", err, got)
			strWant = fmt.Sprintf("error of '%s'", tt.wantErr)
			assert.Error(suite.T(), err, "Multiply should return an error but did not")
			assert.Equal(suite.T(), tt.wantErr, err.Error(), "Divide did not return expected error")
		} else {
			strResult = fmt.Sprintf("= %s", got)
			strWant = tt.want
			assert.NoError(suite.T(), err, "Multiply should not return an error")
			assert.Equal(suite.T(), tt.want, got, "Multiply did not return expected result")
		}
		fmt.Printf("\tTEST & RESULT: %s multiplied by %s %s (expected %s)\n", tt.a, tt.b, strResult, strWant)
	}
}

// TestDivide tests the division of two numbers, including edge cases like division by zero.
func (suite *ArithTestSuite) TestDivide() {
	tests := []struct {
		a, b    string
		want    string
		wantErr string
		hasErr  bool
	}{
		{"987654321", "123456789", "", "division not implemented - due to ambiguous integer results", true}, // Always error - due to implementation
	}

	for _, tt := range tests {
		got, err := Divide(tt.a, tt.b)
		strResult := ""
		strWant := ""
		if tt.hasErr {
			strResult = fmt.Sprintf("produces an error of '%s' (with a result of %s)", err, got)
			strWant = fmt.Sprintf("error of '%s'", tt.wantErr)
			assert.Error(suite.T(), err, "Divide should return an error but did not")
			assert.Equal(suite.T(), tt.wantErr, err.Error(), "Divide did not return expected error")
		} else {
			strResult = fmt.Sprintf("= %s", got)
			strWant = tt.want
			assert.NoError(suite.T(), err, "Divide should not return an error")
			assert.Equal(suite.T(), tt.want, got, "Divide did not return expected result")
		}
		fmt.Printf("\tTEST & RESULT: %s divided by %s %s (expected %s)\n", tt.a, tt.b, strResult, strWant)
	}
}

// TestDivideInField tests the division of two numbers in a finite field, including edge cases.
func (suite *ArithTestSuite) TestDivideInField() {
	tests := []struct {
		a, b, p string
		want    string
		wantErr string
		hasErr  bool
	}{
		{"10", "3", "13", "12", "", false}, // Inverse of 3 mod 13 is 9; 10 * 9 mod 13 = 12
		{"3", "6", "17", "9", "", false},   // Inverse of 6 mod 17 is 3; 3 * 3 mod 17 = 9
		{"3", "0", "17", "", "error finding inverse: invalid input: a is ZERO - no modular multiplicative inverse", true},                // Division by zero scenario
		{"2", "4", "15", "", "error finding inverse: invalid input: modulus 15 is not prime", true},                                      // Non-prime modulus, might not have an inverse
		{"abc", "3", "13", "", "invalid input: a = abc, b = 3, p = 13 - cannot create all the integers required, from this input", true}, // Invalid numeric input
	}

	for _, tt := range tests {
		got, err := DivideInField(tt.a, tt.b, tt.p)
		strResult := ""
		strWant := ""
		if tt.hasErr {
			strResult = fmt.Sprintf("produces an error of '%s' (with a result of %s)", err, got)
			strWant = fmt.Sprintf("error of '%s'", tt.wantErr)
			assert.Error(suite.T(), err, "DivideInField should return an error but did not")
			assert.Equal(suite.T(), tt.wantErr, err.Error(), "DivideInField did not return expected error")
		} else {
			strResult = fmt.Sprintf("= %s", got)
			strWant = tt.want
			assert.NoError(suite.T(), err, "DivideInField should not return an error")
			assert.Equal(suite.T(), tt.want, got, fmt.Sprintf("DivideInField did not return expected result - %s divided by %s modulo %s returned %s but was expected to be %s\n", tt.a, tt.b, tt.p, got, tt.want))
		}
		fmt.Printf("\tTEST & RESULT: %s divided by %s modulo %s %s (expected %s)\n", tt.a, tt.b, tt.p, strResult, strWant)
	}
}

// TestModularInverse tests finding the modular inverse, including edge cases.
func (suite *ArithTestSuite) TestModularInverse() {
	tests := []struct {
		a, p    string
		want    string
		wantErr string
		hasErr  bool
	}{
		{"3", "13", "9", "", false},                                     // Inverse of 3 mod 13 is 9
		{"6", "17", "3", "", false},                                     // Inverse of 6 mod 17 is 3
		{"3", "18", "", "invalid input: modulus 18 is not prime", true}, // 18 is not prime, inverse might not exist
		{"abc", "13", "", "invalid input: a = abc, p = 13 - cannot create all the integers required, from this input", true}, // Invalid input
		{"0", "13", "", "invalid input: a is ZERO - no modular multiplicative inverse", true},                                // Invalid input
	}

	for _, tt := range tests {
		got, err := ModularInverse(tt.a, tt.p)
		strResult := ""
		strWant := ""
		if tt.hasErr {
			strResult = fmt.Sprintf("produces an error of '%s' (with a result of %s)", err, got)
			strWant = fmt.Sprintf("error of '%s'", tt.wantErr)
			assert.Error(suite.T(), err, "ModularInverse should return an error but did not")
			assert.Equal(suite.T(), tt.wantErr, err.Error(), "ModularInverse did not return expected error")
		} else {
			strResult = fmt.Sprintf("= %s", got)
			strWant = tt.want
			assert.NoError(suite.T(), err, "ModularInverse should not return an error")
			assert.Equal(suite.T(), tt.want, got, "ModularInverse did not return expected result")
		}
		fmt.Printf("\tTEST & RESULT: Inverse of %s modulo %s %s (expected %s)\n", tt.a, tt.p, strResult, strWant)
	}
}

// SetupTestSuite runs before the tests in the suite are executed.
func (suite *ArithTestSuite) SetupTest() {
	// Initialize anything required for the test suite
}

// AfterTest runs after each test in the suite.
func (suite *ArithTestSuite) AfterTest(_, _ string) {
	// Clean up after each test if necessary
}

// This function hooks up testify's suite with the Go testing framework.
func TestArithTestSuite(t *testing.T) {
	suite.Run(t, new(ArithTestSuite))
}
