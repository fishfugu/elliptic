// package main

// import (
// 	"fmt"
// 	"image/color"
// 	"math/big"
// 	"os"

// 	"github.com/sirupsen/logrus"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/canvas"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"

// 	"elliptic/pkg/ellipticcurve"
// 	"elliptic/pkg/finiteintfield"
// 	"elliptic/pkg/utils"
// )

// func main() {
// 	logger := utils.InitialiseLogger("[ECVIS]")
// 	logger.Debug("[ECVIS] starting function main")

// 	err := run(logger)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
// 		panic(fmt.Sprintf("panic, error: %v\n", err))
// 	}
// }

// func run(logger *logrus.Logger) error {
// 	logger.Debug("[ECVIS] starting function run")

// 	myApp := app.New()
// 	finiteFieldVisWindow := myApp.NewWindow("Elliptic Curve Visualisation in Finite Field")

// 	// Default curve parameters
// 	a, b, p := big.NewInt(-2), big.NewInt(1), big.NewInt(13)

// 	// Create input fields for A and B
// 	aEntry := widget.NewEntry()
// 	aEntry.SetText(a.String())
// 	aEntry.Validator = func(s string) error {
// 		_, ok := new(big.Int).SetString(s, 10)
// 		if !ok {
// 			return fmt.Errorf("invalid integer")
// 		}
// 		return nil
// 	}

// 	bEntry := widget.NewEntry()
// 	bEntry.SetText(b.String())
// 	bEntry.Validator = func(s string) error {
// 		_, ok := new(big.Int).SetString(s, 10)
// 		if !ok {
// 			return fmt.Errorf("invalid integer")
// 		}
// 		return nil
// 	}

// 	pEntry := widget.NewEntry()
// 	pEntry.SetText(p.String())
// 	pEntry.Validator = func(s string) error {
// 		_, ok := new(big.Int).SetString(s, 10)
// 		if !ok {
// 			return fmt.Errorf("invalid integer")
// 		}
// 		return nil
// 	}

// 	graphContent := container.NewWithoutLayout()

// 	// Function to redraw the graph
// 	redraw := func() {
// 		// Parse A and B values from inputs
// 		aValue, ok := new(big.Int).SetString(aEntry.Text, 10)
// 		if !ok {
// 			logger.Error("Invalid value for A")
// 			return
// 		}

// 		bValue, ok := new(big.Int).SetString(bEntry.Text, 10)
// 		if !ok {
// 			logger.Error("Invalid value for B")
// 			return
// 		}

// 		pValue, ok := new(big.Int).SetString(pEntry.Text, 10)
// 		if !ok {
// 			logger.Error("Invalid value for B")
// 			return
// 		}

// 		// Create the curve
// 		curve := ellipticcurve.NewFiniteFieldEC(aValue, bValue, pValue)

// 		// Calculate finitePoints
// 		finitePoints, _, err := finiteintfield.CalculatePoints(curve)
// 		if err != nil {
// 			logger.Error("Error calculating finitePoints:", err)
// 			return
// 		}

// 		// Clear and redraw finitePoints
// 		graphContent.Objects = nil
// 		for _, point := range finitePoints {
// 			x := new(big.Float).SetInt(point[0])
// 			y := new(big.Float).SetInt(point[1])

// 			// Convert to canvas coordinates
// 			xCanvas, _ := (x.Add(x, big.NewFloat(6)).Quo(x, big.NewFloat(13))).Float64()
// 			xCanvas *= 600
// 			yCanvas, _ := (y.Add(y, big.NewFloat(6)).Quo(y, big.NewFloat(13))).Float64()
// 			yCanvas = 600 - yCanvas*600

// 			pointVisual := canvas.NewCircle(color.NRGBA{R: 255, G: 100, B: 100, A: 255})
// 			pointVisual.Resize(fyne.NewSize(7, 7))
// 			pointVisual.Move(fyne.NewPos(float32(xCanvas), float32(yCanvas)))

// 			graphContent.Add(pointVisual)
// 		}

// 		graphContent.Refresh()
// 	}

// 	// Add a button to redraw the graph
// 	redrawButton := widget.NewButton("Redraw", func() {
// 		redraw()
// 	})

