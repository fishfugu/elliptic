package main

import (
	"fmt"
	"image/color"
	"math/big"

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
	logger := utils.InitialiseLogger("[ECVIS]")

	err := run()
	utils.LogOnError(logger, err, fmt.Sprintf("error percolated to main, error: %v\n", err), false)
}

func run() error {
	myApp := app.New()
	myWindow := myApp.NewWindow("Elliptic Curve Visualisation in Finite Field")

	// Initialise curve parameters and create curve
	a, b, p := new(big.Int).SetInt64(1), new(big.Int).SetInt64(1), new(big.Int).SetInt64(13)
	curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

	// Calculate points on the curve
	points, err := finiteintfield.CalculatePoints(*curve)
	if err != nil {
		myWindow.SetContent(canvas.NewText("Error calculating points: "+err.Error(), color.White))
		myWindow.ShowAndRun()
		return nil
	}

	// Set up the canvas
	w, h := float32(600), float32(600)
	fynePoints := make([]fyne.CanvasObject, 0)
	pointDetails := ""

	// Plot the points and collect point details
	for _, point := range points {
		x := new(big.Int).Set(point[0]) // dupe point[0] into x
		y := new(big.Int).Set(point[1]) // dupe point[1] into y
		// Convert p to int64
		// TODO: deal with "accuracy" details - i.e. the "_" return value
		pFloat, _ := new(big.Int).Set(p).Float64()
		xCanvas := float32(x.Int64()) * (float32(w) / float32(pFloat))
		yCanvas := float32(h) - (float32(y.Int64()) * (float32(h) / float32(pFloat))) // Flipping y to have the origin at the bottom left
		fynePoint := canvas.NewCircle(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		fynePoint.Resize(fyne.NewSize(5, 5))
		fynePoint.Move(fyne.NewPos(xCanvas, yCanvas))
		fynePoints = append(fynePoints, fynePoint)
		pointDetails += fmt.Sprintf("(%v, %v), ", x, y)
	}

	// Create a text label for displaying the point details
	pointLabel := widget.NewLabel(pointDetails)
	pointLabel.Wrapping = fyne.TextWrapWord

	// Create canvas for x-y axes
	axes := canvas.NewLine(color.Gray{Y: 123})
	axes.StrokeWidth = 1
	axes.Position1 = fyne.NewPos(0, h/2) // horizontal line
	axes.Position2 = fyne.NewPos(w, h/2)
	axes2 := canvas.NewLine(color.Gray{Y: 123})
	axes2.StrokeWidth = 1
	axes2.Position1 = fyne.NewPos(w/2, 0) // vertical line
	axes2.Position2 = fyne.NewPos(w/2, h)

	// Assemble content with layout
	split := container.NewHSplit(
		container.NewVScroll(pointLabel),
		container.NewWithoutLayout(append(fynePoints, axes, axes2)...),
	)
	split.Offset = 0.3
	myWindow.SetContent(split)
	myWindow.Resize(fyne.NewSize(w, h))
	myWindow.ShowAndRun()

	return nil
}
