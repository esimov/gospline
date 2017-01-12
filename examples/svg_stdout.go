package main

import (
	"image/color"
	_ "image/jpeg"
	"time"
	"math/rand"
	"github.com/esimov/gospline"
	"os"
)

func main()  {
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	var points [][]float64

	for i := 0; i < 20; i++ {
		x := randInt(0, 800, rng)
		y := randInt(0, 800, rng)
		point := []float64{float64(x), float64(y)}
		points = append(points, point)
	}

	svg := &spline.SVG{
		Width: 800,
		Height: 800,
		Title: "BSpline",
		Lines: []spline.Line{},
		Color: color.NRGBA{R:255,G:0,B:0,A:255},
		Description: "Convert straight lines to curves",
		StrokeWidth: 2,
		StrokeLineCap: "round", //butt, round, square
	}

	svg.Draw(os.Stdout, points, false)
}

func randInt(min, max int, rng *rand.Rand) int {
	return rng.Intn(max-min) + min
}