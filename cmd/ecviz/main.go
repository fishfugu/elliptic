package main

import (
	"elliptic/pkg/bigarith"
	"elliptic/pkg/ellipticcurve"
	"elliptic/pkg/finiteintfield"
	"fmt"
	"image/color"
	"math/big"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Elliptic Curve Visualisation in Finite Field")

	// Initialise curve parameters
	aEntry := widget.NewEntry()
	aEntry.SetPlaceHolder("Coefficient A")

	bEntry := widget.NewEntry()
	bEntry.SetPlaceHolder("Coefficient B")

	pEntry := widget.NewEntry()
	pEntry.SetPlaceHolder("Prime modulus P")

	// Create a container to hold the entries and button
	inputs := container.NewVBox(
		widget.NewLabel("Enter Elliptic Curve Parameters:"),
		widget.NewForm(
			widget.NewFormItem("A:", aEntry),
			widget.NewFormItem("B:", bEntry),
			widget.NewFormItem("P:", pEntry),
		),
	)

	// Create a container to hold the visualisation
	visualisation := container.NewWithoutLayout()

	// Redraw function to update the visualisation
	redraw := func() {
		a, okA := new(big.Int).SetString(aEntry.Text, 10)
		b, okB := new(big.Int).SetString(bEntry.Text, 10)
		p, okP := new(big.Int).SetString(pEntry.Text, 10)
		if !okA || !okB || !okP {
			visualisation.Objects = []fyne.CanvasObject{
				canvas.NewText("Invalid input: Please enter valid integers for A, B, and P.", color.White),
			}
			visualisation.Refresh()
			return
		}

		okPrime := bigarith.ProbablyPrime(p)
		if !okPrime {
			visualisation.Objects = []fyne.CanvasObject{
				canvas.NewText("Invalid input: P does not qppear to be prime.", color.White),
			}
			visualisation.Refresh()
			return
		}

		curve := ellipticcurve.NewFiniteFieldEC(a, b, p)

		// Calculate points on the curve
		points, err := finiteintfield.CalculatePoints(*curve)
		if err != nil {
			visualisation.Objects = []fyne.CanvasObject{
				canvas.NewText("Error calculating points: "+err.Error(), color.White),
			}
			visualisation.Refresh()
			return
		}

		// Get the size of the box inside the window
		w := float32(myWindow.Canvas().Size().Width) - 100
		h := float32(myWindow.Canvas().Size().Height) - 250 // Adjust height to account for input fields and button
		fynePoints := make([]fyne.CanvasObject, 0)

		// Add axes
		axes := addAxes(p.Int64(), w, h)
		fynePoints = append(fynePoints, axes...)

		// Plot the points
		for _, point := range points {
			x, _ := new(big.Int).SetString(point[0], 10)
			y, _ := new(big.Int).SetString(point[1], 10)
			xCanvas := float32(x.Int64()) * (float32(w) / float32(p.Int64()))
			yCanvas := float32(h) - (float32(y.Int64()) * (float32(h) / float32(p.Int64()))) - 2 // Flipping y to have the origin at the bottom left
			newPoint := canvas.NewCircle(color.NRGBA{R: 255, G: 100, B: 100, A: 255})
			newPoint.Resize(fyne.NewSize(3, 3))
			newPoint.Move(fyne.NewPos(xCanvas, yCanvas))
			fynePoints = append(fynePoints, newPoint)
		}

		visualisation.Objects = fynePoints
		visualisation.Refresh()
	}

	// Redraw button
	redrawButton := widget.NewButton("Redraw", func() {
		redraw()
	})

	// Layout the inputs and visualisation in a vertical box
	content := container.NewVBox(
		inputs,
		redrawButton,
		visualisation,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 700))

	// Refresh the visualisation when the window is resized
	// myWindow.Resize() = func(_ fyne.Size) {
	// 	redraw()
	// }

	go func() {
		widthBefore := myWindow.Content().Size().Width
		heightBefore := myWindow.Content().Size().Height
		for {
			width := myWindow.Content().Size().Width
			height := myWindow.Content().Size().Height
			if width != widthBefore || height != heightBefore {
				// window was resized
				fmt.Println("Window was resized!")
				fmt.Printf("Width: %s", strconv.Itoa(int(width)))
				fmt.Printf("height: %s", strconv.Itoa(int(height)))
			}
			widthBefore = width
			heightBefore = height
			time.Sleep(50 * time.Millisecond) // TODO: what should this be?
		}
	}()

	myWindow.ShowAndRun()
}

func addAxes(p int64, w, h float32) []fyne.CanvasObject {
	axes := make([]fyne.CanvasObject, 0)

	// Add x-axis and y-axis lines
	xAxis := canvas.NewLine(color.Gray{Y: 0x77})
	xAxis.Position1 = fyne.NewPos(0, h)
	xAxis.Position2 = fyne.NewPos(w, h)
	axes = append(axes, xAxis)

	yAxis := canvas.NewLine(color.Gray{Y: 0x77})
	yAxis.Position1 = fyne.NewPos(700, 0)
	yAxis.Position2 = fyne.NewPos(700, h)
	axes = append(axes, yAxis)

	// Add x-axis labels
	for i := int64(0); i <= p; i++ {
		label := canvas.NewText(fmt.Sprintf("%d", i), color.White)
		label.TextSize = 12
		label.Move(fyne.NewPos(float32(i)*(w/float32(p)), h+5))
		axes = append(axes, label)
	}

	// Add y-axis labels
	for i := int64(0); i <= p; i++ {
		label := canvas.NewText(fmt.Sprintf("%d", i), color.White)
		label.TextSize = 12
		label.Move(fyne.NewPos(750, h-float32(i)*(h/float32(p))-6))
		axes = append(axes, label)
	}

	return axes
}