// 	// Layout
// 	inputContent := container.NewVBox(
// 		widget.NewForm(
// 			widget.NewFormItem("A", aEntry),
// 			widget.NewFormItem("B", bEntry),
// 			widget.NewFormItem("p", pEntry),
// 		),
// 		redrawButton,
// 	)

// 	finiteFieldContent := container.NewHSplit(
// 		inputContent, // Left: Inputs for A and B
// 		graphContent, // Right: Graph finiteFieldContent
// 	)
// 	finiteFieldContent.Offset = 0.3

// 	finiteFieldVisWindow.SetContent(finiteFieldContent)
// 	finiteFieldVisWindow.Resize(fyne.NewSize(900, 650))
// 	redraw() // Draw the initial graph
// 	finiteFieldVisWindow.ShowAndRun()

// 	return nil
// }

package main

import (
	"fmt"
	"image/color"
	"math"
	"math/big"
	"os"

	"github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"
)

const (
	graphSize    = 600.0 // assuming a square graph
	displayWidth = 300.0 // width of left hand side display details
)

var (
	graphSizeInt, oneInt  *big.Int = big.NewInt(graphSize), big.NewInt(1)
	graphSizeRat, fourRat *big.Rat = big.NewRat(graphSize, 1), big.NewRat(4, 1)
	halfRat, graphStep    *big.Rat = big.NewRat(1, 2), big.NewRat(1, graphSize)
)

