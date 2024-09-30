package ellipticcurve_test

import (
	"elliptic/pkg/bigarith"
	"elliptic/pkg/ellipticcurve"
	"math/rand"
	"strconv"
	"testing"
)

func TestSolveCubic_OneRealRoot(t *testing.T) {
	// Test for the case where there is 1 real root (discriminant > 0)
	A := bigarith.NewInt("6")
	B := bigarith.NewInt("10")
	ec := ellipticcurve.NewEllipticCurve(A, B)

	roots, err := ec.SolveCubic()
	if err != nil {
		t.Fatalf("Error solving cubic: %v", err)
	}

	// There should be only 1 root
	if len(roots) != 1 {
		t.Errorf("Expected 1 real root, got %d", len(roots))
	}
}

func TestSolveCubic_ThreeRealRoots(t *testing.T) {
	// Test for the case where there are 3 real roots (discriminant < 0)
	A := bigarith.NewInt("-3")
	B := bigarith.NewInt("-1")
	ec := ellipticcurve.NewEllipticCurve(A, B)

	roots, err := ec.SolveCubic()
	if err != nil {
		t.Fatalf("Error solving cubic: %v", err)
	}

	// There should be 3 roots
	if len(roots) != 3 {
		t.Errorf("Expected 3 real roots, got %d", len(roots))
	}
}

func TestSolveCubic_TwoEqualRealRoots(t *testing.T) {
	// Test for the case where there are 2 equal real roots (discriminant = 0)
	A := bigarith.NewInt("2")
	B := bigarith.NewInt("1")
	ec := ellipticcurve.NewEllipticCurve(A, B)

	roots, err := ec.SolveCubic()
	if err != nil {
		t.Fatalf("Error solving cubic: %v", err)
	}

	// There should be 2 equal real roots
	if len(roots) != 3 {
		t.Errorf("Expected 3 real roots (2 equal), got %d", len(roots))
	}

	if roots[1] != roots[2] {
		t.Errorf("Expected two equal roots, but got %s and %s", roots[1], roots[2])
	}
}

func TestSolveCubic_RandomSmallIntegers(t *testing.T) {
	// Test with random small integers to validate the function over a wide range of inputs
	for i := 0; i < 50; i++ {
		A := bigarith.NewInt(randomInt())
		B := bigarith.NewInt(randomInt())
		ec := ellipticcurve.NewEllipticCurve(A, B)

		roots, err := ec.SolveCubic()
		if err != nil {
			t.Errorf("Error solving cubic for A = %s, B = %s: %v", A.Val(), B.Val(), err)
		}

		// Validate that at least one real root is returned in all cases
		if len(roots) < 1 {
			t.Errorf("Expected at least 1 real root, got %d for A = %s, B = %s", len(roots), A.Val(), B.Val())
		}
	}
}

// Helper function to generate small random integers
func randomInt() string {
	return strconv.Itoa(rand.Intn(20) - 10) // Random integers between -10 and 10
}
