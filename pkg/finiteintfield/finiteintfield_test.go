package finiteintfield

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"elliptic/pkg/ellipticcurve"
)

// Test suite
type FiniteIntFieldSuite struct {
	suite.Suite
}

func (suite *FiniteIntFieldSuite) TestCalculatePoints() {
	A := big.NewInt(2)
	B := big.NewInt(3)
	P := big.NewInt(5)
	ffec := ellipticcurve.NewFiniteFieldEC(A, B, P)

	expectedPoints := [][2]string{
		{"1", "1"}, {"1", "4"}, {"2", "0"}, {"3", "1"}, {"3", "4"}, {"4", "0"},
	}

	points, err := CalculatePoints(*ffec)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), points, len(expectedPoints))

	pointMap := make(map[[2]string]bool)
	for _, point := range points {
		pointMap[point] = true
	}

	for _, expectedPoint := range expectedPoints {
		assert.True(suite.T(), pointMap[expectedPoint], "Expected point %v not found", expectedPoint)
	}
}

func (suite *FiniteIntFieldSuite) TestFormatPoints() {
	points := [][2]string{
		{"0", "2"}, {"1", "1"}, {"2", "4"},
	}
	expectedOutput := `List of Points on the Curve:
Point (x, y): (0, 2)
Point (x, y): (1, 1)
Point (x, y): (2, 4)
`

	formattedPoints := FormatPoints(points)
	assert.Equal(suite.T(), expectedOutput, formattedPoints)
}

func (suite *FiniteIntFieldSuite) TestVisualisePoints() {
	points := [][2]string{
		{"0", "2"}, {"1", "1"}, {"2", "4"},
	}
	p := 5
	expectedOutput := `
2D Plane Visualization with Cartesian Axes, Reflection Line, and Scale:
|            5
|   *        4
| . . . . .  5/2 3
*            2
| *          1
+ / / / / /  0
`
	visualisedPoints := VisualisePoints(points, p)
	assert.Equal(suite.T(), expectedOutput, visualisedPoints)
}

func (suite *FiniteIntFieldSuite) TestMax() {
	int := max(1, 2)
	suite.Require().Equal(int, 2)

	int = max(2, 1)
	suite.Require().Equal(int, 2)

	int = max(3, 3)
	suite.Require().Equal(int, 3)
}

func TestFiniteIntFieldSuite(t *testing.T) {
	suite.Run(t, new(FiniteIntFieldSuite))
}
