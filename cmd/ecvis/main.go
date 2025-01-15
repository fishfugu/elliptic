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
	logger.Debug("starting function main")

	err := run(logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

func run(logger *logrus.Logger) error {
	logger.Debug("starting function run")

	myApp := app.New()
	myWindow := myApp.NewWindow("Elliptic Curve Visualisation in Finite Field")

	// initialise curve parameters and create the curve
	a, b, p := big.NewInt(-2), big.NewInt(1), big.NewInt(13)
	aRat, bRat, pRat := new(big.Rat).SetInt(a), new(big.Rat).SetInt(b), new(big.Rat).SetInt(p)
	halfPRat := new(big.Rat).Mul(pRat, halfRat)
	curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

	// get the roots of the curve
	roots, err := curve.SolveCubic(big.NewInt(-6)) // shift down by p/2 (rounded)
	if err != nil {
		return err
	}

	// calculate points on the curve
	points, err := finiteintfield.CalculatePoints(*curve, big.NewInt(-6), big.NewInt(-6)) // shift down by p/2 (rounded)
	if err != nil {
		myWindow.SetContent(canvas.NewText("Error calculating points: "+err.Error(), color.White))
		return nil
	}

	// constants for canvas dimensions
	scaleRat := new(big.Rat).SetFrac(graphSizeInt, p)

	negP := new(big.Int).Neg(p)
	minP := finiteintfield.Div2RoundUp(negP)
	maxP := finiteintfield.Div2RoundUp(p)

	// create points and display details
	fynePoints := make([]fyne.CanvasObject, 0)
	graphDetails := fmt.Sprintf("A: %s, B: %s, P: %s", a.String(), b.String(), p.String()) // for text output to screen

	// TODO: split loop sections into functions?
	// TODO: poition points and line with offset so it cetnres where it should better
	// TODO: grid
	// TODO: ticks at whole number or steps of whole numbers (based on size of p)
	pointDetails := ""
	numZeros := big.NewInt(0)
	for _, point := range points {
		x := new(big.Rat).SetInt(point[0])
		y := new(big.Rat).SetInt(point[1])
		logger.Debugf("x: %s, y: %s", x, y)

		// count zeros
		if y.Sign() == 0 {
			numZeros.Add(numZeros, oneInt)
		}

		pointDetails += fmt.Sprintf("(%v, %v), ", x, y)

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
		fynePoints = append(fynePoints, fynePoint)
	}

	logger.Warnf("roots[0]: %s, halfP: %s, graphStep: %s", roots[0].FloatString(10), halfPRat, graphStep)
	logger.Warnf("roots: %+v", roots)
	for x := new(big.Rat).SetFrac(minP, oneInt); x.Cmp(new(big.Rat).SetFrac(maxP, oneInt)) < 0; x.Add(x, graphStep) {
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
		fynePoints = append(fynePoints, fynePoint)

		// create visual representation for the -ve point on curve line
		negYCanvasFloat, _ := negYCanvas.Float32() // number should be 0 <= yCanvasFloat <= graphSize - don't worry about exactness
		negFynePoint := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
		negFynePoint.Resize(fyne.NewSize(1, 1))

		negFynePoint.Move(fyne.NewPos(xCanvasFloat, negYCanvasFloat))
		fynePoints = append(fynePoints, negFynePoint)
	}

	// create a label for displaying point details
	pointLabel := widget.NewLabel(fmt.Sprintf(
		"Graph Deatils: %s\nPoint Details: %s",
		graphDetails,
		pointDetails,
	))
	pointLabel.Wrapping = fyne.TextWrapWord

	// create axes for the canvas
	xAxis := canvas.NewLine(color.Gray{Y: 123})
	xAxis.StrokeWidth = 1
	xAxis.Position1 = fyne.NewPos(0, graphSize/2) // Horizontal line
	xAxis.Position2 = fyne.NewPos(graphSize, graphSize/2)

	scaleFloat, _ := scaleRat.Float32()

	yAxis := canvas.NewLine(color.Gray{Y: 123})
	yAxis.StrokeWidth = 1
	yAxis.Position1 = fyne.NewPos((float32(p.Int64())/2)*scaleFloat, 0) // Vertical line
	yAxis.Position2 = fyne.NewPos((float32(p.Int64())/2)*scaleFloat, graphSize)

	// Assemble content with layout
	content := container.NewHSplit(
		container.NewVScroll(pointLabel),                                // Left: Point details
		container.NewWithoutLayout(append(fynePoints, xAxis, yAxis)...), // Right: Visualisation
	)
	content.Offset = 0.3 // Set split proportion
	myWindow.SetContent(content)

	// Set window size and display
	myWindow.Resize(fyne.NewSize(graphSize+displayWidth, graphSize))

	// Show the window
	myWindow.Show()

	// Run the event loop explicitly
	go func() {
		myWindow.SetOnClosed(func() {
			myApp.Quit()
		})
	}()

	myApp.Run()

	return nil
}
