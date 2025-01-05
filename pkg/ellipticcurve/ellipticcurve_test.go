package ellipticcurve

import (
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

func BenchmarkSolveCubic(b *testing.B) {
	// Example coefficients for benchmarking
	testCases := []struct {
		A, B *big.Int
	}{
		{new(big.Int).SetInt64(0), new(big.Int).SetInt64(0)},
		{new(big.Int).SetInt64(0), new(big.Int).SetInt64(1)},
		{new(big.Int).SetInt64(1), new(big.Int).SetInt64(1)},
		{new(big.Int).SetInt64(16), new(big.Int).SetInt64(40)},
		{new(big.Int).SetInt64(-48), new(big.Int).SetInt64(128)},
		{new(big.Int).SetInt64(-75), new(big.Int).SetInt64(250)},
		{new(big.Int).SetInt64(-12), new(big.Int).SetInt64(16)},
		{new(big.Int).SetInt64(-28), new(big.Int).SetInt64(48)},
		{new(big.Int).SetInt64(49), new(big.Int).SetInt64(260)},
		{new(big.Int).SetInt64(529), new(big.Int).SetInt64(-292)},
		{new(big.Int).SetInt64(-796), new(big.Int).SetInt64(-591)},
		{new(big.Int).SetInt64(-511), new(big.Int).SetInt64(392)},
		{new(big.Int).SetInt64(-914), new(big.Int).SetInt64(-245)},
		{new(big.Int).SetInt64(-588), new(big.Int).SetInt64(-908)},
		{new(big.Int).SetInt64(-600), new(big.Int).SetInt64(40000)},
	}

	for _, tc := range testCases {
		b.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(b *testing.B) {
			// Create the elliptic curve
			ec := NewEllipticCurve(tc.A, tc.B)

			// Reset the timer to exclude setup time
			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				_, err := ec.SolveCubic()
				if err != nil {
					b.Fatalf("Error solving cubic: %v", err)
				}
			}
		})
	}
}

func BenchmarkSolveCubicDetails(b *testing.B) {
	testCases := []struct {
		A, B *big.Int
	}{
		{new(big.Int).SetInt64(0), new(big.Int).SetInt64(0)},
		{new(big.Int).SetInt64(0), new(big.Int).SetInt64(1)},
		{new(big.Int).SetInt64(1), new(big.Int).SetInt64(1)},
		{new(big.Int).SetInt64(16), new(big.Int).SetInt64(40)},
		{new(big.Int).SetInt64(-48), new(big.Int).SetInt64(128)}, // Example with double root
		{new(big.Int).SetInt64(-75), new(big.Int).SetInt64(250)}, // Example with double root
		{new(big.Int).SetInt64(-12), new(big.Int).SetInt64(16)},  // Example with double root
	}

	for _, tc := range testCases {
		ec := NewEllipticCurve(tc.A, tc.B)

		b.Run(fmt.Sprintf("Discriminant A=%s, B=%s", tc.A, tc.B), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				A, B := ec.GetDetailsAsRats()
				aCubed := new(big.Rat).Mul(A, new(big.Rat).Mul(A, A))
				bSquared := new(big.Rat).Mul(B, B)
				fourACubed := new(big.Rat).Mul(new(big.Rat).SetInt64(4), aCubed)
				twentySevenBCubed := new(big.Rat).Mul(new(big.Rat).SetInt64(27), bSquared)
				_ = new(big.Rat).Add(fourACubed, twentySevenBCubed) // Calculate discriminant
			}
		})

		b.Run(fmt.Sprintf("Newton A=%s, B=%s", tc.A, tc.B), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				A, B := ec.GetDetails()
				_, _ = newtonCubic(A, B)
			}
		})

		b.Run(fmt.Sprintf("Double Root A=%s, B=%s", tc.A, tc.B), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				root1 := big.NewRat(4, 1) // Simulate root1
				xSquared := new(big.Rat).Mul(root1, root1)
				threeXSquared := new(big.Rat).Mul(threeRat, xSquared)
				gradient := new(big.Rat).Add(threeXSquared, new(big.Rat).SetInt(tc.A))
				_ = approximateRat(gradient) // Simulate operations
			}
		})
	}
}

