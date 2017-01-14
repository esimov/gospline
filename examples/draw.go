package main

import (
	"image/color"
	_ "image/jpeg"
	"time"
	"math/rand"
	"net/http"
	"github.com/esimov/gospline"
	"os"
	"log"
)

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main()  {
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	var points [][]float64
	var width, height int = 1200, 800

	for i := 0; i < 40; i++ {
		x := randInt(0, width, rng)
		y := randInt(0, height, rng)
		point := []float64{float64(x), float64(y)}
		points = append(points, point)
	}
	svg := &spline.SVG{
		Width: width,
		Height: height,
		Title: "BSpline",
		Lines: []spline.Line{},
		Color: color.NRGBA{R:255,G:0,B:0,A:255},
		Description: "Convert straight lines to curves",
		StrokeWidth: 4,
		StrokeLineCap: "round", //butt, round, square
	}

	raster := &spline.Image{
		Width : width,
		Height : height,
		Color : color.NRGBA{R:255,G:0,B:0,A:255},
	}
	output, _ := os.Create("./samples/curve.png")
	defer output.Close()
	raster.Draw(output, points, true)

	drawers := []spline.ImageDrawer{svg, raster}
	for _, drawer := range drawers {
		switch drawer.(type) {
		case *spline.SVG:
			if (len(os.Args) > 1 && os.Args[1] == "--web") {
				handler := func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/svg+xml")
					drawer.Draw(w, points, false)
				}
				http.HandleFunc("/", handler)
				log.Fatal(http.ListenAndServe("localhost:8000", nil))
			}
		case *spline.Image:
			output, _ := os.Create("./samples/curve_" + randSeq(4, rng) +".png")
			defer output.Close()
			drawer.Draw(output, points, true)
		}
	}
}

func randInt(min, max int, rng *rand.Rand) int {
	return rng.Intn(max-min) + min
}

func randSeq(n int, rng *rand.Rand) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randInt(0, len(letters), rng)]
	}
	return string(b)
}