package ellipticcurve

import (
	"flag"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var verbose bool

// EllipticCurveTestSuite defines a suite of tests for the EllipticCurve package
type EllipticCurveTestSuite struct {
	suite.Suite
	Verbose bool
}

func TestMain(m *testing.M) {
	// Define the verbose flag for testing.
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output for tests")
	flag.Parse()

	// Now run the tests.
	os.Exit(m.Run())
}

// SetupTestSuite runs before the tests in the suite are executed.
func (suite *EllipticCurveTestSuite) SetupTest() {
	suite.Verbose = verbose // Use the global verbose flag
}

// TestNewEllipticCurve tests the creation of a new EllipticCurve.
func (suite *EllipticCurveTestSuite) TestNewEllipticCurve() {
	a, b := 1.0, -1.0
	curve := NewEllipticCurve(a, b)
	gotA, gotB := curve.GetDetails()

	assert.Equal(suite.T(), a, gotA, "Coefficient A does not match")
	assert.Equal(suite.T(), b, gotB, "Coefficient B does not match")
}

// TestNewFiniteFieldEC tests the creation of a new FiniteFieldEC.
func (suite *EllipticCurveTestSuite) TestNewFiniteFieldEC() {
	a, b, p := big.NewInt(1), big.NewInt(-1), big.NewInt(17)
	curve := NewFiniteFieldEC(a, b, p)
	gotA, gotB, gotP := curve.GetDetails()

	assert.True(suite.T(), a.Cmp(gotA) == 0, "Coefficient A does not match")
	assert.True(suite.T(), b.Cmp(gotB) == 0, "Coefficient B does not match")
	assert.True(suite.T(), p.Cmp(gotP) == 0, "Modulus P does not match")
}

// TestInvalidParameters tests creation of curves with invalid parameters.
func (suite *EllipticCurveTestSuite) TestInvalidParameters() {
	// TODO: Can I think of anything?
}

// SetupSuite prepares any contexts needed before running the suite
func (suite *EllipticCurveTestSuite) SetupSuite() {
	// Initialization logic before the tests run
}

// TearDownSuite cleans up any contexts once the suite is done
func (suite *EllipticCurveTestSuite) TearDownSuite() {
	// Cleanup logic after all tests run
}

// This function hooks up the suite with the Go testing framework.
func TestEllipticCurveTestSuite(t *testing.T) {
	suite.Run(t, new(EllipticCurveTestSuite))
}
