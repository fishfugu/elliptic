package ellipticcurve_test

import (
	"elliptic/pkg/ellipticcurve"
	"fmt"
	"math/big"
	"testing"
)

// Define a struct to hold the test case inputs and expected output
type cubicTestCase struct {
	A             *big.Int   // Coefficient for x term
	B             *big.Int   // Constant term
	expectedRoots []*big.Rat // List of expected roots
}

// Test function for multiple test cases
func TestSolveCubic_OneRealRoot(t *testing.T) {
	// Define a list of test cases
	testCases := []cubicTestCase{
		{
			A: new(big.Int).SetInt64(0),
			B: new(big.Int).SetInt64(1),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-1),
			},
		},
		{
			A: new(big.Int).SetInt64(0),
			B: new(big.Int).SetInt64(8),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		},
		{
			A: new(big.Int).SetInt64(0),
			B: new(big.Int).SetInt64(27),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-3),
			},
		},
		{
			A: new(big.Int).SetInt64(0),
			B: new(big.Int).SetInt64(64),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-4),
			},
		},
		{
			A: new(big.Int).SetInt64(0),
			B: new(big.Int).SetInt64(125),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-5),
			},
		}, {
			A: new(big.Int).SetInt64(16),
			B: new(big.Int).SetInt64(40),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		}, {
			A: new(big.Int).SetInt64(36),
			B: new(big.Int).SetInt64(80),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		}, {
			A: new(big.Int).SetInt64(49),
			B: new(big.Int).SetInt64(106),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		}, {
			A: new(big.Int).SetInt64(49),
			B: new(big.Int).SetInt64(174),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-3),
			},
		}, {
			A: new(big.Int).SetInt64(49),
			B: new(big.Int).SetInt64(260),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-4),
			},
		}, {
			A: new(big.Int).SetInt64(64),
			B: new(big.Int).SetInt64(219),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-3),
			},
		},
	}

	// testCases := []cubicTestCase{
	// {
	// 	A:                     "49",
	// 	B:                     "260",
	// 	expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 	expectedRoots:         []string{"-4"},
	// },
	// {
	// 	A:                     "64",
	// 	B:                     "219",
	// 	expectedNumberOfRoots: 1, // Expected single real root at x = -1 for equation x^3 + 0x + 64 = 0
	// 	expectedRoots:         []string{"-3"},
	// },
	// }

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Create a new elliptic curve object with A and B
			ec := ellipticcurve.NewEllipticCurve(tc.A, tc.B)

			// Solve the cubic equation
			roots, err := ec.SolveCubic()
			if err != nil {
				t.Fatalf("Error solving cubic: %v", err)
			}

			// Check if the number of real roots matches the expected value
			if len(roots) != len(tc.expectedRoots) {
				t.Errorf("Expected %d real root(s), got %d", len(tc.expectedRoots), len(roots))
			}

			// Check if real roots match the expected value
			for i := range roots {
				foundMatch := false
				lastErrorMessage := "NONE"
				for j := range tc.expectedRoots {
					if roots[i].Cmp(tc.expectedRoots[j]) == 0 {
						foundMatch = true
					} else {
						diff := new(big.Rat).Abs(new(big.Rat).Sub(roots[i], tc.expectedRoots[j]))
						diffFloat := new(big.Float).SetRat(diff)
						lastErrorMessage = fmt.Sprintf("Expedcted first root %s, got %s, diff %s (as a float: %v)", tc.expectedRoots[j], roots[i], diff, diffFloat)
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
			A: new(big.Int).SetInt64(-48),
			B: new(big.Int).SetInt64(128),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-8),
				new(big.Rat).SetInt64(4),
				new(big.Rat).SetInt64(4),
			},
		},
		// {
		// 	A: new(big.Int).SetInt64(-7),
		// 	B: new(big.Int).SetInt64(-6),
		// 	expectedRoots: []*big.Rat{
		// 		new(big.Rat).SetInt64(-2),
		// 		new(big.Rat).SetInt64(-1),
		// 		new(big.Rat).SetInt64(3),
		// 	},
		// },
		// {
		// 	A: new(big.Int).SetInt64(-75),
		// 	B: new(big.Int).SetInt64(20),
		// 	expectedRoots: []*big.Rat{
		// 		new(big.Rat).SetInt64(-10),
		// 		new(big.Rat).SetInt64(5),
		// 		new(big.Rat).SetInt64(5),
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Create a new elliptic curve object with A and B
			ec := ellipticcurve.NewEllipticCurve(tc.A, tc.B)

			// Solve the cubic equation
			roots, err := ec.SolveCubic()
			if err != nil {
				t.Fatalf("Error solving cubic: %v", err)
			}

			// Check if the number of real roots matches the expected value
			if len(roots) != len(tc.expectedRoots) {
				t.Errorf("Expected %d real root(s), got %d", len(tc.expectedRoots), len(roots))
			}

			// Check if real roots match the expected value
			for i := range roots {
				foundMatch := false
				lastErrorMessage := "NONE"
				for j := range tc.expectedRoots {
					if roots[i].Cmp(tc.expectedRoots[j]) == 0 {
						foundMatch = true
					} else {
						diff := new(big.Rat).Abs(new(big.Rat).Sub(roots[i], tc.expectedRoots[j]))
						lastErrorMessage = fmt.Sprintf("Expedcted roots: %s, got %s, diff %s", tc.expectedRoots, roots[i], diff)
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
		// {
		// 	A:             "-28",
		// 	B:             "48",
		// 	expectedRoots: []string{"-6", "2", "4"},
		// },
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Create a new elliptic curve object with A and B
			ec := ellipticcurve.NewEllipticCurve(tc.A, tc.B)

			// Solve the cubic equation
			roots, err := ec.SolveCubic()
			if err != nil {
				t.Fatalf("Error solving cubic: %v", err)
			}

			// Check if the number of real roots matches the expected value
			if len(roots) != len(tc.expectedRoots) {
				t.Errorf("Expected %d real root(s), got %d", len(tc.expectedRoots), len(roots))
			}

			// Check if real roots match the expected value
			for i := range roots {
				foundMatch := false
				lastErrorMessage := "NONE"
				for j := range tc.expectedRoots {
					if roots[i].Cmp(tc.expectedRoots[j]) == 0 {
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
