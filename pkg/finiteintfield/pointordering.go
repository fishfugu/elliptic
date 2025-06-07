package finiteintfield

import (
	"math/big"
	"sort"

	"elliptic/pkg/ellipticcurve"
)

// OrderIntPoints returns the list of points handed in... ordered like this:
// Assuming we have a list of points as: [][2]*big.Int
// order those points such that, all the points where y is positive in reverse x position order
// followed by all the points where y is negative in positive x position order
func OrderIntPoints(ffec *ellipticcurve.FiniteFieldEC, points [][2]*big.Int) [][2]*big.Int {
	// separate points based on y > 0 and y < 0
	var yPositive, yNegative [][2]*big.Int
	for _, point := range points {
		if point[1].Cmp(big.NewInt(0)) >= 0 {
			yPositive = append(yPositive, point)
		} else {
			yNegative = append(yNegative, point)
		}
	}

	// sort yPositive in reverse x order
	sort.Slice(yPositive, func(i, j int) bool {
		return yPositive[i][0].Cmp(yPositive[j][0]) >= 0
	})

	// sort yNegative in positive x order
	sort.Slice(yNegative, func(i, j int) bool {
		return yNegative[i][0].Cmp(yNegative[j][0]) < 0
	})

	// concatenate the two sorted slices
	orderedPoints := append(yPositive, yNegative...)

	// return the ordered points
	return orderedPoints
}

// OrderRealPoints returns the list of points handed in... ordered like this:
// Assuming we have a list of points as: [][2]*big.Rat
// order those points such that, all the points where y is positive in reverse x position order
// followed by all the points where y is negative in positive x position order
// NOTE: position of finite field window is moved to:
func OrderRealPoints(ec *ellipticcurve.EllipticCurve, points [][2]*big.Rat) [][2]*big.Rat {
	// separate points based on y > 0, y == 0, and y < 0
	var yPositive, yNegative, yZero [][2]*big.Rat
	for _, point := range points {
		if point[1].Cmp(big.NewRat(0, 1)) > 0 {
			yPositive = append(yPositive, point)
		} else if point[1].Cmp(big.NewRat(0, 1)) == 0 {
			yZero = append(yZero, point)
		} else {
			yNegative = append(yNegative, point)
		}
	}

	// sort yPositive in reverse x order
	sort.Slice(yPositive, func(i, j int) bool {
		return yPositive[i][0].Cmp(yPositive[j][0]) > 0
	})

	// sort yNegative in positive x order
	sort.Slice(yNegative, func(i, j int) bool {
		return yNegative[i][0].Cmp(yNegative[j][0]) < 0
	})

	// concatenate the two sorted slices
	orderedPoints := append(yPositive, yZero...)
	orderedPoints = append(orderedPoints, yNegative...)

	// return the ordered points
	return orderedPoints
}
