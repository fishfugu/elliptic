package main

import (
	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"image/color"
	"math/big"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Elliptic Curve Visualization in Finite Field")

	// Initialize curve parameters and create curve
	a, b, p := big.NewInt(1), big.NewInt(1), big.NewInt(13)
	curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

	// Calculate points on the curve
	points, err := finiteintfield.CalculatePoints(*curve)
	if err != nil {
		myWindow.SetContent(canvas.NewText("Error calculating points: "+err.Error(), color.White))
		myWindow.ShowAndRun()
		return
	}

	// Set up the canvas
	w, h := float32(800), float32(600)
	fynePoints := make([]fyne.CanvasObject, 0)

	// Plot the points
	for _, point := range points {
		x, _ := new(big.Int).SetString(point[0], 10)
		y, _ := new(big.Int).SetString(point[1], 10)
		xCanvas := float32(x.Int64()) * (float32(w) / float32(p.Int64()))
		yCanvas := float32(h) - (float32(y.Int64()) * (float32(h) / float32(p.Int64()))) // Flipping y to have the origin at the bottom left
		fynePoint := canvas.NewCircle(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		fynePoint.Resize(fyne.NewSize(5, 5))
		fynePoint.Move(fyne.NewPos(xCanvas, yCanvas))
		fynePoints = append(fynePoints, fynePoint)
	}

	content := container.NewWithoutLayout(fynePoints...)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(w, h))
	myWindow.ShowAndRun()
}