func BenchmarkNewtonCubic(b *testing.B) {
	testCases := []struct {
		A, B *big.Int
	}{
		{new(big.Int).SetInt64(0), new(big.Int).SetInt64(0)},
		{new(big.Int).SetInt64(0), new(big.Int).SetInt64(1)},
		{new(big.Int).SetInt64(1), new(big.Int).SetInt64(1)},
		{new(big.Int).SetInt64(16), new(big.Int).SetInt64(40)},
		{new(big.Int).SetInt64(-48), new(big.Int).SetInt64(128)}, // Example with double root
		{new(big.Int).SetInt64(-75), new(big.Int).SetInt64(250)}, // Example with double root
		{new(big.Int).SetInt64(-12), new(big.Int).SetInt64(16)},  // Example with double root
	}

	for _, tc := range testCases {
		b.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = newtonCubic(tc.A, tc.B)
			}
		})
	}
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
		},
		{
			A: new(big.Int).SetInt64(16),
			B: new(big.Int).SetInt64(40),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		},
		{
			A: new(big.Int).SetInt64(36),
			B: new(big.Int).SetInt64(80),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		},
		{
			A: new(big.Int).SetInt64(49),
			B: new(big.Int).SetInt64(106),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
			},
		},
		{
			A: new(big.Int).SetInt64(49),
			B: new(big.Int).SetInt64(174),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-3),
			},
		},
		{
			A: new(big.Int).SetInt64(49),
			B: new(big.Int).SetInt64(260),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-4),
			},
		},
		{
			A: new(big.Int).SetInt64(64),
			B: new(big.Int).SetInt64(219),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-3),
			},
		},
		{
			A: new(big.Int).SetInt64(-600),
			B: new(big.Int).SetInt64(40000),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-40),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Create a new elliptic curve object with A and B
			ec := NewEllipticCurve(tc.A, tc.B)

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
						lastErrorMessage = fmt.Sprintf("Expedcted roots: %s, got %s, diff %s", tc.expectedRoots, roots[i].FloatString(100), diff.FloatString(100))
					}
				}
				if !foundMatch {
					t.Error(lastErrorMessage)
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
		{
			A: new(big.Int).SetInt64(-75),
			B: new(big.Int).SetInt64(250),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-10),
				new(big.Rat).SetInt64(5),
				new(big.Rat).SetInt64(5),
			},
		},
		{
			A: new(big.Int).SetInt64(-12),
			B: new(big.Int).SetInt64(16),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-4),
				new(big.Rat).SetInt64(2),
				new(big.Rat).SetInt64(2),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Create a new elliptic curve object with A and B
			ec := NewEllipticCurve(tc.A, tc.B)

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
						lastErrorMessage = fmt.Sprintf("Expedcted roots: %s, got %s, diff %s", tc.expectedRoots, roots[i].FloatString(100), diff.FloatString(100))
					}
				}
				if !foundMatch {
					t.Error(lastErrorMessage)
				}
			}
		})
	}
}

func TestSolveCubic_ThreeRealRoots(t *testing.T) {
	// Define a list of test cases
	testCases := []cubicTestCase{
		{
			A: new(big.Int).SetInt64(-28),
			B: new(big.Int).SetInt64(48),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-6),
				new(big.Rat).SetInt64(2),
				new(big.Rat).SetInt64(4),
			},
		},
		{
			A: new(big.Int).SetInt64(-7),
			B: new(big.Int).SetInt64(-6),
			expectedRoots: []*big.Rat{
				new(big.Rat).SetInt64(-2),
				new(big.Rat).SetInt64(-1),
				new(big.Rat).SetInt64(3),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("A=%s, B=%s", tc.A, tc.B), func(t *testing.T) {
			// Create a new elliptic curve object with A and B
			ec := NewEllipticCurve(tc.A, tc.B)

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
				lastErrorMessage := "NONE - should not have triggered "
				for j := range tc.expectedRoots {
					if roots[i].Cmp(tc.expectedRoots[j]) == 0 {
						foundMatch = true
					} else {
						diff := new(big.Rat).Abs(new(big.Rat).Sub(roots[i], tc.expectedRoots[j]))
						lastErrorMessage = fmt.Sprintf("Expedcted roots: %v, got %s, diff %s", tc.expectedRoots, roots[i].FloatString(100), diff.FloatString(100))
					}
				}
				if !foundMatch {
					t.Error(lastErrorMessage)
				}
			}
		})
	}
}
