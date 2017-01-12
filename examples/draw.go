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

	for i := 0; i < 20; i++ {
		x := randInt(0, 800, rng)
		y := randInt(0, 800, rng)
		point := []float64{float64(x), float64(y)}
		points = append(points, point)
	}
	/*points = [][]float64{
		[]float64{ 200, 50},
		[]float64{ 125, 75},
		[]float64{ 125, 190},
		[]float64{ 200, 200},
		[]float64{ 225, 150},
		[]float64{ 255, 150},
		[]float64{ 225, 150},
		[]float64{ 220, 200},
		[]float64{ 255, 200},
		[]float64{ 255, 150},
		[]float64{ 230, 150},
		[]float64{ 260, 155},
		[]float64{ 300, 150},
		[]float64{ 270, 150},
		[]float64{ 270, 200},
		[]float64{ 300, 200},
		[]float64{ 300, 190},
		[]float64{ 310, 50},
		[]float64{ 300, 75},
		[]float64{ 305, 200},
		[]float64{ 305, 190},
		[]float64{ 335, 170},
		[]float64{ 355, 150},
		[]float64{ 330, 150},
		[]float64{ 320, 190},
		[]float64{ 350, 200},
		[]float64{ 370, 180},
		[]float64{ 367, 150},
		[]float64{ 370, 230},
		[]float64{ 367, 200},
		[]float64{ 375, 150},
		[]float64{ 405, 150},
		[]float64{ 400, 180},
		[]float64{ 375, 170},
		[]float64{ 400, 180},
		[]float64{ 445, 150},
		[]float64{ 420, 150},
		[]float64{ 415, 190},
		[]float64{ 440, 200},
		[]float64{ 460, 170},
		[]float64{ 460, 150},
		[]float64{ 460, 200},
		[]float64{ 460, 155},
		[]float64{ 490, 155},
		[]float64{ 490, 200},
	}*/
	/*points := [][]float64{
		[]float64{ 100.0, 75.0 },
		[]float64{ 165.0, 50.0 },
		[]float64{ 435.0, 50.0 },
		[]float64{ 500.0, 75.0 },
		[]float64{ 500.0, 175.0 },
		[]float64{ 450.0, 200.0 },
		[]float64{ 290.0, 200.0 },
		[]float64{ 250.0, 210.0 },
		[]float64{ 220.0, 250.0 },
		[]float64{ 130.0, 270.0 },
		[]float64{ 185.0, 240.0 },
		[]float64{ 200.0, 205.0 },
		[]float64{ 170.0, 200.0 },
		[]float64{ 125.0, 200.0 },
		[]float64{ 100.0, 175.0 },
	}*/

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

	raster := &spline.Image{
		Width : 800,
		Height : 800,
		Color : color.NRGBA{R:255,G:0,B:0,A:255},
	}

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
				return
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