package spline

import (
	"image"
	"image/color"
	"image/png"
	_ "image/jpeg"
	"os"
	"fmt"
	"text/template"
	"io"
	"image/draw"
)

type ImageDrawer interface {
	Draw(io.Writer, [][]float64, bool)
}

type Point struct {
	X float64
	Y float64
}

type Line struct {
	Start 	Point
	End 	Point
}

type SVG struct {
	Width 		int
	Height 		int
	Title 		string
	Lines 		[]Line
	Color 		color.NRGBA
	Description 	string
	StrokeLineCap 	string
	StrokeWidth 	float64
}

type Image struct {
	Width 	int
	Height 	int
	Color 	color.NRGBA
}

func Pt(x float64, y float64) Point {
	return Point{X: x, Y: y}
}

func (svg *SVG) Draw(output io.Writer, points [][]float64, debug bool) {
	const TEMPLATE = `<?xml version="1.0" ?>
	<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN"
	  "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
	<svg width="{{.Width}}px" height="{{.Height}}px" viewBox="0 0 {{.Width}} {{.Height}}"
	     xmlns="http://www.w3.org/2000/svg" version="1.1">
	  <title>{{.Title}}</title>
	  <desc>{{.Description}}</desc>
	  <!-- Points -->
	  <g stroke="rgba({{.Color.R}},{{.Color.G}},{{.Color.B}},{{.Color.A}})" stroke-linecap="{{.StrokeLineCap}}" stroke-width="{{.StrokeWidth}}" fill="none">
	    {{range .Lines}}
		<path d="M{{.Start.X}},{{.Start.Y}} L{{.End.X}},{{.End.Y}}" />
	    {{end}}</g>
	</svg>`

	var interpol []float64
	var lines []Line

	spline := NewBSpline(points, 3, false)
	spline.Init()
	oldX := spline.Interpolate(0.0, 0.1)[0]
	oldY := spline.Interpolate(0.0, 0.1)[1]
	pointStart := Pt(oldX, oldY)

	for t := 0.0; t <= 1.0; t += 0.002 {
		interpol = spline.Interpolate(t, 0.1)
		x := interpol[0]
		y := interpol[1]
		pointEnd := Pt(x, y)
		oldX, oldY = x, y
		lines = append(lines, []Line{Line{pointStart, pointEnd}}...)
		pointStart = Pt(oldX, oldY)
	}

	svg.Lines = lines

	tmpl := template.Must(template.New("svg").Parse(TEMPLATE))
	if err := tmpl.Execute(output, svg); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (img *Image) Draw(output io.Writer, points [][]float64, debug bool) {
	matrix := Identity()
	canvas := Canvas{
		image.NewRGBA(image.Rect(0, 0, img.Width, img.Height)),
		matrix,
	}
	spline := NewBSpline(points, 3, false)
	spline.Init()
	oldX := spline.Interpolate(0.0, 0.5)[0]
	oldY := spline.Interpolate(0.0, 0.5)[1]
	var interpol []float64
	var raster *Canvas

	// Draw white background
	background := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas.Image, canvas.Image.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)

	// Draw original line
	if debug {
		for i:=0; i < len(points)-1; i++ {
			raster = canvas.DrawLine(
				points[i][0],
				points[i][1],
				points[i+1][0],
				points[i+1][1],
				color.NRGBA{R:155,G:155,B:155,A:70 }, false)
		}
	}

	for t := 0.0; t <= 1.0; t += 0.001 {
		interpol = spline.Interpolate(t, 0.5)
		x := interpol[0]
		y := interpol[1]
		raster = canvas.DrawLine(oldX, oldY, x, y, img.Color, true)
		oldX, oldY = x, y
	}

	png.Encode(output, raster)
}