func main() {
	logger := utils.InitialiseLogger("[ECVIS]")
	logger.Debug("[ECVIS] starting function main")

	err := run(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

func run(logger *logrus.Logger) error {
	logger.Debug("[ECVIS] starting function run")

	myApp := app.New()
	finiteFieldVisWindow := myApp.NewWindow("Elliptic Curve Visualisation in Finite Field")
	realToInfinityVisWindow := myApp.NewWindow("Elliptic Curve Visualisation extended to Infinity")

	// initialise curve parameters and create the curve
	a, b, p := big.NewInt(-2), big.NewInt(1), big.NewInt(13)
	// aRat, bRat, pRat := new(big.Rat).SetInt(a), new(big.Rat).SetInt(b), new(big.Rat).SetInt(p)
	// halfPRat := new(big.Rat).Mul(pRat, halfRat)
	curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

	// get the roots of the curve
	roots, err := curve.SolveCubic(big.NewInt(-6)) // shift down by p/2 (rounded)
	if err != nil {
		return err
	}

	// calculate finitePoints on the curve
	logger.Debug("[ECVIS] calculating finitePoints")
	finitePoints, realPoints, err := finiteintfield.CalculatePoints(curve)
	if err != nil {
		finiteFieldVisWindow.SetContent(canvas.NewText("Error calculating finitePoints: "+err.Error(), color.White))
		return err
	}

	// TODO: not sure if this works - have to revisit
	logger.Debug("[ECVIS] get ordered finitePoints")
	orderedPoints := finiteintfield.OrderIntPoints(curve, finitePoints)
	orderedRealPoints := finiteintfield.OrderRealPoints(curve.GetEC(), realPoints)

	logger.Debug("[ECVIS] creating canvases")

	// FINITE FIELD CANVAS
	finiteFieldContent := drawFiniteFieldCanvas(logger, a, b, p, roots, finitePoints, orderedPoints, realPoints, orderedRealPoints)

	logger.Debug("Setting Finite Field Canvas Split Proportion")
	finiteFieldContent.Offset = 0.3 // Set split proportion
	logger.Debug("Setting Finite Field Canvas Content")
	finiteFieldVisWindow.SetContent(finiteFieldContent)

	logger.Debug("Setting Finite Field Canvas window size")
	// Set window size and display
	finiteFieldVisWindow.Resize(fyne.NewSize(graphSize+displayWidth, graphSize))

	// Show the window
	finiteFieldVisWindow.Show()

	// Run the event loop explicitly
	go func() {
		finiteFieldVisWindow.SetOnClosed(func() {
			myApp.Quit()
		})
	}()

	// REAL NUMBER CANVAS
	// ...
	realNumberContent := drawRealNumberCanvas(logger, a, b, p, roots, finitePoints)

	logger.Debug("Setting Real Number Canvas Split Proportion")
	realNumberContent.Offset = 0.3 // Set split proportion
	logger.Debug("Setting Real Number Canvas Content")
	realToInfinityVisWindow.SetContent(realNumberContent)

	logger.Debug("Setting Real Number Canvas window size")
	// Set window size and display
	realToInfinityVisWindow.Resize(fyne.NewSize(graphSize+displayWidth, graphSize))

	// Show the window
	realToInfinityVisWindow.Show()

	// Run the event loop explicitly
	go func() {
		realToInfinityVisWindow.SetOnClosed(func() {
			myApp.Quit()
		})
	}()

	// START APP
	logger.Debug("Running App")
	myApp.Run()

	return nil
}

func drawFiniteFieldCanvas(logger *logrus.Logger, a, b, p *big.Int, roots []*big.Rat, finitePoints, orderedPoints [][2]*big.Int, realPoints, orderedRealPoints [][2]*big.Rat) *container.Split {
	logger.Debug("[ECVIS] starting function drawFiniteFieldCanvas")

	aRat, bRat, pRat := new(big.Rat).SetInt(a), new(big.Rat).SetInt(b), new(big.Rat).SetInt(p)
	halfPRat := new(big.Rat).Mul(pRat, halfRat)

	// constants for finite field canvas dimensions
	scaleRat := new(big.Rat).SetFrac(graphSizeInt, p)
	negP := new(big.Int).Neg(p)
	minP := finiteintfield.Div2RoundUp(negP)
	maxP := finiteintfield.Div2RoundUp(p)

	// create finitePoints and display details
	finiteFieldDrawingPoints := make([]fyne.CanvasObject, 0)
	finiteFieldGraphDetails := fmt.Sprintf("A: %s, B: %s, P: %s", a.String(), b.String(), p.String()) // for text output to screen

	logger.Debug("Starting Finite Field point loop")

	// TODO: split loop sections into functions?
	// TODO: poition finitePoints and line with offset so it cetnres where it should better
	// TODO: grid
	// TODO: ticks at whole number or steps of whole numbers (based on size of p)
	finiteFieldPointDetails := ""
	numZeros := big.NewInt(0)
	for _, point := range finitePoints {
		// add text description of point
		finiteFieldPointDetails += fmt.Sprintf("(%v, %v), ", point[0], point[1])

		// count zeros (when y is 0)
		if point[1].Sign() == 0 {
			numZeros.Add(numZeros, oneInt)
		}

		// convert to big.Rat for later calcaultion
		x := new(big.Rat).SetInt(point[0])
		y := new(big.Rat).SetInt(point[1])

		// Convert x, y to canvas coordinates
		xCanvas := new(big.Rat).Add(x, halfPRat)
		xCanvas.Mul(xCanvas, scaleRat)
		xCanvas.Sub(xCanvas, fourRat)

		yCanvas := new(big.Rat).Add(y, halfPRat)
		yCanvas.Mul(yCanvas, scaleRat)
		yCanvas.Sub(graphSizeRat, yCanvas) // Flip y-axis to origin at the bottom-left
		yCanvas.Sub(yCanvas, fourRat)

		xCanvasFloat, _ := xCanvas.Float32()
		yCanvasFloat, _ := yCanvas.Float32()

		// Create a visual representation for the point
		fynePoint := canvas.NewCircle(color.NRGBA{R: 255, G: 100, B: 100, A: 255})
		fynePoint.Resize(fyne.NewSize(7, 7))
		fynePoint.Move(fyne.NewPos(xCanvasFloat, yCanvasFloat))

		// Store the point's visual and details
		finiteFieldDrawingPoints = append(finiteFieldDrawingPoints, fynePoint)
	}

	orderedPointDetails := ""
	for _, orderedPoint := range orderedPoints {
		x := new(big.Int).Set(orderedPoint[0])
		y := new(big.Int).Set(orderedPoint[1])

		// add text description of point
		orderedPointDetails += fmt.Sprintf("(%s, %s), ", x, y)
	}

	finiteFieldRealPontDetails := ""
	for _, realPoint := range realPoints {
		x := new(big.Rat).Set(realPoint[0])
		y := new(big.Rat).Set(realPoint[1])

		// add text description of point
		yFloat, _ := y.Float32()
		finiteFieldRealPontDetails += fmt.Sprintf("(%v, %f), ", x, yFloat)
	}

	finiteFieldOrderedRealPointsDetails := ""
	for _, orderedRealPoint := range orderedRealPoints {
		x := new(big.Rat).Set(orderedRealPoint[0])
		y := new(big.Rat).Set(orderedRealPoint[1])

		// add text description of point
		yFloat, _ := y.Float32()
		finiteFieldOrderedRealPointsDetails += fmt.Sprintf("(%v, %f), ", x, yFloat)
	}

	for x := new(big.Rat).SetInt(minP); x.Cmp(new(big.Rat).SetInt(maxP)) < 0; x.Add(x, graphStep) {
		axPlusB := new(big.Rat).Mul(aRat, x) // ax
		axPlusB.Add(axPlusB, bRat)           // ax + b

		y := new(big.Rat).Mul(x, x) // x^2
		y.Mul(y, x)                 // x^3
		y.Add(y, axPlusB)           // x^3 + ax + b
		yFloat64, _ := y.Float64()  // (x^3 + ax + b) - as float64
		if yFloat64 < 0 {           // if calculated y position is < 0
			continue // just go to next x value
		}
		y = new(big.Rat).SetFloat64(math.Sqrt(yFloat64)) // sqrt(x^3 + ax + b)
		negY := new(big.Rat).Neg(y)                      // -sqrt(x^3 + ax + b)

		// Convert x, y to canvas coordinates
		xCanvas := new(big.Rat).Add(x, halfPRat) // x + 1/2 p
		xCanvas.Mul(xCanvas, scaleRat)           // (x + 1/2 p) * scale

		yCanvas := new(big.Rat).Add(y, halfPRat) // y + 1/2 p
		yCanvas.Mul(yCanvas, scaleRat)           // (y + 1/2 p) * scale
		yCanvas.Sub(graphSizeRat, yCanvas)       // graphSize - ((y + 1/2 p) * scale) (flip y-axis to origin at the bottom-left - add half the circle size)

		negYCanvas := new(big.Rat).Add(negY, halfPRat) // -y + 1/2 p
		negYCanvas.Mul(negYCanvas, scaleRat)           // (-y + 1/2 p) * scale
		negYCanvas.Sub(graphSizeRat, negYCanvas)       // graphSize - ((-y + 1/2 p) * scale) (flip y-axis to origin at the bottom-left - add half the circle size)

		xCanvasFloat, _ := xCanvas.Float32() // number should be 0 <= xCanvasFloat <= graphSize - don't worry about exactness
		yCanvasFloat, _ := yCanvas.Float32() // number should be 0 <= yCanvasFloat <= graphSize - don't worry about exactness

		// create visual representation for the +ve point on curve line
		fynePoint := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		fynePoint.Resize(fyne.NewSize(1, 1))

		fynePoint.Move(fyne.NewPos(xCanvasFloat, yCanvasFloat))
		finiteFieldDrawingPoints = append(finiteFieldDrawingPoints, fynePoint)

		// create visual representation for the -ve point on curve line
		negYCanvasFloat, _ := negYCanvas.Float32() // number should be 0 <= yCanvasFloat <= graphSize - don't worry about exactness
		negFynePoint := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		negFynePoint.Resize(fyne.NewSize(1, 1))

		negFynePoint.Move(fyne.NewPos(xCanvasFloat, negYCanvasFloat))
		finiteFieldDrawingPoints = append(finiteFieldDrawingPoints, negFynePoint)
	}

	logger.Debugf("roots: %+v", roots)

	// create a label for displaying point details
	numberOfPoints := len(finitePoints) + 1 // plus 1 for the point at inf
	pointLabel := widget.NewLabel(fmt.Sprintf(
		"Graph Deatils: %s\n\nPoint Details: %s\n\nOrdered Point Details: %s, \n\nReal Point Details: %s\n\nOrdered Real Points Details: <redacted for now>\n\nNumber of Points (including point at inf): %d\n\n",
		finiteFieldGraphDetails,
		finiteFieldPointDetails,
		orderedPointDetails,
		finiteFieldRealPontDetails,
		// finiteFieldOrderedRealPointsDetails,
		numberOfPoints,
	))
	pointLabel.Wrapping = fyne.TextWrapWord

	logger.Debugf("finiteFieldGraphDetails: %+v", finiteFieldGraphDetails)
	logger.Debugf("finiteFieldPointDetails: %+v", finiteFieldPointDetails)
	logger.Debugf("finiteFieldRealPontDetails: %+v", finiteFieldRealPontDetails)
	logger.Debugf("finiteFieldOrderedRealPointsDetails: %+v", finiteFieldOrderedRealPointsDetails)

	logger.Debug("Creating Finite Field Canvas")
	// create axes for the finite field canvas
	xAxis := canvas.NewLine(color.Gray{Y: 123})
	xAxis.StrokeWidth = 1
	xAxis.Position1 = fyne.NewPos(0, graphSize/2) // Horizontal line
	xAxis.Position2 = fyne.NewPos(graphSize, graphSize/2)

	scaleFloat, _ := scaleRat.Float32()

	yAxis := canvas.NewLine(color.Gray{Y: 123})
	yAxis.StrokeWidth = 1
	yAxis.Position1 = fyne.NewPos((float32(p.Int64())/2)*scaleFloat, 0) // Vertical line
	yAxis.Position2 = fyne.NewPos((float32(p.Int64())/2)*scaleFloat, graphSize)

	logger.Debug("Assembling finiteFieldContent for layout")
	// Assemble finiteFieldContent with layout
	content := container.NewHSplit(
		container.NewVScroll(pointLabel),                                              // Left: Point details
		container.NewWithoutLayout(append(finiteFieldDrawingPoints, xAxis, yAxis)...), // Right: Visualisation
	)

	return content
}

func drawRealNumberCanvas(logger *logrus.Logger, a, b, p *big.Int, roots []*big.Rat, finitePoints [][2]*big.Int) *container.Split {
	logger.Debug("[ECVIS] starting function drawRealNumberCanvas")

	aRat, bRat, pRat := new(big.Rat).SetInt(a), new(big.Rat).SetInt(b), new(big.Rat).SetInt(p)
	halfPRat := new(big.Rat).Mul(pRat, halfRat)

	// constants for finite field canvas dimensions
	scaleRat := new(big.Rat).SetFrac(graphSizeInt, p)
	// negP := new(big.Int).Neg(p)
	minGraphPos := -int64(graphSize) / 2
	// minGraphBigInt := new(big.Int).SetInt64(minGraphPos)
	maxGraphPos := int64(graphSize) / 2
	// maxGraphBigInt := new(big.Int).SetInt64(maxGraphPos)

	// create realPoints and display details
	realNumberDrawingPoints := make([]fyne.CanvasObject, 0)

	infinityPoint := float64(maxGraphPos) + 1.0 // maxGrpahPos + 1 is the infinity point
	for xPos := minGraphPos; xPos < maxGraphPos; xPos++ {
		logger.Debugf("[ECVIS] calculating for xPos: %d", xPos)
		// convert pos to -inf < pos < inf value
		posPercentile := float64(xPos) / infinityPoint      // -1 < xPercentile < 1
		twoPiPosPercentile := 2.0 * math.Pi * posPercentile // -2 pi < twoPiXPercentile < 2 pi
		x := math.Tan(twoPiPosPercentile)                   // -inf < tan(twoPiXPercentile) < inf

		// importantly, x represents a number, -inf < x < inf - when pos == 0, x == 0
		// x is smooth, and monotonically increasting across minGraphPos < pos < maxGraphPos

		xRat := new(big.Rat).SetFloat64(x)      // x as a big.Rat
		axPlusB := new(big.Rat).Mul(aRat, xRat) // ax
		axPlusB.Add(axPlusB, bRat)              // ax + b

		// xRat gets used as the the x value for y calculation
		y := new(big.Rat).Mul(xRat, xRat) // x^2
		y.Mul(y, xRat)                    // x^3
		y.Add(y, axPlusB)                 // x^3 + ax + b
		yFloat64, _ := y.Float64()        // (x^3 + ax + b) - as float64
		if yFloat64 < 0 {                 // if calculated y position is < 0
			continue // just go to next x value
		}
		y = new(big.Rat).SetFloat64(math.Sqrt(yFloat64)) // sqrt(x^3 + ax + b)
		negY := new(big.Rat).Neg(y)                      // -sqrt(x^3 + ax + b)

		// convert y back into yPos
		arctanY, err := utils.ArctanAGM(y)
		if err != nil {
			logger.Errorf("error calculating arctan: %+v", arctanY)
		}
		arctanYFloat, _ := arctanY.Float32()
		yPos := 2 * math.Pi * arctanYFloat

		// Convert x, y to canvas coordinates
		// pos get used to position the new point (not xRat)
		xCanvas := new(big.Rat).Add(new(big.Rat).SetInt64(xPos), halfPRat) // x + 1/2 p
		xCanvas.Mul(xCanvas, scaleRat)                                     // (x + 1/2 p) * scale

		pFloat, _ := pRat.Float32()
		scaleFloat, _ := scaleRat.Float32()
		yCanvasFloat := yPos + (0.5 * pFloat)    // y + 1/2 p
		yCanvasFloat = yCanvasFloat * scaleFloat // (y + 1/2 p) * scale
		yCanvasFloat = graphSize - yCanvasFloat  // graphSize - ((y + 1/2 p) * scale) (flip y-axis to origin at the bottom-left - add half the circle size)

		negYCanvas := new(big.Rat).Add(negY, halfPRat) // -y + 1/2 p
		negYCanvas.Mul(negYCanvas, scaleRat)           // (-y + 1/2 p) * scale
		negYCanvas.Sub(graphSizeRat, negYCanvas)       // graphSize - ((-y + 1/2 p) * scale) (flip y-axis to origin at the bottom-left - add half the circle size)

		xCanvasFloat, _ := xCanvas.Float32() // number should be 0 <= xCanvasFloat <= graphSize - don't worry about exactness

		// create visual representation for the +ve point on curve line
		fynePoint := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		fynePoint.Resize(fyne.NewSize(1, 1))

		fynePoint.Move(fyne.NewPos(xCanvasFloat, yCanvasFloat))
		realNumberDrawingPoints = append(realNumberDrawingPoints, fynePoint)

		// create visual representation for the -ve point on curve line
		negYCanvasFloat, _ := negYCanvas.Float32() // number should be 0 <= yCanvasFloat <= graphSize - don't worry about exactness
		negFynePoint := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		negFynePoint.Resize(fyne.NewSize(1, 1))

		negFynePoint.Move(fyne.NewPos(xCanvasFloat, negYCanvasFloat))
		realNumberDrawingPoints = append(realNumberDrawingPoints, negFynePoint)
	}

	logger.Debugf("roots: %+v", roots)

	// create a label for displaying point details
	numberOfPoints := len(finitePoints) + 1 // plus 1 for the point at inf
	pointLabel := widget.NewLabel(fmt.Sprintf(
		"Number of Points (including point at inf): %d\n\n",
		numberOfPoints,
	))
	pointLabel.Wrapping = fyne.TextWrapWord

	logger.Debug("Creating Real Number Canvas")
	// create axes for the finite field canvas
	xAxis := canvas.NewLine(color.Gray{Y: 123})
	xAxis.StrokeWidth = 1
	xAxis.Position1 = fyne.NewPos(0, graphSize/2) // Horizontal line
	xAxis.Position2 = fyne.NewPos(graphSize, graphSize/2)

	scaleFloat, _ := scaleRat.Float32()

	yAxis := canvas.NewLine(color.Gray{Y: 123})
	yAxis.StrokeWidth = 1
	yAxis.Position1 = fyne.NewPos((float32(p.Int64())/2)*scaleFloat, 0) // Vertical line
	yAxis.Position2 = fyne.NewPos((float32(p.Int64())/2)*scaleFloat, graphSize)

	logger.Debug("Assembling realNumberContent for layout")
	// Assemble finiteFieldContent with layout
	content := container.NewHSplit(
		container.NewVScroll(pointLabel),                                             // Left: Point details
		container.NewWithoutLayout(append(realNumberDrawingPoints, xAxis, yAxis)...), // Right: Visualisation
	)

	return content
}
