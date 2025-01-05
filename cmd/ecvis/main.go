package main

import (
	"fmt"
	"image/color"
	"math/big"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"elliptic/pkg/utils"
)

func main() {
	logger := utils.InitialiseLogger("[ECVIS/MAIN]")
	logger.Debug("starting function main")

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic, error: %v\n", err)
		panic(fmt.Sprintf("panic, error: %v\n", err))
	}
}

func run() error {
	logger := utils.InitialiseLogger("[ECVIS/RUN]")
	logger.Debug("starting function run")

	myApp := app.New()
	myWindow := myApp.NewWindow("Elliptic Curve Visualisation in Finite Field")

	// Initialise curve parameters and create the curve
	a, b, p := big.NewInt(1), big.NewInt(1), big.NewInt(13)
	curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

	// Calculate points on the curve
	points, err := finiteintfield.CalculatePoints(*curve)
	if err != nil {
		myWindow.SetContent(canvas.NewText("Error calculating points: "+err.Error(), color.White))
		myWindow.ShowAndRun()
		return nil
	}

	// Constants for canvas dimensions
	const canvasSize = 600.0 // Assuming a square canvas (600x600)
	pFloat, _ := p.Float64() // TODO: Deal with "accuracy" details
	scale := canvasSize / pFloat

	// Create points and display details
	fynePoints := make([]fyne.CanvasObject, len(points))
	pointDetails := ""

	for i, point := range points {
		x, _ := point[0].Float64()
		y, _ := point[1].Float64()

		// Convert x, y to canvas coordinates
		xCanvas := x * scale
		yCanvas := canvasSize - (y * scale) // Flip y-axis to origin at the bottom-left

		// Create a visual representation for the point
		fynePoint := canvas.NewCircle(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		fynePoint.Resize(fyne.NewSize(5, 5))
		fynePoint.Move(fyne.NewPos(float32(xCanvas), float32(yCanvas)))

		// Store the point's visual and details
		fynePoints[i] = fynePoint
		pointDetails += fmt.Sprintf("(%v, %v), ", x, y)
	}

	// Create a label for displaying point details
	pointLabel := widget.NewLabel(pointDetails)
	pointLabel.Wrapping = fyne.TextWrapWord

	// Create axes for the canvas
	xAxis := canvas.NewLine(color.Gray{Y: 123})
	xAxis.StrokeWidth = 1
	xAxis.Position1 = fyne.NewPos(0, canvasSize/2) // Horizontal line
	xAxis.Position2 = fyne.NewPos(canvasSize, canvasSize/2)

	yAxis := canvas.NewLine(color.Gray{Y: 123})
	yAxis.StrokeWidth = 1
	yAxis.Position1 = fyne.NewPos(canvasSize/2, 0) // Vertical line
	yAxis.Position2 = fyne.NewPos(canvasSize/2, canvasSize)

	// Assemble content with layout
	content := container.NewHSplit(
		container.NewVScroll(pointLabel),                                // Left: Point details
		container.NewWithoutLayout(append(fynePoints, xAxis, yAxis)...), // Right: Visualisation
	)
	content.Offset = 0.3 // Set split proportion
	myWindow.SetContent(content)

	// Set window size and display
	myWindow.Resize(fyne.NewSize(canvasSize, canvasSize))
	myWindow.ShowAndRun()

	return nil
}
