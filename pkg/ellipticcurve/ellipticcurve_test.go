package ellipticcurve_test

import (
	"elliptic/pkg/bigarith"
	"elliptic/pkg/ellipticcurve"
	"fmt"
	"testing"
)

// Define a struct to hold the test case inputs and expected output
type cubicTestCase struct {
	A                     string   // Coefficient for x term
	B                     string   // Constant term
	expectedNumberOfRoots int      // Expected number of real roots
	expectedRoots         []string // List of expected roots
}

// Test function for multiple test cases
func TestSolveCubic_OneRealRoot(t *testing.T) {
	// Define a list of test cases
	// testCases := []cubicTestCase{
	// 	{
	// 		A:                     "0",
	// 		B:                     "1",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 1 = 0
	// 		expectedRoots:         []string{"-1"},
	// 	}, {
	// 		A:                     "0",
	// 		B:                     "8",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 8 = 0
	// 		expectedRoots:         []string{"-2"},
	// 	}, {
	// 		A:                     "0",
	// 		B:                     "27",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 27 = 0
	// 		expectedRoots:         []string{"-3"},
	// 	}, {
	// 		A:                     "0",
	// 		B:                     "64",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-4"},
	// 	}, {
	// 		A:                     "0",
	// 		B:                     "125",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-5"},
	// 	}, {
	// 		A:                     "16",
	// 		B:                     "40",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-2"},
	// 	}, {
	// 		A:                     "36",
	// 		B:                     "80",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-2"},
	// 	}, {
	// 		A:                     "49",
	// 		B:                     "106",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-2"},
	// 	}, {
	// 		A:                     "49",
	// 		B:                     "174",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-3"},
	// 	}, {
	// 		A:                     "49",
	// 		B:                     "260",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-4"},
	// 	}, {
	// 		A:                     "64",
	// 		B:                     "219",
	// 		expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 		expectedRoots:         []string{"-3"},
	// 	},
	// }

	testCases := []cubicTestCase{
		{
			A:                     "49",
			B:                     "260",
			expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
			expectedRoots:         []string{"-4"},
		}, {
			A:                     "64",
			B:                     "219",
			expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
			expectedRoots:         []string{"-3"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Convert strings to big integers
			A := bigarith.NewInt(tc.A)
			B := bigarith.NewInt(tc.B)

			// Create a new elliptic curve object with A and B
			ec := ellipticcurve.NewEllipticCurve(A, B)

			// Solve the cubic equation
			roots, err := ec.SolveCubic()
			if err != nil {
				t.Fatalf("Error solving cubic: %v", err)
			}

			// Check if the number of real roots matches the expected value
			if len(roots) != tc.expectedNumberOfRoots {
				t.Errorf("Expected %d real root(s), got %d", tc.expectedNumberOfRoots, len(roots))
			}

			// Check if real roots match the expected value
			for i := range roots {
				foundMatch := false
				lastErrorMessage := "NONE"
				for j := range tc.expectedRoots {
					if bigarith.NewFloat(roots[i]).Compare(tc.expectedRoots[j]) == 0 {
						foundMatch = true
					} else {
						lastErrorMessage = fmt.Sprintf("Expedcted first root %s, got %s", tc.expectedRoots[j], roots[i])
					}
				}
				if !foundMatch {
					t.Errorf(lastErrorMessage)
				}
			}
		})
	}
}

// Test function for multiple test cases
func TestSolveCubic_DoubleRoot(t *testing.T) {
	// Define a list of test cases
	testCases := []cubicTestCase{
		{
			A:                     "-48",
			B:                     "128",
			expectedNumberOfRoots: 3, // Expected single real root at x = -1 for equation x^3 + 0x + 1 = 0
			expectedRoots:         []string{"-8", "4"},
		}, {
			A:                     "-75",
			B:                     "250",
			expectedNumberOfRoots: 3, // Expected single real root at x = -1 for equation x^3 + 0x + 1 = 0
			expectedRoots:         []string{"-10", "5"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Convert strings to big integers
			A := bigarith.NewInt(tc.A)
			B := bigarith.NewInt(tc.B)

			// Create a new elliptic curve object with A and B
			ec := ellipticcurve.NewEllipticCurve(A, B)

			// Solve the cubic equation
			roots, err := ec.SolveCubic()
			if err != nil {
				t.Fatalf("Error solving cubic: %v", err)
			}

			// Check if the number of real roots matches the expected value
			if len(roots) != tc.expectedNumberOfRoots {
				t.Errorf("Expected %d real root(s), got %d", tc.expectedNumberOfRoots, len(roots))
			}

			// Check if real roots match the expected value
			for i := range roots {
				foundMatch := false
				lastErrorMessage := "NONE"
				for j := range tc.expectedRoots {
					if bigarith.NewFloat(roots[i]).Compare(tc.expectedRoots[j]) == 0 {
						foundMatch = true
					} else {
						lastErrorMessage = fmt.Sprintf("Expedcted roots: %s, got %s", tc.expectedRoots, roots[i])
					}
				}
				if !foundMatch {
					t.Errorf(lastErrorMessage)
				}
			}
		})
	}
}

func TestSolveCubic_ThreeRealRoots(t *testing.T) {
	// Define a list of test cases
	testCases := []cubicTestCase{
		{
			A:                     "-28",
			B:                     "48",
			expectedNumberOfRoots: 3, // Expected single real root at x = -1 for equation x^3 + 0x + 1 = 0
			expectedRoots:         []string{"-6", "2", "4"},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Convert strings to big integers
			A := bigarith.NewInt(tc.A) // e.g. -28
			B := bigarith.NewInt(tc.B) // e.g. 48

			// Create a new elliptic curve object with A and B
			ec := ellipticcurve.NewEllipticCurve(A, B)

			// Solve the cubic equation
			roots, err := ec.SolveCubic()
			if err != nil {
				t.Fatalf("Error solving cubic: %v", err)
			}

			// Check if the number of real roots matches the expected value
			if len(roots) != tc.expectedNumberOfRoots {
				t.Errorf("Expected %d real root(s), got %d", tc.expectedNumberOfRoots, len(roots))
			}

			// Check if real roots match the expected value
			for i := range roots {
				foundMatch := false
				lastErrorMessage := "NONE"
				for j := range tc.expectedRoots {
					if bigarith.NewFloat(roots[i]).Compare(tc.expectedRoots[j]) == 0 {
						foundMatch = true
					} else {
						lastErrorMessage = fmt.Sprintf("Expedcted roots: %s, got %s", tc.expectedRoots, roots[i])
					}
				}
				if !foundMatch {
					t.Errorf(lastErrorMessage)
				}
			}
		})
	}
}
